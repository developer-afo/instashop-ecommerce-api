package order_service

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/developer-afo/instashop-ecommerce-api/dto"
	payment_gateway_dto "github.com/developer-afo/instashop-ecommerce-api/dto/payment_gateway"
	"github.com/developer-afo/instashop-ecommerce-api/lib/constants"
	"github.com/developer-afo/instashop-ecommerce-api/lib/helper"
	"github.com/developer-afo/instashop-ecommerce-api/models"
	"github.com/developer-afo/instashop-ecommerce-api/repository"
	order_repository "github.com/developer-afo/instashop-ecommerce-api/repository/order"
	core_service "github.com/developer-afo/instashop-ecommerce-api/service/core"
	finance_service "github.com/developer-afo/instashop-ecommerce-api/service/finance"
	payment_gateway_service "github.com/developer-afo/instashop-ecommerce-api/service/finance/payment_gateway"
	userService "github.com/developer-afo/instashop-ecommerce-api/service/user"
)

var (
	TransactionDescription = "Payment for order"
	TransactionShortDesc   = "order"
)

type OrderServiceInterface interface {
	CheckoutOrder(order dto.CreateOrderDTO) (string, int, error)
	CancelOrder(orderId uuid.UUID) error
	FindOrderById(uuid uuid.UUID) (dto.OrderDTO, error)
	FindOrderByReference(reference string) (dto.OrderDTO, error)
	FindAllOrders(pageable order_repository.OrderPageable) ([]dto.OrderDTO, repository.Pagination, error)
	VerifyOrderPayment(reference string) error
	ProcessOrder(orderId uuid.UUID) error
	OutForDelivery(orderId uuid.UUID) error
	Delivered(orderId uuid.UUID) error
}

type orderService struct {
	orderRepository           order_repository.OrderRepositoryInterface
	orderItemService          OrderItemServiceInterface
	orderStatusService        OrderStatusServiceInterface
	orderStatusHistoryService OrderStatusHistoryServiceInterface
	productService            core_service.ProductServiceInterface
	transactionService        finance_service.TransactionServiceInterface
	paymentGatewayService     payment_gateway_service.PaymentGatewayServiceInterface
	userService               userService.UserServiceInterface
}

func NewOrderService(
	orderRepository order_repository.OrderRepositoryInterface,
	orderItemService OrderItemServiceInterface,
	orderStatusService OrderStatusServiceInterface,
	orderStatusHistoryService OrderStatusHistoryServiceInterface,
	productService core_service.ProductServiceInterface,
	transactionService finance_service.TransactionServiceInterface,
	paymentGatewayService payment_gateway_service.PaymentGatewayServiceInterface,
	userService userService.UserServiceInterface,
) OrderServiceInterface {

	return &orderService{
		orderRepository:           orderRepository,
		orderItemService:          orderItemService,
		orderStatusService:        orderStatusService,
		orderStatusHistoryService: orderStatusHistoryService,
		productService:            productService,
		transactionService:        transactionService,
		paymentGatewayService:     paymentGatewayService,
		userService:               userService,
	}
}

func (o *orderService) ConvertToDTO(order models.Order) dto.OrderDTO {
	var orderDTO dto.OrderDTO

	orderDTO.ID = order.ID
	orderDTO.UserID = order.UserID
	orderDTO.TransactionID = order.TransactionID
	orderDTO.PaymentMethod = order.PaymentMethod
	orderDTO.Reference = order.Reference
	orderDTO.TotalPrice = order.TotalPrice
	orderDTO.StatusUUID = order.StatusID
	orderDTO.User = o.userService.ConvertToDTO(order.User)
	orderDTO.Status = o.orderStatusService.ConvertToDTO(order.Status)
	orderDTO.Transaction = o.transactionService.ConvertToDTO(order.Transaction)
	for _, item := range order.OrderItems {
		orderDTO.Items = append(orderDTO.Items, o.orderItemService.ConvertToDTO(item))
	}
	for _, item := range order.StatusHistory {
		orderDTO.StatusHistory = append(orderDTO.StatusHistory, o.orderStatusHistoryService.ConvertToDTO(item))
	}
	orderDTO.CreatedAt = order.CreatedAt
	orderDTO.UpdatedAt = order.UpdatedAt
	orderDTO.DeletedAt = order.DeletedAt.Time

	return orderDTO
}

func (o *orderService) ConvertToModel(orderDTO dto.OrderDTO) models.Order {
	var order models.Order

	order.ID = orderDTO.ID
	order.UserID = orderDTO.UserID
	order.TransactionID = orderDTO.TransactionID
	order.PaymentMethod = orderDTO.PaymentMethod
	order.Reference = orderDTO.Reference
	order.TotalPrice = orderDTO.TotalPrice
	order.StatusID = orderDTO.StatusUUID
	order.CreatedAt = orderDTO.CreatedAt
	order.UpdatedAt = orderDTO.UpdatedAt
	order.DeletedAt.Time = orderDTO.DeletedAt

	return order
}

