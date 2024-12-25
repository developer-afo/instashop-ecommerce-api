package payment_gateway_dto

type PaymentInitializationDTO struct {
	Amount    float64 `json:"amount"`
	Email     string  `json:"email"`
	Reference string  `json:"reference"`
	Gateway   string  `json:"gateway"`
}

type PaymentInitializationResponseDTO struct {
	Status     bool   `json:"status"`
	PaymentURL string `json:"payment_url"`
}

type PaymentVerifyResponseDTO struct {
	Status        bool   `json:"status"`
	PaymentStatus string `json:"payment_status"`
	Message       string `json:"message"`
}
