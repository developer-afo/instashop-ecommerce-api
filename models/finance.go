package models

import (
	"github.com/google/uuid"

	"github.com/developer-afo/instashop-ecommerce-api/lib/database"
)

type Transaction struct {
	database.BaseModel

	UserID      uuid.UUID `json:"user_id"`
	Amount      float64   `json:"amount"`
	Type        string    `json:"type"` // credit or debit
	Reference   string    `json:"reference"`
	Description string    `json:"description"` // use this to differentiate what the transaction is for
	ShortDesc   string    `json:"short_desc" gorm:"column:purpose"`
	Status      string    `json:"status"`
	Method      string    `json:"method"`
	Vendor      string    `json:"vendor"`
}