// CheckoutOrder implements OrderServiceInterface.
func (o *orderService) CheckoutOrder(order dto.CreateOrderDTO) (string, int, error) {
	var orderDto dto.OrderDTO
	var totalPrice float64
	var paymentUrl string
	totalPrice, calcErr := o.CalculateTotalPrice(order)

	snowflake, err := helper.GenerateSnowflakeID()

	if err != nil {
		return "", constants.ServerErrorServiceUnavailable, err
	}

	if calcErr != nil {
		return "", constants.ServerErrorServiceUnavailable, calcErr
	}

	trans, paymentUrl, err := o.PayWithGateway(order.UserID, totalPrice, order.PaymentMethod)
	if err != nil {
		if err == payment_gateway_service.ErrPaymentInitialization {
			return "", constants.PaymentGatewayError, err
		}

		return "", constants.ServerErrorServiceUnavailable, err
	}

	orderStatus, err := o.orderStatusService.StatusOrderPlaced()
	if err != nil {
		return "", constants.ServerErrorServiceUnavailable, err
	}

	// Create order
	orderDto.UserID = order.UserID
	orderDto.TransactionID = trans.ID
	orderDto.CouponID = order.CouponID
	orderDto.StatusUUID = orderStatus.ID
	orderDto.PaymentMethod = order.PaymentMethod
	orderDto.Reference = helper.Int64ToString(snowflake)
	orderDto.TotalPrice = totalPrice

	// Save order
	newOrder := o.ConvertToModel(orderDto)

	newOrder, err = o.orderRepository.CreateOrder(newOrder)

	if err != nil {
		return "", constants.ServerErrorServiceUnavailable, err
	}

	// create order items
	if err := o.CreateOrderItems(newOrder.ID, order.Items); err != nil {
		return "", constants.ServerErrorServiceUnavailable, err
	}

	// remove quantity from product stock
	if err := o.RemoveQuantityFromProductStock(order.Items); err != nil {
		return "", constants.ServerErrorServiceUnavailable, err
	}

	// create status history
	_, err = o.orderStatusHistoryService.CreateOrderStatusHistory(dto.OrderStatusHistoryDTO{
		OrderUUID:  newOrder.ID,
		StatusUUID: newOrder.StatusID,
	})

	if err != nil {
		return "", constants.ServerErrorServiceUnavailable, err
	}

	return paymentUrl, constants.OrderPlacedSuccessfully, nil

}

// CancelOrder implements OrderServiceInterface.
func (o *orderService) CancelOrder(orderId uuid.UUID) error {
	// get order
	order, err := o.FindOrderById(orderId)

	if err != nil {
		return err
	}

	// check order statuses
	switch order.Status.ShortName {
	case ORDER_PLACED:
		break
	case AWAITING_CONFIRMATION:
		return fmt.Errorf("order is awaiting confirmation")
	case ORDER_PROCESSING:
		return fmt.Errorf("order is processing")
	case OUT_FOR_DELIVERY:
		return fmt.Errorf("order is out for delivery")
	case DELIVERED:
		return fmt.Errorf("order is already delivered")
	case CANCELLED:
		return fmt.Errorf("order is already cancelled")
	}

	// update order status to cancelled
	status, err := o.orderStatusService.StatusCancelled()
	if err != nil {
		return err
	}

	// update order status
	if err = o.UpdateOrderStatus(orderId, status.ID); err != nil {
		return err
	}

	return err
}

// Update order status to awaiting confirmation
func (o *orderService) ConfirmOrder(id uuid.UUID) error {

	// update order status to awaiting confirmation
	statusConfirm, err := o.orderStatusService.StatusAwaitingConfirmation()
	if err != nil {
		return err
	}

	err = o.UpdateOrderStatus(id, statusConfirm.ID)

	if err != nil {
		return err
	}

	return err
}

func (o *orderService) ProcessOrder(orderId uuid.UUID) error {
	// get order
	order, err := o.FindOrderById(orderId)

	if err != nil {
		return err
	}

	// check order statuses
	switch order.Status.ShortName {
	case ORDER_PLACED:
		return fmt.Errorf("order is not yet confirmed")
	case ORDER_PROCESSING:
		return fmt.Errorf("order is already processing")
	case OUT_FOR_DELIVERY:
		return fmt.Errorf("order is out for delivery")
	case DELIVERED:
		return fmt.Errorf("order is already delivered")
	case CANCELLED:
		return fmt.Errorf("order is cancelled")
	}

	// update order status to order processing
	statusProcessing, err := o.orderStatusService.StatusOrderProcessing()
	if err != nil {
		return err
	}

	err = o.UpdateOrderStatus(orderId, statusProcessing.ID)

	if err != nil {
		return err
	}

	return err
}

