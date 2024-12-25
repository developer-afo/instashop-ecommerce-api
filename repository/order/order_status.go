package order_repository

import (
	"github.com/google/uuid"

	"github.com/developer-afo/instashop-ecommerce-api/lib/database"
	"github.com/developer-afo/instashop-ecommerce-api/models"
)

type OrderStatusRepositoryInterface interface {
	GetOrderStatuses() ([]models.OrderStatus, error)
	CreateOrderStatus(orderStatus models.OrderStatus) (models.OrderStatus, error)
	FindOrderStatusById(uuid uuid.UUID) (models.OrderStatus, error)
	FindOrderStatusByShortName(shortName string) (models.OrderStatus, error)
	FindOrderStatuses() ([]models.OrderStatus, error)
	UpdateOrderStatus(orderStatus models.OrderStatus) (models.OrderStatus, error)
	DeleteOrderStatus(uuid uuid.UUID) error
}

type orderStatusRepository struct {
	database database.DatabaseInterface
}

func NewOrderStatusRepository(database database.DatabaseInterface) OrderStatusRepositoryInterface {
	return &orderStatusRepository{database: database}
}

// GetOrderStatuses implements OrderStatusRepositoryInterface.
func (o *orderStatusRepository) GetOrderStatuses() (orderStatuses []models.OrderStatus, err error) {

	err = o.database.Connection().Model(&models.OrderStatus{}).Find(&orderStatuses).Error

	return orderStatuses, err
}

// CreateOrderStatus implements OrderStatusRepositoryInterface.
func (o *orderStatusRepository) CreateOrderStatus(orderStatus models.OrderStatus) (models.OrderStatus, error) {
	orderStatus.Prepare()

	err := o.database.Connection().Create(&orderStatus).Error

	return orderStatus, err
}

// FindOrderStatusById implements OrderStatusRepositoryInterface.
func (o *orderStatusRepository) FindOrderStatusById(uuid uuid.UUID) (orderStatus models.OrderStatus, err error) {

	err = o.database.Connection().Model(&models.OrderStatus{}).Where("id = ?", uuid).First(&orderStatus).Error

	return orderStatus, err
}

// FindOrderStatusByShortName implements OrderStatusRepositoryInterface.
func (o *orderStatusRepository) FindOrderStatusByShortName(shortName string) (orderStatus models.OrderStatus, err error) {

	err = o.database.Connection().Model(&models.OrderStatus{}).Where("short_name = ?", shortName).First(&orderStatus).Error

	return orderStatus, err
}

// FindOrderStatuses implements OrderStatusRepositoryInterface.
func (o *orderStatusRepository) FindOrderStatuses() (orderStatuses []models.OrderStatus, err error) {

	err = o.database.Connection().Model(&models.OrderStatus{}).Find(&orderStatuses).Error

	return orderStatuses, err
}

// UpdateOrderStatus implements OrderStatusRepositoryInterface.
func (o *orderStatusRepository) UpdateOrderStatus(orderStatus models.OrderStatus) (models.OrderStatus, error) {

	err := o.database.Connection().Save(&orderStatus).Error

	return orderStatus, err
}

// DeleteOrderStatus implements OrderStatusRepositoryInterface.
func (o *orderStatusRepository) DeleteOrderStatus(uuid uuid.UUID) error {
	err := o.database.Connection().Where("id = ?", uuid).Delete(&models.OrderStatus{}).Error

	return err
}
