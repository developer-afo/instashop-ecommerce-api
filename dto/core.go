package dto

import "github.com/google/uuid"

type ProductDTO struct {
	DTO

	Slug          string  `json:"slug"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Specification string  `json:"specification"`
	Price         float64 `json:"price"`
	SlashPrice    float64 `json:"slash_price"`
	Stock         int     `json:"stock"`
	Sales         int     `json:"sales"`

	Images []ImageDTO `json:"images"`
}

type ImageDTO struct {
	DTO

	ProductUUID uuid.UUID `json:"product_id"`
	Key         string    `json:"key"`
}