func (o *orderService) OutForDelivery(orderId uuid.UUID) error {
	// get order
	order, err := o.FindOrderById(orderId)

	if err != nil {
		return err
	}

	// check order statuses
	switch order.Status.ShortName {
	case ORDER_PLACED:
		return fmt.Errorf("order is not yet confirmed")
	case AWAITING_CONFIRMATION:
		return fmt.Errorf("order is awaiting confirmation")
	case OUT_FOR_DELIVERY:
		return fmt.Errorf("order is already out for delivery")
	case DELIVERED:
		return fmt.Errorf("order is already delivered")
	case CANCELLED:
		return fmt.Errorf("order is cancelled")
	}

	// update order status to out for delivery
	status, err := o.orderStatusService.StatusOutForDelivery()
	if err != nil {
		return err
	}

	// update order status
	if err = o.UpdateOrderStatus(orderId, status.ID); err != nil {
		return err
	}

	return err
}

func (o *orderService) Delivered(orderId uuid.UUID) error {
	// get order
	order, err := o.FindOrderById(orderId)

	if err != nil {
		return err
	}

	// check order statuses
	switch order.Status.ShortName {
	case ORDER_PLACED:
		return fmt.Errorf("order is not yet confirmed")
	case AWAITING_CONFIRMATION:
		return fmt.Errorf("order is awaiting confirmation")
	case ORDER_PROCESSING:
		return fmt.Errorf("order is processing")
	case DELIVERED:
		return fmt.Errorf("order is already delivered")
	case CANCELLED:
		return fmt.Errorf("order is cancelled")
	}

	// update order status to delivered
	status, err := o.orderStatusService.StatusDelivered()
	if err != nil {
		return err
	}

	// update order status
	if err = o.UpdateOrderStatus(orderId, status.ID); err != nil {
		return err
	}

	return err
}

// CreateOrderItems
func (o *orderService) CreateOrderItems(orderId uuid.UUID, items []dto.CreateOrderItemDTO) error {
	var orderItemDtos []dto.OrderItemDTO

	for _, item := range items {
		var price float64

		// Get product
		product, err := o.productService.FindProductByUUID(item.ProductUUID)
		if err != nil {
			return err
		}

		price = product.Price * float64(item.Quantity)

		orderItemDtos = append(orderItemDtos, dto.OrderItemDTO{
			OrderUUID:   orderId,
			ProductUUID: product.ID,
			Quantity:    item.Quantity,
			Price:       price,
		})
	}

	err := o.orderItemService.BatchCreateOrderItem(orderId.String(), orderItemDtos)

	if err != nil {
		return err
	}

	return nil
}

// verify order by payment reference
func (o *orderService) VerifyOrderPayment(reference string) error {
	// get transaction
	transaction, err := o.transactionService.FindTransactionByReference(reference)

	if err != nil {
		return err
	}

	if transaction.Status == finance_service.TransactionStatusSuccess {
		return nil
	}

	// verify transaction from payment gateway
	gatewayResp, err := o.paymentGatewayService.VerifyPayment(reference, transaction.Vendor)

	if err != nil {
		return err
	}

	order, err := o.orderRepository.FindOrderByTransactionId(transaction.ID)

	if err != nil {
		return err
	}

	if !gatewayResp.Status {
		return fmt.Errorf("payment verification failed: %s", gatewayResp.Message)
	}

	if gatewayResp.PaymentStatus == finance_service.TransactionStatusPending {
		return fmt.Errorf("payment verification is still pending: %s", gatewayResp.Message)
	}

	if gatewayResp.PaymentStatus == finance_service.TransactionStatusFailed {
		orderStatus, err := o.orderStatusService.StatusCancelled()
		if err != nil {
			return err
		}

		if err = o.UpdateOrderStatus(order.ID, orderStatus.ID); err != nil {
			return err
		}

		_, err = o.transactionService.FailTransaction(transaction.ID.String())
		if err != nil {
			return err
		}

		return nil
	}

	_, err = o.transactionService.ConfirmTransaction(transaction.ID.String())

	if err != nil {
		return err
	}

	err = o.ConfirmOrder(order.ID)

	return err
}

