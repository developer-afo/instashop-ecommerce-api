package response

import (
	"time"

	"github.com/google/uuid"
)

type TransactionResponse struct {
	ID          uuid.UUID `json:"id"`
	Amount      float64   `json:"amount"`
	Type        string    `json:"type"`
	Reference   string    `json:"reference"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Method      string    `json:"method"`
	Vendor      string    `json:"vendor"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
