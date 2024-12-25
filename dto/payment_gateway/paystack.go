package payment_gateway_dto

import "time"

type Paystack struct {
	Amount    float64 `json:"amount"`
	Email     string  `json:"email"`
	Reference string  `json:"reference"`
}

type InitializePaystackResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		AuthorizationURL string `json:"authorization_url"`
		AccessCode       string `json:"access_code"`
		Reference        string `json:"reference"`
	}
}

type VerifyPaystackResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Status          string    `json:"status"`
		Message         string    `json:"message"`
		ID              int       `json:"id"`
		Domain          string    `json:"domain"`
		Reference       string    `json:"reference"`
		Amount          int       `json:"amount"`
		GatewayResponse string    `json:"gateway_response"`
		CreatedAt       time.Time `json:"created_at"`
		Channel         string    `json:"channel"`
		Currency        string    `json:"currency"`
		IPAddress       string    `json:"ip_address"`
		Customer        struct {
			ID        int     `json:"id"`
			FirstName *string `json:"first_name"`
			LastName  *string `json:"last_name"`
			Email     string  `json:"email"`
		} `json:"customer"`
		TransactionDate time.Time `json:"transaction_date"`
	}
}
