package request

type CreateProductRequest struct {
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Specification string   `json:"specification"`
	Price         int      `json:"price"`
	Stock         int      `json:"stock"`
	SlashPrice    int      `json:"slash_price"`
	Images        []string `json:"images"`
}

type UpdateProductRequest struct {
	CreateProductRequest
}

type ImageRequest struct {
	Key string `json:"key"`
}