// FindOrderById implements OrderServiceInterface.
func (o *orderService) FindOrderById(uuid uuid.UUID) (dto.OrderDTO, error) {

	order, err := o.orderRepository.FindOrderById(uuid)
	if err != nil {
		return dto.OrderDTO{}, err
	}

	return o.ConvertToDTO(order), nil

}

// FindOrderByReference implements OrderServiceInterface.
func (o *orderService) FindOrderByReference(reference string) (dto.OrderDTO, error) {

	order, err := o.orderRepository.FindOrderByReference(reference)
	if err != nil {
		return dto.OrderDTO{}, err
	}

	return o.ConvertToDTO(order), nil
}

// FindAllOrders implements OrderServiceInterface.
func (o *orderService) FindAllOrders(pageable order_repository.OrderPageable) ([]dto.OrderDTO, repository.Pagination, error) {
	orders := []dto.OrderDTO{}

	orderModels, pagination, err := o.orderRepository.FindAllOrders(pageable)

	if err != nil {
		return nil, pagination, err
	}

	for _, order := range orderModels {
		orders = append(orders, o.ConvertToDTO(order))
	}

	return orders, pagination, nil

}

func (o *orderService) UpdateOrderStatus(orderId uuid.UUID, statusId uuid.UUID) error {

	order, err := o.orderRepository.FindOrderById(orderId)

	if err != nil {
		return err
	}

	orderStatus, err := o.orderStatusService.FindOrderStatusById(statusId)

	if err != nil {
		return err
	}

	order.StatusID = orderStatus.ID

	_, err = o.orderRepository.UpdateOrder(order)

	if err != nil {
		return err
	}

	_, err = o.orderStatusHistoryService.CreateOrderStatusHistory(dto.OrderStatusHistoryDTO{
		OrderUUID:  order.ID,
		StatusUUID: orderStatus.ID,
	})

	return err

}

func (o *orderService) CalculateTotalPrice(order dto.CreateOrderDTO) (float64, error) {
	var totalPrice float64

	// calculate total price in order items
	totalItemsPrice, err := o.CalculateTotalItemsPrice(order.Items)

	if err != nil {
		return 0, err
	}

	totalPrice += totalItemsPrice

	return totalPrice, nil
}

func (o *orderService) RemoveQuantityFromProductStock(items []dto.CreateOrderItemDTO) error {
	for _, item := range items {
		product, err := o.productService.FindProductByUUID(item.ProductUUID)

		if err != nil {
			return err
		}

		product.Stock = product.Stock - item.Quantity

		_, err = o.productService.UpdateProduct(product)

		if err != nil {
			return err
		}
	}

	return nil
}

func (o *orderService) CalculateTotalItemsPrice(items []dto.CreateOrderItemDTO) (float64, error) {
	var totalPrice float64

	for _, item := range items {
		var price float64

		// Get product
		product, err := o.productService.FindProductByUUID(item.ProductUUID)
		if err != nil {
			return 0, fmt.Errorf("product: %s is not found on this platform", product.Name)
		}

		// check product stock
		if product.Stock < item.Quantity {
			return 0, fmt.Errorf("product: %s stock is not enough", product.Name)
		}

		// calculate price
		// check if sales price is not zero
		if product.SlashPrice > 0 {
			price = product.SlashPrice * float64(item.Quantity)
		} else {
			price = product.Price * float64(item.Quantity)
		}

		// Calculate total price
		totalPrice += price
	}

	return totalPrice, nil
}

func (o *orderService) PayWithGateway(UserID uuid.UUID, amount float64, gateway string) (dto.TransactionDTO, string, error) {
	user, err := o.userService.FindUserById(UserID.String())

	if err != nil {
		return dto.TransactionDTO{}, "", err
	}

	// Create transaction
	newTransaction, err := o.transactionService.CreateTransaction(dto.TransactionDTO{
		UserID:      UserID,
		Amount:      amount,
		Type:        finance_service.TransactionTypeDebit,
		Description: TransactionDescription,
		ShortDesc:   TransactionShortDesc,
		Status:      finance_service.TransactionStatusPending,
		Method:      finance_service.TransactionMethodGateway,
		Vendor:      finance_service.TransactionVendorPayStack,
	})

	if err != nil {
		return dto.TransactionDTO{}, "", err
	}

	// Initialize payment
	initialize, err := o.paymentGatewayService.InitializePayment(payment_gateway_dto.PaymentInitializationDTO{
		Amount:    amount,
		Email:     user.Email,
		Reference: newTransaction.Reference,
		Gateway:   gateway,
	})

	if err != nil {
		return dto.TransactionDTO{}, "", err
	}

	if !initialize.Status {
		return dto.TransactionDTO{}, "", payment_gateway_service.ErrPaymentInitialization
	}

	return newTransaction, initialize.PaymentURL, nil
}
