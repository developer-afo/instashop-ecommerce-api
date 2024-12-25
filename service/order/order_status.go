package order_service

import (
	"github.com/google/uuid"

	"github.com/developer-afo/instashop-ecommerce-api/dto"
	"github.com/developer-afo/instashop-ecommerce-api/models"
	order_repository "github.com/developer-afo/instashop-ecommerce-api/repository/order"
)

var (
	ONGOING               = "ongoing"
	ORDER_PLACED          = "order_placed"
	AWAITING_CONFIRMATION = "awaiting_confirmation"
	ORDER_PROCESSING      = "order_processing"
	OUT_FOR_DELIVERY      = "out_for_delivery"
	DELIVERED             = "delivered"
	CANCELLED             = "cancelled"
)

type OrderStatusServiceInterface interface {
	CreateOrderStatus(orderStatus models.OrderStatus) (dto.OrderStatusDTO, error)
	GetOrderStatuses() ([]dto.OrderStatusDTO, error)
	FindOrderStatusById(id uuid.UUID) (dto.OrderStatusDTO, error)
	StatusOrderPlaced() (dto.OrderStatusDTO, error)
	StatusAwaitingConfirmation() (dto.OrderStatusDTO, error)
	StatusOrderProcessing() (dto.OrderStatusDTO, error)
	StatusOutForDelivery() (dto.OrderStatusDTO, error)
	StatusDelivered() (dto.OrderStatusDTO, error)
	StatusCancelled() (dto.OrderStatusDTO, error)
	ConvertToDTO(orderStatus models.OrderStatus) dto.OrderStatusDTO
}

type orderStatusService struct {
	orderStatusRepository order_repository.OrderStatusRepositoryInterface
}

func NewOrderStatusService(orderStatusRepository order_repository.OrderStatusRepositoryInterface) OrderStatusServiceInterface {
	return &orderStatusService{orderStatusRepository: orderStatusRepository}
}

func (o *orderStatusService) ConvertToDTO(orderStatus models.OrderStatus) dto.OrderStatusDTO {
	var orderStatusDTO dto.OrderStatusDTO

	orderStatusDTO.ID = orderStatus.ID
	orderStatusDTO.Name = orderStatus.Name
	orderStatusDTO.ShortName = orderStatus.ShortName

	return orderStatusDTO
}

// convert to model
func (o *orderStatusService) ConvertToModel(orderStatusDTO dto.OrderStatusDTO) models.OrderStatus {
	var orderStatus models.OrderStatus

	orderStatus.ID = orderStatusDTO.ID
	orderStatus.Name = orderStatusDTO.Name
	orderStatus.ShortName = orderStatusDTO.ShortName

	return orderStatus
}

// CreateOrderStatus implements OrderStatusServiceInterface.
func (o *orderStatusService) CreateOrderStatus(orderStatus models.OrderStatus) (dto.OrderStatusDTO, error) {
	orderStatus, err := o.orderStatusRepository.CreateOrderStatus(orderStatus)

	if err != nil {
		return dto.OrderStatusDTO{}, err
	}

	return o.ConvertToDTO(orderStatus), nil
}

// GetOrderStatuses implements OrderStatusServiceInterface.
func (o *orderStatusService) GetOrderStatuses() ([]dto.OrderStatusDTO, error) {
	orderStatuses, err := o.orderStatusRepository.GetOrderStatuses()

	if err != nil {
		return []dto.OrderStatusDTO{}, err
	}

	var orderStatusesDTO []dto.OrderStatusDTO

	for _, orderStatus := range orderStatuses {
		orderStatusesDTO = append(orderStatusesDTO, o.ConvertToDTO(orderStatus))
	}

	return orderStatusesDTO, nil
}

// FindOrderStatusById implements OrderStatusServiceInterface.
func (o *orderStatusService) FindOrderStatusById(id uuid.UUID) (dto.OrderStatusDTO, error) {
	orderStatus, err := o.orderStatusRepository.FindOrderStatusById(id)

	if err != nil {
		return dto.OrderStatusDTO{}, err
	}

	return o.ConvertToDTO(orderStatus), nil
}

// StatusOrderPlaced implements OrderStatusServiceInterface.
func (o *orderStatusService) StatusOrderPlaced() (dto.OrderStatusDTO, error) {
	orderStatus, err := o.orderStatusRepository.FindOrderStatusByShortName(ORDER_PLACED)

	if err != nil {
		return dto.OrderStatusDTO{}, err
	}

	return o.ConvertToDTO(orderStatus), nil
}

// StatusAwaitingConfirmation implements OrderStatusServiceInterface.
func (o *orderStatusService) StatusAwaitingConfirmation() (dto.OrderStatusDTO, error) {
	orderStatus, err := o.orderStatusRepository.FindOrderStatusByShortName(AWAITING_CONFIRMATION)

	if err != nil {
		return dto.OrderStatusDTO{}, err
	}

	return o.ConvertToDTO(orderStatus), nil
}

// StatusOrderProcessing implements OrderStatusServiceInterface.
func (o *orderStatusService) StatusOrderProcessing() (dto.OrderStatusDTO, error) {
	orderStatus, err := o.orderStatusRepository.FindOrderStatusByShortName(ORDER_PROCESSING)

	if err != nil {
		return dto.OrderStatusDTO{}, err
	}

	return o.ConvertToDTO(orderStatus), nil
}

// StatusOutForDelivery implements OrderStatusServiceInterface.
func (o *orderStatusService) StatusOutForDelivery() (dto.OrderStatusDTO, error) {
	orderStatus, err := o.orderStatusRepository.FindOrderStatusByShortName(OUT_FOR_DELIVERY)

	if err != nil {
		return dto.OrderStatusDTO{}, err
	}

	return o.ConvertToDTO(orderStatus), nil
}

// StatusDelivered implements OrderStatusServiceInterface.
func (o *orderStatusService) StatusDelivered() (dto.OrderStatusDTO, error) {
	orderStatus, err := o.orderStatusRepository.FindOrderStatusByShortName(DELIVERED)

	if err != nil {
		return dto.OrderStatusDTO{}, err
	}

	return o.ConvertToDTO(orderStatus), nil
}

// StatusCancelled implements OrderStatusServiceInterface.
func (o *orderStatusService) StatusCancelled() (dto.OrderStatusDTO, error) {
	orderStatus, err := o.orderStatusRepository.FindOrderStatusByShortName(CANCELLED)

	if err != nil {
		return dto.OrderStatusDTO{}, err
	}

	return o.ConvertToDTO(orderStatus), nil
}
