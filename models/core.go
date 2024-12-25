package models

import (
	"github.com/google/uuid"

	"github.com/developer-afo/instashop-ecommerce-api/lib/database"
)

type Product struct {
	database.BaseModel

	Slug          string  `json:"slug"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Specification string  `json:"specification"`
	Price         float64 `json:"price"`
	SlashPrice    float64 `json:"slash_price"`
	Stock         int     `json:"stock"`
	Brand         string  `json:"brand"`

	Sales  int     `json:"sales" gorm:"->"`
	Images []Image `json:"images" gorm:"foreignKey:ProductID;references:ID"`
}

type Image struct {
	database.BaseModel

	ProductID uuid.UUID `json:"product_id" gorm:"type:uuid"`
	Key       string    `json:"key"`
}
