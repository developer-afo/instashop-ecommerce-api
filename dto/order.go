package dto

import "github.com/google/uuid"

type OrderDTO struct {
	DTO

	UserID        uuid.UUID  `json:"user_id"`
	TransactionID uuid.UUID  `json:"transaction_id"`
	CouponID      *uuid.UUID `json:"coupon_id"`
	PaymentMethod string     `json:"payment_method"`
	Reference     string     `json:"reference"`
	TotalPrice    float64    `json:"total_price"`
	StatusUUID    uuid.UUID  `json:"status_id"`

	User          UserDTO                 `json:"user"`
	Items         []OrderItemDTO          `json:"items"`
	Transaction   TransactionDTO          `json:"transaction"`
	Status        OrderStatusDTO          `json:"status"`
	StatusHistory []OrderStatusHistoryDTO `json:"status_history"`
}

type OrderItemDTO struct {
	DTO

	OrderUUID   uuid.UUID `json:"order_id"`
	ProductUUID uuid.UUID `json:"product_id"`
	Quantity    int       `json:"quantity"`
	Price       float64   `json:"price"`

	Product ProductDTO `json:"product"`
}

type OrderStatusDTO struct {
	DTO

	Name      string `json:"name"`
	ShortName string `json:"short_name"`
}
type OrderStatusHistoryDTO struct {
	DTO

	OrderUUID  uuid.UUID `json:"order_id"`
	StatusUUID uuid.UUID `json:"status_id"`

	Status OrderStatusDTO `json:"status"`
}

type CreateOrderDTO struct {
	UserID            uuid.UUID            `json:"user_id"`
	ShippingAddressID uuid.UUID            `json:"shipping_address_id"`
	CouponID          *uuid.UUID           `json:"coupon_id"`
	ShippingTypeID    uuid.UUID            `json:"shipping_type_id"`
	PaymentMethod     string               `json:"payment_method"`
	Items             []CreateOrderItemDTO `json:"items"`
}

type CreateOrderItemDTO struct {
	ProductUUID string `json:"product_id"`
	Quantity    int    `json:"quantity"`
}
