package response

import (
	"time"

	"github.com/google/uuid"
)

type OrderResponse struct {
	ID            uuid.UUID                    `json:"id"`
	CreatedAt     time.Time                    `json:"created_at"`
	UpdatedAt     time.Time                    `json:"updated_at"`
	PaymentMethod string                       `json:"payment_method"`
	Reference     string                       `json:"reference"`
	TotalPrice    float64                      `json:"total_price"`
	OrderItems    []OrderItemResponse          `json:"order_items"`
	Status        OrderStatusResponse          `json:"status"`
	StatusHistory []OrderStatusHistoryResponse `json:"status_history"`
	Transaction   TransactionResponse          `json:"transaction"`
	User          UserResponseData             `json:"user"`
}

type OrderItemResponse struct {
	Product  ProductResponse `json:"product"`
	Quantity int             `json:"quantity"`
	Price    float64         `json:"price"`
}

type OrderStatusResponse struct {
	Name      string `json:"name"`
	ShortName string `json:"short_name"`
}

type OrderStatusHistoryResponse struct {
	Status    OrderStatusResponse `json:"status"`
	CreatedAt time.Time           `json:"created_at"`
}
