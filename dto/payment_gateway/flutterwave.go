package payment_gateway_dto

import "time"

// Currency       string         `json:"currency"`
type InitializeFlutterwaveRequest struct {
	TxRef          string         `json:"tx_ref"`
	Amount         float64        `json:"amount"`
	RedirectURL    string         `json:"redirect_url"`
	Customer       Customer       `json:"customer"`
	Customizations Customizations `json:"customizations"`
}

type Customer struct {
	Email string `json:"email"`
}

type Customizations struct {
	Title string `json:"title"`
}

type InitializeFlutterwaveResponse struct {
	Status  string `json:"status"` // success || error
	Message string `json:"message"`
	Data    struct {
		Link string `json:"link"`
	} `json:"data"`
}

type VerifyFlutterwaveResponse struct {
	Status  string `json:"status"` // success || error
	Message string `json:"message"`
	Data    struct {
		ID                int       `json:"id"`
		TxRef             string    `json:"tx_ref"`
		FlwRef            string    `json:"flw_ref"`
		DeviceFingerprint string    `json:"device_fingerprint"`
		Amount            float64   `json:"amount"`
		Currency          string    `json:"currency"`
		ChargedAmount     float64   `json:"charged_amount"`
		AppFee            float64   `json:"app_fee"`
		MerchantFee       float64   `json:"merchant_fee"`
		ProcessorResponse string    `json:"processor_response"`
		AuthModel         string    `json:"auth_model"`
		IP                string    `json:"ip"`
		Narration         string    `json:"narration"`
		Status            string    `json:"status"`
		PaymentType       string    `json:"payment_type"`
		CreatedAt         time.Time `json:"created_at"`
		AccountId         int       `json:"account_id"`
		Meta              struct {
			CheckoutInitAddress     string `json:"__CheckoutInitAddress"`
			ConsumerID              string `json:"consumer_id"`
			OriginatorAccountNumber string `json:"originatoraccountnumber"`
			OriginatorName          string `json:"originatorname"`
			BankName                string `json:"bankname"`
			OriginatorAmount        string `json:"originatoramount"`
		} `json:"meta"`
		AmountSettled float64 `json:"amount_settled"`
		Customer      struct {
			ID          int       `json:"id"`
			Name        string    `json:"name"`
			PhoneNumber string    `json:"phone_number"`
			Email       string    `json:"email"`
			CreatedAt   time.Time `json:"created_at"`
		} `json:"customer"`
	} `json:"data"`
}
