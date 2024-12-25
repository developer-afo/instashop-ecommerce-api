package dto

import (
	"github.com/google/uuid"
)

type TransactionDTO struct {
	DTO

	UserID      uuid.UUID `json:"user_id"`
	Amount      float64   `json:"amount"`
	Type        string    `json:"type"`
	Reference   string    `json:"reference"`
	Description string    `json:"description"`
	ShortDesc   string    `json:"short_desc"`
	Status      string    `json:"status"`
	Method      string    `json:"method"`
	Vendor      string    `json:"vendor"`
}
