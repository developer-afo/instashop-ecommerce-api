package response

import (
	"time"

	"github.com/google/uuid"
)

type ProductResponse struct {
	UUID          uuid.UUID       `json:"id"`
	Slug          string          `json:"slug"`
	Name          string          `json:"name"`
	Description   string          `json:"description"`
	Specification string          `json:"specification"`
	Price         float64         `json:"price"`
	SlashPrice    float64         `json:"slash_price"`
	Stock         int             `json:"stock"`
	Sales         int             `json:"sales"`
	Images        []ImageResponse `json:"images"`
	CreatedAt     time.Time       `json:"created_at"`
}

type ImageResponse struct {
	Key string `json:"key"`
}
