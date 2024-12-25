package order_handler

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/developer-afo/instashop-ecommerce-api/dto"
	"github.com/developer-afo/instashop-ecommerce-api/handler"
	"github.com/developer-afo/instashop-ecommerce-api/lib/constants"
	"github.com/developer-afo/instashop-ecommerce-api/payload/request"
	"github.com/developer-afo/instashop-ecommerce-api/payload/response"
	order_repository "github.com/developer-afo/instashop-ecommerce-api/repository/order"
	order_service "github.com/developer-afo/instashop-ecommerce-api/service/order"
	order_validator "github.com/developer-afo/instashop-ecommerce-api/validator/order"
)

type orderHandler struct {
	orderService              order_service.OrderServiceInterface
	orderStatusHistoryService order_service.OrderStatusHistoryServiceInterface
	orderStatusService        order_service.OrderStatusServiceInterface
	validator                 order_validator.OrderValidator
}

type OrderHandlerInterface interface {
	CreateOrder(c *fiber.Ctx) error
	CancelOrder(c *fiber.Ctx) error
	GetUserOrders(c *fiber.Ctx) error
	GetAllOrders(c *fiber.Ctx) error
	VerifyOrderPayment(c *fiber.Ctx) error
	OrderProcessing(c *fiber.Ctx) error
	OutForDelivery(c *fiber.Ctx) error
	Delivered(c *fiber.Ctx) error
	StatusHistoryByOrderId(c *fiber.Ctx) error
	GetOrderStatuses(c *fiber.Ctx) error
}

func NewOrderHandler(
	orderService order_service.OrderServiceInterface,
	orderStatusHistoryService order_service.OrderStatusHistoryServiceInterface,
	orderStatusService order_service.OrderStatusServiceInterface,
) OrderHandlerInterface {
	return &orderHandler{
		orderService:              orderService,
		orderStatusHistoryService: orderStatusHistoryService,
		orderStatusService:        orderStatusService,
	}
}

func ConvertOrderStatusDTOToResponse(orderStatusDto dto.OrderStatusDTO) response.OrderStatusResponse {
	var resp response.OrderStatusResponse

	resp.Name = orderStatusDto.Name
	resp.ShortName = orderStatusDto.ShortName

	return resp
}

func ConvertOrderStatusHistoryDTOToResponse(orderStatusHistoryDto dto.OrderStatusHistoryDTO) response.OrderStatusHistoryResponse {
	var resp response.OrderStatusHistoryResponse

	resp.CreatedAt = orderStatusHistoryDto.CreatedAt
	resp.Status = ConvertOrderStatusDTOToResponse(orderStatusHistoryDto.Status)

	return resp
}

func (h *orderHandler) ConvertDTOtoResponse(orderDto dto.OrderDTO) response.OrderResponse {
	var orderResponse response.OrderResponse

	orderResponse.ID = orderDto.ID

	orderResponse.CreatedAt = orderDto.CreatedAt
	orderResponse.UpdatedAt = orderDto.UpdatedAt
	orderResponse.PaymentMethod = orderDto.PaymentMethod
	orderResponse.Reference = orderDto.Reference
	orderResponse.TotalPrice = orderDto.TotalPrice
	orderResponse.Transaction = response.TransactionResponse{
		ID:          orderDto.Transaction.ID,
		Reference:   orderDto.Transaction.Reference,
		Amount:      orderDto.Transaction.Amount,
		Status:      orderDto.Transaction.Status,
		Type:        orderDto.Transaction.Type,
		Description: orderDto.Transaction.Description,
		Method:      orderDto.Transaction.Method,
		Vendor:      orderDto.Transaction.Vendor,
		CreatedAt:   orderDto.Transaction.CreatedAt,
		UpdatedAt:   orderDto.Transaction.UpdatedAt,
	}
	for _, item := range orderDto.Items {
		orderResponse.OrderItems = append(orderResponse.OrderItems, response.OrderItemResponse{
			Product: response.ProductResponse{
				UUID:        item.Product.ID,
				Name:        item.Product.Name,
				Description: item.Product.Description,
				Price:       item.Product.Price,
				Images: func() []response.ImageResponse {
					var images []response.ImageResponse

					for _, image := range item.Product.Images {
						images = append(images, response.ImageResponse{
							Key: image.Key,
						})
					}

					return images
				}(),
			},
			Quantity: item.Quantity,
			Price:    item.Price,
		})
	}
	for _, statusHistory := range orderDto.StatusHistory {
		orderResponse.StatusHistory = append(orderResponse.StatusHistory, ConvertOrderStatusHistoryDTOToResponse(statusHistory))
	}
	orderResponse.Status = ConvertOrderStatusDTOToResponse(orderDto.Status)
	orderResponse.User = response.UserResponseData{
		FirstName: orderDto.User.FirstName,
		LastName:  orderDto.User.LastName,
		Email:     orderDto.User.Email,
	}

	return orderResponse
}

