package order_repository

import (
	"github.com/google/uuid"

	"github.com/developer-afo/instashop-ecommerce-api/lib/database"
	"github.com/developer-afo/instashop-ecommerce-api/models"
)

type OrderItemRepositoryInterface interface {
	CreateOrderItem(orderItem models.OrderItem) (models.OrderItem, error)
	BatchCreateOrderItem(items []models.OrderItem) error
	FindOrderItemById(uuid uuid.UUID) (models.OrderItem, error)
	FindOrderItemsByOrderId(orderId uuid.UUID) ([]models.OrderItem, error)
	UpdateOrderItem(orderItem models.OrderItem) (models.OrderItem, error)
}

type orderItemRepository struct {
	database database.DatabaseInterface
}

func NewOrderItemRepository(database database.DatabaseInterface) OrderItemRepositoryInterface {
	return &orderItemRepository{database: database}
}

// CreateOrderItem implements OrderItemRepositoryInterface.
func (o *orderItemRepository) CreateOrderItem(orderItem models.OrderItem) (models.OrderItem, error) {
	orderItem.Prepare()

	err := o.database.Connection().Create(&orderItem).Error

	return orderItem, err
}

// BatchCreateOrderItem implements OrderItemRepositoryInterface.
func (o *orderItemRepository) BatchCreateOrderItem(items []models.OrderItem) error {

	err := o.database.Connection().Create(&items).Error

	return err
}

// FindOrderItemById implements OrderItemRepositoryInterface.
func (o *orderItemRepository) FindOrderItemById(uuid uuid.UUID) (orderItem models.OrderItem, err error) {

	err = o.database.Connection().Model(&models.OrderItem{}).Where("id = ?", uuid).First(&orderItem).Error

	return orderItem, err
}

// FindOrderItemsByOrderId implements OrderItemRepositoryInterface.
func (o *orderItemRepository) FindOrderItemsByOrderId(orderId uuid.UUID) (orderItems []models.OrderItem, err error) {

	err = o.database.Connection().Model(&models.OrderItem{}).Where("order_id = ?", orderId).Find(&orderItems).Error

	return orderItems, err
}

// UpdateOrderItem implements OrderItemRepositoryInterface.
func (o *orderItemRepository) UpdateOrderItem(orderItem models.OrderItem) (models.OrderItem, error) {

	err := o.database.Connection().Save(&orderItem).Error

	return orderItem, err
}
