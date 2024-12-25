package order_repository

import (
	"strings"

	"github.com/google/uuid"

	"github.com/developer-afo/instashop-ecommerce-api/lib/database"
	"github.com/developer-afo/instashop-ecommerce-api/models"
	"github.com/developer-afo/instashop-ecommerce-api/repository"
)

type OrderPageable struct {
	repository.Pageable

	Status   string
	UserID   uuid.UUID
	FromDate string
	ToDate   string
}

type OrderRepositoryInterface interface {
	CreateOrder(order models.Order) (models.Order, error)
	FindOrderById(uuid uuid.UUID) (models.Order, error)
	FindOrderByReference(reference string) (models.Order, error)
	FindOrderByTransactionId(transactionId uuid.UUID) (models.Order, error)
	FindAllOrders(pageable OrderPageable) ([]models.Order, repository.Pagination, error)
	CheckOrderExistByCouponId(couponId uuid.UUID) (bool, error)
	UpdateOrder(order models.Order) (models.Order, error)
	DeleteOrder(uuid uuid.UUID) error
}

type orderRepository struct {
	database database.DatabaseInterface
}

func NewOrderRepository(database database.DatabaseInterface) OrderRepositoryInterface {
	return &orderRepository{database: database}
}

// CreateOrder implements OrderRepositoryInterface.
func (o *orderRepository) CreateOrder(order models.Order) (models.Order, error) {
	order.Prepare()

	err := o.database.Connection().Create(&order).Error

	return order, err
}

// FindOrderById implements OrderRepositoryInterface.
func (o *orderRepository) FindOrderById(uuid uuid.UUID) (order models.Order, err error) {

	err = o.database.Connection().
		Model(&models.Order{}).
		Preload("User").
		Preload("Status").
		Preload("Transaction").
		Preload("OrderItems").
		Preload("OrderItems.Product").
		Where("id = ?", uuid).
		First(&order).Error

	return order, err
}

// FindOrderByReference implements OrderRepositoryInterface.
func (o *orderRepository) FindOrderByReference(reference string) (order models.Order, err error) {

	err = o.database.Connection().
		Model(&models.Order{}).
		Preload("Status").
		Preload("Transaction").
		Preload("OrderItems").
		Where("reference = ?", reference).
		First(&order).Error

	return order, err
}

// FindOrderByTransactionId implements OrderRepositoryInterface.
func (o *orderRepository) FindOrderByTransactionId(transactionId uuid.UUID) (order models.Order, err error) {

	err = o.database.Connection().Model(&models.Order{}).Where("transaction_id = ?", transactionId).First(&order).Error

	return order, err
}

// FindAllOrders implements OrderRepositoryInterface.
func (o *orderRepository) FindAllOrders(pageable OrderPageable) ([]models.Order, repository.Pagination, error) {
	var orders []models.Order
	var order models.Order
	var pagination repository.Pagination
	var errCount error

	pagination.CurrentPage = int64(pageable.Page)
	pagination.TotalItems = 0
	pagination.TotalPages = 1

	offset := (pageable.Page - 1) * pageable.Size
	model := o.database.Connection().
		Model(&order).
		Preload("User").
		Preload("Status").
		Preload("StatusHistory").
		Preload("StatusHistory.Status").
		Preload("Transaction").
		Preload("OrderItems").
		Preload("OrderItems.Product").
		Preload("OrderItems.Product.Images")

	// Apply search filters
	if len(strings.TrimSpace(pageable.Search)) > 0 {
		model = model.Where("orders.reference LIKE ?", "%"+strings.TrimSpace(pageable.Search)+"%")
	}

	// if pageable.Status is ongoing then check for 4 statuses which are order_placed, awaiting_confirmation, order_processing, out_for_delivery
	if pageable.Status == "ongoing" {
		model = model.Joins("JOIN order_statuses ON orders.status_id = order_statuses.id").
			Where("order_statuses.short_name IN ('order_placed', 'awaiting_confirmation', 'order_processing', 'out_for_delivery')")
	} else if pageable.Status != "" {
		model = model.Joins("JOIN order_statuses ON orders.status_id = order_statuses.id").
			Where("order_statuses.short_name = ?", pageable.Status)
	}

	if pageable.UserID != uuid.Nil {
		model = model.Where("orders.user_id = ?", pageable.UserID)
	}

	// Count total items for pagination
	if errCount = model.Count(&pagination.TotalItems).Error; errCount != nil {
		return nil, pagination, errCount
	}

	// Apply pagination
	paginatedQuery := model.Offset(int(offset)).Limit(int(pageable.Size)).Order(pageable.SortBy + " " + pageable.SortDirection)

	// Execute the query
	if err := paginatedQuery.Find(&orders).Error; err != nil {
		return nil, pagination, err
	}

	// Calculate total pages
	if pagination.TotalItems > 0 {
		pagination.TotalPages = (pagination.TotalItems + int64(pageable.Size) - 1) / int64(pageable.Size)
	} else {
		pagination.TotalPages = 1
	}

	return orders, pagination, nil
}

// CheckOrderExistByCouponId implements OrderRepositoryInterface.
func (o *orderRepository) CheckOrderExistByCouponId(couponId uuid.UUID) (bool, error) {

	var count int64

	err := o.database.Connection().
		Model(&models.Order{}).
		Where("coupon_id = ?", couponId).
		Count(&count).Error

	return count > 0, err
}

// UpdateOrder implements OrderRepositoryInterface.
func (o *orderRepository) UpdateOrder(order models.Order) (models.Order, error) {

	err := o.database.Connection().
		Model(&models.Order{}).
		Where("id = ?", order.ID).
		Updates(order).Error

	return order, err
}

// DeleteOrder implements OrderRepositoryInterface.
func (o *orderRepository) DeleteOrder(uuid uuid.UUID) error {

	err := o.database.Connection().Delete(&models.Order{}, uuid).Error

	return err
}
