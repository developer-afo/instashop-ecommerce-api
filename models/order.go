package models

import (
	"github.com/google/uuid"

	"github.com/developer-afo/instashop-ecommerce-api/lib/database"
)

type Order struct {
	database.BaseModel

	UserID        uuid.UUID `json:"user_id"`
	TransactionID uuid.UUID `json:"transaction_id"`
	PaymentMethod string    `json:"payment_method"`
	Reference     string    `json:"reference"`
	TotalPrice    float64   `json:"total_price"`
	StatusID      uuid.UUID `json:"status_id"`

	User          User                 `json:"user" gorm:"foreignKey:UserID;references:ID"`
	OrderItems    []OrderItem          `json:"order_items" gorm:"foreignKey:OrderID;references:ID"`
	Status        OrderStatus          `json:"status" gorm:"foreignKey:StatusID;references:ID"`
	Transaction   Transaction          `json:"transaction" gorm:"foreignKey:TransactionID;references:ID"`
	StatusHistory []OrderStatusHistory `json:"status_history" gorm:"foreignKey:OrderID;references:ID"`
}

type OrderItem struct {
	database.BaseModel

	OrderID   uuid.UUID `json:"order_id"`
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
	Price     float64   `json:"price"`

	Product Product `json:"product" gorm:"foreignKey:ProductID;references:ID"`
}

type OrderStatus struct {
	database.BaseModel

	Name      string `json:"name"`
	ShortName string `json:"short_name"`
}

type OrderStatusHistory struct {
	database.BaseModel

	OrderID  uuid.UUID `json:"order_id"`
	StatusID uuid.UUID `json:"status_id"`

	Status OrderStatus `json:"status" gorm:"foreignKey:StatusID;references:ID"`
}
