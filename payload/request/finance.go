package request

import "time"

type PaystackWebhook struct {
	Event string `json:"event"`
	Data  struct {
		ID              int         `json:"id"`
		Domain          string      `json:"domain"`
		Status          string      `json:"status"`
		Reference       string      `json:"reference"`
		Amount          int         `json:"amount"`
		Message         interface{} `json:"message"`
		GatewayResponse string      `json:"gateway_response"`
		PaidAt          time.Time   `json:"paid_at"`
		CreatedAt       time.Time   `json:"created_at"`
		Channel         string      `json:"channel"`
		Currency        string      `json:"currency"`
		IPAddress       string      `json:"ip_address"`
		Metadata        interface{} `json:"metadata"`
		Log             struct {
			TimeSpent      int           `json:"time_spent"`
			Attempts       int           `json:"attempts"`
			Authentication string        `json:"authentication"`
			Errors         int           `json:"errors"`
			Success        bool          `json:"success"`
			Mobile         bool          `json:"mobile"`
			Input          []interface{} `json:"input"`
			Channel        interface{}   `json:"channel"`
			History        []struct {
				Type    string `json:"type"`
				Message string `json:"message"`
				Time    int    `json:"time"`
			} `json:"history"`
		} `json:"log"`
		Fees     interface{} `json:"fees"`
		Customer struct {
			ID           int         `json:"id"`
			FirstName    string      `json:"first_name"`
			LastName     string      `json:"last_name"`
			Email        string      `json:"email"`
			CustomerCode string      `json:"customer_code"`
			Phone        interface{} `json:"phone"`
			Metadata     interface{} `json:"metadata"`
			RiskAction   string      `json:"risk_action"`
		} `json:"customer"`
		Authorization struct {
			AuthorizationCode string `json:"authorization_code"`
			Bin               string `json:"bin"`
			Last4             string `json:"last4"`
			ExpMonth          string `json:"exp_month"`
			ExpYear           string `json:"exp_year"`
			CardType          string `json:"card_type"`
			Bank              string `json:"bank"`
			CountryCode       string `json:"country_code"`
			Brand             string `json:"brand"`
			AccountName       string `json:"account_name"`
		} `json:"authorization"`
		Plan struct{} `json:"plan"`
	} `json:"data"`
}

type FlutterwaveWebhook struct {
	Event string `json:"event"`
	Data  struct {
		ID                int         `json:"id"`
		TxRef             string      `json:"tx_ref"`
		FlwRef            string      `json:"flw_ref"`
		DeviceFingerprint string      `json:"device_fingerprint"`
		Amount            int         `json:"amount"`
		Currency          string      `json:"currency"`
		ChargedAmount     int         `json:"charged_amount"`
		AppFee            interface{} `json:"app_fee"`
		MerchantFee       int         `json:"merchant_fee"`
		ProcessorResponse string      `json:"processor_response"`
		AuthModel         string      `json:"auth_model"`
		IP                string      `json:"ip"`
		Narration         string      `json:"narration"`
		Status            string      `json:"status"`
		PaymentType       string      `json:"payment_type"`
		CreatedAt         time.Time   `json:"created_at"`
		AccountID         int         `json:"account_id"`
		Customer          struct {
			ID          int         `json:"id"`
			Name        string      `json:"name"`
			PhoneNumber interface{} `json:"phone_number"`
			Email       string      `json:"email"`
			CreatedAt   time.Time   `json:"created_at"`
		} `json:"customer"`
		Card struct {
			First6Digits string `json:"first_6digits"`
			Last4Digits  string `json:"last_4digits"`
			Issuer       string `json:"issuer"`
			Country      string `json:"country"`
			Type         string `json:"type"`
			Expiry       string `json:"expiry"`
		} `json:"card"`
	} `json:"data"`
}
