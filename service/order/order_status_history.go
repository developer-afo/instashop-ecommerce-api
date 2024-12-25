package order_service

import (
	"github.com/google/uuid"

	"github.com/developer-afo/instashop-ecommerce-api/dto"
	"github.com/developer-afo/instashop-ecommerce-api/models"
	order_repository "github.com/developer-afo/instashop-ecommerce-api/repository/order"
)

type OrderStatusHistoryServiceInterface interface {
	CreateOrderStatusHistory(orderStatusHistory dto.OrderStatusHistoryDTO) (dto.OrderStatusHistoryDTO, error)
	FindOrderStatusHistoryById(uuid uuid.UUID) (dto.OrderStatusHistoryDTO, error)
	FindOrderStatusHistoriesByOrderId(orderId string) ([]dto.OrderStatusHistoryDTO, error)
	ConvertToDTO(orderStatusHistory models.OrderStatusHistory) dto.OrderStatusHistoryDTO
}

type orderStatusHistoryService struct {
	orderStatusHistoryRepository order_repository.OrderStatusHistoryRepositoryInterface
	orderStatusService           OrderStatusServiceInterface
}

func NewOrderStatusHistoryService(
	orderStatusHistoryRepository order_repository.OrderStatusHistoryRepositoryInterface,
	orderStatusService OrderStatusServiceInterface,
) OrderStatusHistoryServiceInterface {
	return &orderStatusHistoryService{
		orderStatusHistoryRepository: orderStatusHistoryRepository,
		orderStatusService:           orderStatusService,
	}
}

func (s *orderStatusHistoryService) ConvertToDTO(orderStatusHistory models.OrderStatusHistory) (orderStatusHistoryDto dto.OrderStatusHistoryDTO) {

	orderStatusHistoryDto.ID = orderStatusHistory.ID
	orderStatusHistoryDto.OrderUUID = orderStatusHistory.OrderID
	orderStatusHistoryDto.StatusUUID = orderStatusHistory.StatusID
	orderStatusHistoryDto.Status = s.orderStatusService.ConvertToDTO(orderStatusHistory.Status)
	orderStatusHistoryDto.CreatedAt = orderStatusHistory.CreatedAt
	orderStatusHistoryDto.UpdatedAt = orderStatusHistory.UpdatedAt
	orderStatusHistoryDto.DeletedAt = orderStatusHistory.DeletedAt.Time

	return orderStatusHistoryDto
}

func (s *orderStatusHistoryService) ConvertToModel(orderStatusHistoryDto dto.OrderStatusHistoryDTO) (orderStatusHistory models.OrderStatusHistory) {

	orderStatusHistory.ID = orderStatusHistoryDto.ID
	orderStatusHistory.OrderID = orderStatusHistoryDto.OrderUUID
	orderStatusHistory.StatusID = orderStatusHistoryDto.StatusUUID
	orderStatusHistory.CreatedAt = orderStatusHistoryDto.CreatedAt
	orderStatusHistory.UpdatedAt = orderStatusHistoryDto.UpdatedAt
	orderStatusHistory.DeletedAt.Time = orderStatusHistoryDto.DeletedAt

	return orderStatusHistory
}

// CreateOrderStatusHistory implements OrderStatusHistoryServiceInterface.
func (s *orderStatusHistoryService) CreateOrderStatusHistory(orderStatusHistoryDtoArg dto.OrderStatusHistoryDTO) (dto.OrderStatusHistoryDTO, error) {

	orderStatusHistoryModel := s.ConvertToModel(orderStatusHistoryDtoArg)
	orderStatusHistoryModel, err := s.orderStatusHistoryRepository.CreateOrderStatusHistory(orderStatusHistoryModel)
	if err != nil {
		return dto.OrderStatusHistoryDTO{}, err
	}

	return s.ConvertToDTO(orderStatusHistoryModel), nil
}

// FindOrderStatusHistoryById implements OrderStatusHistoryServiceInterface.
func (s *orderStatusHistoryService) FindOrderStatusHistoryById(uuid uuid.UUID) (dto.OrderStatusHistoryDTO, error) {

	orderStatusHistoryModel, err := s.orderStatusHistoryRepository.FindOrderStatusHistoryById(uuid)
	if err != nil {
		return dto.OrderStatusHistoryDTO{}, err
	}

	return s.ConvertToDTO(orderStatusHistoryModel), nil
}

// FindOrderStatusHistoriesByOrderId implements OrderStatusHistoryServiceInterface.
func (s *orderStatusHistoryService) FindOrderStatusHistoriesByOrderId(orderId string) ([]dto.OrderStatusHistoryDTO, error) {

	orderStatusHistories, err := s.orderStatusHistoryRepository.FindOrderStatusHistoriesByOrderId(uuid.MustParse(orderId))
	if err != nil {
		return []dto.OrderStatusHistoryDTO{}, err
	}

	var orderStatusHistoriesDto []dto.OrderStatusHistoryDTO
	for _, orderStatusHistory := range orderStatusHistories {
		orderStatusHistoriesDto = append(orderStatusHistoriesDto, s.ConvertToDTO(orderStatusHistory))
	}

	return orderStatusHistoriesDto, nil
}
