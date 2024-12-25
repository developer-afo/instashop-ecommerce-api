package order_service

import (
	"github.com/developer-afo/instashop-ecommerce-api/dto"
	"github.com/developer-afo/instashop-ecommerce-api/models"
	order_repository "github.com/developer-afo/instashop-ecommerce-api/repository/order"
	core_service "github.com/developer-afo/instashop-ecommerce-api/service/core"
	"github.com/google/uuid"
)

type OrderItemServiceInterface interface {
	BatchCreateOrderItem(orderUuid string, items []dto.OrderItemDTO) error
	ConvertToDTO(orderItem models.OrderItem) dto.OrderItemDTO
}

type orderItemService struct {
	orderItemRepository order_repository.OrderItemRepositoryInterface
	productService      core_service.ProductServiceInterface
}

func NewOrderItemService(
	orderItemRepository order_repository.OrderItemRepositoryInterface,
	productService core_service.ProductServiceInterface,
) OrderItemServiceInterface {
	return &orderItemService{
		orderItemRepository: orderItemRepository,
		productService:      productService,
	}
}

func (s *orderItemService) ConvertToDTO(orderItem models.OrderItem) dto.OrderItemDTO {
	var orderItemDTO dto.OrderItemDTO

	orderItemDTO.OrderUUID = orderItem.OrderID
	orderItemDTO.ProductUUID = orderItem.ProductID
	orderItemDTO.Quantity = orderItem.Quantity
	orderItemDTO.Price = orderItem.Price

	orderItemDTO.Product = s.productService.ConvertToDTO(orderItem.Product)

	return orderItemDTO
}

func (s *orderItemService) ConvertToModel(orderItemDTO dto.OrderItemDTO) models.OrderItem {
	var orderItem models.OrderItem

	orderItem.ID = orderItemDTO.ID
	orderItem.OrderID = orderItemDTO.OrderUUID
	orderItem.ProductID = orderItemDTO.ProductUUID
	orderItem.Quantity = orderItemDTO.Quantity
	orderItem.Price = orderItemDTO.Price

	return orderItem
}

func (s *orderItemService) BatchCreateOrderItem(orderUuid string, items []dto.OrderItemDTO) error {
	var orderItems []models.OrderItem

	orderId, err := uuid.Parse(orderUuid)

	if err != nil {
		return err
	}

	for _, item := range items {
		var orderItem models.OrderItem

		orderItem.ID, _ = uuid.NewV7()
		orderItem.OrderID = orderId
		orderItem.ProductID = item.ProductUUID
		orderItem.Quantity = item.Quantity
		orderItem.Price = item.Price

		orderItems = append(orderItems, orderItem)
	}

	err = s.orderItemRepository.BatchCreateOrderItem(orderItems)

	return err
}