func (h *orderHandler) GeneratePageable(c *fiber.Ctx) (pageable order_repository.OrderPageable) {
	var resp response.Response
	basePageable := handler.GeneratePageable(c)

	pageable.Page = basePageable.Page
	pageable.Size = basePageable.Size
	pageable.SortBy = basePageable.SortBy
	pageable.SortDirection = basePageable.SortDirection
	pageable.Search = basePageable.Search

	pageable.Status = ""
	pageable.UserID = uuid.Nil
	pageable.FromDate = ""
	pageable.ToDate = ""

	if status := c.Query("status", ""); status != "" {
		pageable.Status = status
	}

	userID := c.Query("user_id", "")
	if userID != "" {
		userId, err := uuid.Parse(userID)

		pageable.UserID = userId

		if err != nil {
			resp.Status = constants.ClientUnProcessableEntity
			resp.Message = "User ID is not a valid UUID format"

			c.Status(http.StatusUnprocessableEntity).JSON(resp)
		}
	}

	fromDate := c.Query("from_date", "")
	if fromDate != "" {
		from, err := time.Parse("2006-01-02", fromDate)
		if err != nil {
			resp.Status = constants.ClientUnProcessableEntity
			resp.Message = "From date is not a valid date format"

			c.Status(http.StatusUnprocessableEntity).JSON(resp)
		}

		pageable.FromDate = from.Format("2006-01-02")
	}

	toDate := c.Query("to_date", "")
	if toDate != "" {
		to, err := time.Parse("2006-01-02", toDate)
		if err != nil {
			resp.Status = constants.ClientUnProcessableEntity
			resp.Message = "To date is not a valid date format"

			c.Status(http.StatusUnprocessableEntity).JSON(resp)
		}

		pageable.ToDate = to.Format("2006-01-02")
	}

	return pageable
}

