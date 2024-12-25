package order_repository

import (
	"github.com/google/uuid"

	"github.com/developer-afo/instashop-ecommerce-api/lib/database"
	"github.com/developer-afo/instashop-ecommerce-api/models"
)

type OrderStatusHistoryRepositoryInterface interface {
	CreateOrderStatusHistory(orderStatusHistory models.OrderStatusHistory) (models.OrderStatusHistory, error)
	FindOrderStatusHistoryById(uuid uuid.UUID) (models.OrderStatusHistory, error)
	FindOrderStatusHistoriesByOrderId(orderId uuid.UUID) ([]models.OrderStatusHistory, error)
	UpdateOrderStatusHistory(orderStatusHistory models.OrderStatusHistory) (models.OrderStatusHistory, error)
}

type orderStatusHistoryRepository struct {
	database database.DatabaseInterface
}

func NewOrderStatusHistoryRepository(database database.DatabaseInterface) OrderStatusHistoryRepositoryInterface {
	return &orderStatusHistoryRepository{database: database}
}

// CreateOrderStatusHistory implements OrderStatusHistoryRepositoryInterface.
func (o *orderStatusHistoryRepository) CreateOrderStatusHistory(orderStatusHistory models.OrderStatusHistory) (models.OrderStatusHistory, error) {
	orderStatusHistory.Prepare()

	err := o.database.Connection().Create(&orderStatusHistory).Error

	return orderStatusHistory, err
}

// FindOrderStatusHistoryById implements OrderStatusHistoryRepositoryInterface.
func (o *orderStatusHistoryRepository) FindOrderStatusHistoryById(uuid uuid.UUID) (orderStatusHistory models.OrderStatusHistory, err error) {

	err = o.database.Connection().Model(&models.OrderStatusHistory{}).Where("id = ?", uuid).First(&orderStatusHistory).Error

	return orderStatusHistory, err
}

// FindOrderStatusHistoriesByOrderId implements OrderStatusHistoryRepositoryInterface.
func (o *orderStatusHistoryRepository) FindOrderStatusHistoriesByOrderId(orderId uuid.UUID) (orderStatusHistories []models.OrderStatusHistory, err error) {

	err = o.database.Connection().
		Model(&models.OrderStatusHistory{}).
		Where("order_id = ?", orderId).
		Preload("Status").
		Find(&orderStatusHistories).Error

	return orderStatusHistories, err
}

// UpdateOrderStatusHistory implements OrderStatusHistoryRepositoryInterface.
func (o *orderStatusHistoryRepository) UpdateOrderStatusHistory(orderStatusHistory models.OrderStatusHistory) (models.OrderStatusHistory, error) {

	err := o.database.Connection().Save(&orderStatusHistory).Error

	return orderStatusHistory, err
}