func (h *orderHandler) CreateOrder(c *fiber.Ctx) error {
	var createOrderRequest request.CreateOrderRequest
	var createOrderDto dto.CreateOrderDTO
	var resp response.Response

	if err := c.BodyParser(&createOrderRequest); err != nil {
		resp.Status = constants.ClientRequestValidationError
		resp.Message = err.Error()

		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	if validation, err := h.validator.CreateOrderValidate(createOrderRequest); err != nil {
		resp.Status = constants.ClientRequestValidationError
		resp.Message = err.Error()
		resp.Data = validation

		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	createOrderDto.UserID = handler.GetUserId(c)
	createOrderDto.PaymentMethod = createOrderRequest.PaymentMethod

	for _, item := range createOrderRequest.Items {

		createOrderDto.Items = append(createOrderDto.Items, dto.CreateOrderItemDTO{
			ProductUUID: item.ProductID,
			Quantity:    item.Quantity,
		})
	}

	url, status, err := h.orderService.CheckoutOrder(createOrderDto)

	if err != nil {
		resp.Status = uint16(status)
		resp.Message = err.Error()

		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	resp.Status = constants.OrderPlacedSuccessfully
	resp.Message = "Order created successfully"
	resp.Data = map[string]interface{}{"payment_url": url}

	return c.Status(http.StatusCreated).JSON(resp)
}

func (h *orderHandler) CancelOrder(c *fiber.Ctx) error {
	var resp response.Response

	id := c.Params("order_id")

	orderId, err := uuid.Parse(id)
	if err != nil {
		resp.Status = constants.InvalidOrderID
		resp.Message = "Invalid order ID"

		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	err = h.orderService.CancelOrder(orderId)
	if err != nil {
		resp.Status = constants.ServerErrorInternal
		resp.Message = err.Error()

		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	resp.Status = http.StatusOK
	resp.Message = "Order cancelled successfully"

	return c.Status(http.StatusOK).JSON(resp)
}

func (h *orderHandler) GetUserOrders(c *fiber.Ctx) error {
	var resp response.Response
	var orderResponses []response.OrderResponse

	userId := handler.GetUserId(c)
	pageable := h.GeneratePageable(c)

	pageable.UserID = userId

	orders, pagination, err := h.orderService.FindAllOrders(pageable)
	if err != nil {
		resp.Status = constants.ServerErrorInternal
		resp.Message = err.Error()

		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	for _, order := range orders {
		orderResponses = append(orderResponses, h.ConvertDTOtoResponse(order))
	}

	resp.Status = http.StatusOK
	resp.Message = "Success"
	resp.Data = map[string]interface{}{"results": orderResponses, "pagination": pagination}

	return c.Status(http.StatusOK).JSON(resp)
}

func (h *orderHandler) GetAllOrders(c *fiber.Ctx) error {
	var resp response.Response
	var orderResponses []response.OrderResponse

	pageable := h.GeneratePageable(c)

	orders, pagination, err := h.orderService.FindAllOrders(pageable)
	if err != nil {
		resp.Status = constants.ServerErrorInternal
		resp.Message = err.Error()

		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	for _, order := range orders {
		orderResponses = append(orderResponses, h.ConvertDTOtoResponse(order))
	}

	resp.Status = http.StatusOK
	resp.Message = "Success"
	resp.Data = map[string]interface{}{"results": orderResponses, "pagination": pagination}

	return c.Status(http.StatusOK).JSON(resp)
}

func (h *orderHandler) VerifyOrderPayment(c *fiber.Ctx) error {
	var resp response.Response

	reference := c.Params("reference")

	err := h.orderService.VerifyOrderPayment(reference)
	if err != nil {
		resp.Status = http.StatusBadRequest
		resp.Message = err.Error()

		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	resp.Status = http.StatusOK
	resp.Message = "Success"

	return c.Status(http.StatusOK).JSON(resp)
}

func (h *orderHandler) StatusHistoryByOrderId(c *fiber.Ctx) error {
	var resp response.Response
	var orderStatusHistoriesResp []response.OrderStatusHistoryResponse

	orderId := c.Params("order_id")
	orderStatusHistories, err := h.orderStatusHistoryService.FindOrderStatusHistoriesByOrderId(orderId)
	if err != nil {
		resp.Status = fiber.StatusBadRequest
		resp.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	for _, orderStatusHistory := range orderStatusHistories {
		orderStatusHistoriesResp = append(orderStatusHistoriesResp, ConvertOrderStatusHistoryDTOToResponse(orderStatusHistory))
	}

	resp.Status = fiber.StatusOK
	resp.Message = "Success"
	resp.Data = map[string]interface{}{"results": orderStatusHistoriesResp}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h *orderHandler) GetOrderStatuses(c *fiber.Ctx) error {
	var resp response.Response
	var orderStatusesResp []response.OrderStatusResponse

	orderStatuses, err := h.orderStatusService.GetOrderStatuses()
	if err != nil {
		resp.Status = fiber.StatusBadRequest
		resp.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	for _, orderStatus := range orderStatuses {
		orderStatusesResp = append(orderStatusesResp, ConvertOrderStatusDTOToResponse(orderStatus))
	}

	resp.Status = fiber.StatusOK
	resp.Message = "Success"
	resp.Data = map[string]interface{}{"results": orderStatusesResp}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h *orderHandler) OrderProcessing(c *fiber.Ctx) error {
	var resp response.Response

	id := c.Params("order_id")

	orderId, err := uuid.Parse(id)

	if err != nil {
		resp.Status = constants.InvalidOrderID
		resp.Message = "Invalid order ID"

		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	err = h.orderService.ProcessOrder(orderId)
	if err != nil {
		resp.Status = fiber.StatusBadRequest
		resp.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	resp.Status = fiber.StatusOK
	resp.Message = "Success"

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h *orderHandler) OutForDelivery(c *fiber.Ctx) error {
	var resp response.Response

	id := c.Params("order_id")

	orderId, err := uuid.Parse(id)

	if err != nil {
		resp.Status = constants.InvalidOrderID
		resp.Message = "Invalid order ID"

		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	err = h.orderService.OutForDelivery(orderId)
	if err != nil {
		resp.Status = fiber.StatusBadRequest
		resp.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	resp.Status = fiber.StatusOK
	resp.Message = "Success"

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h *orderHandler) Delivered(c *fiber.Ctx) error {
	var resp response.Response

	id := c.Params("order_id")

	orderId, err := uuid.Parse(id)

	if err != nil {
		resp.Status = constants.InvalidOrderID
		resp.Message = "Invalid order ID"

		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	err = h.orderService.Delivered(orderId)
	if err != nil {
		resp.Status = fiber.StatusBadRequest
		resp.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(resp)
	}

	resp.Status = fiber.StatusOK
	resp.Message = "Success"

	return c.Status(fiber.StatusOK).JSON(resp)
}
