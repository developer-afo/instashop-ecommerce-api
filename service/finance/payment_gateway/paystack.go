package payment_gateway_service

import (
	payment_gateway_dto "github.com/developer-afo/instashop-ecommerce-api/dto/payment_gateway"
	"github.com/developer-afo/instashop-ecommerce-api/lib/constants"
	"github.com/developer-afo/instashop-ecommerce-api/service"
)

var (
	PaystackStatusAbandoned  = "abandoned"
	PaystackStatusFailed     = "failed"
	PaystackStatusOngoing    = "ongoing"
	PaystackStatusPending    = "pending"
	PaystackStatusProcessing = "processing"
	PaystackStatusQueued     = "queued"
	PaystackStatusReversed   = "reversed"
	PaystackStatusSuccess    = "success"
)

type PaystackServiceInterface interface {
	InitializePayment(paymentDto payment_gateway_dto.Paystack) (payment_gateway_dto.InitializePaystackResponse, error)
	VerifyPayment(reference string) (payment_gateway_dto.VerifyPaystackResponse, error)
}

type paystackService struct {
	httpService service.HttpServiceInterface
	secretKey   string
	callbackURL string
	baseURL     string
}

func NewPaystackService(httpService service.HttpServiceInterface, env constants.Env) PaystackServiceInterface {
	return &paystackService{
		httpService: httpService,
		secretKey:   env.PAYSTACK_SECRET_KEY,
		callbackURL: env.PAYMENT_CALLBACK_URL,
		baseURL:     "https://api.paystack.co/transaction",
	}
}

func (p *paystackService) InitializePayment(paymentDto payment_gateway_dto.Paystack) (data payment_gateway_dto.InitializePaystackResponse, err error) {
	initializeURL := p.baseURL + "/initialize"

	headers := map[string]string{
		"Authorization": "Bearer " + p.secretKey,
	}

	body := map[string]interface{}{
		"amount":       paymentDto.Amount * 100,
		"email":        paymentDto.Email,
		"reference":    paymentDto.Reference,
		"callback_url": p.callbackURL + paymentDto.Reference,
	}

	resp, err := p.httpService.Post(initializeURL, headers, body)

	if err != nil {
		return data, err
	}

	defer resp.Body.Close()

	if err = p.httpService.BodyToDTO(resp.Body, &data); err != nil {
		return data, err
	}

	return data, nil
}

func (p *paystackService) VerifyPayment(reference string) (payment_gateway_dto.VerifyPaystackResponse, error) {

	verifyURL := p.baseURL + "/verify/" + reference

	headers := map[string]string{
		"Authorization": "Bearer " + p.secretKey,
	}

	resp, err := p.httpService.Get(verifyURL, headers)

	if err != nil {
		return payment_gateway_dto.VerifyPaystackResponse{}, err
	}

	defer resp.Body.Close()

	var data payment_gateway_dto.VerifyPaystackResponse

	if err = p.httpService.BodyToDTO(resp.Body, &data); err != nil {
		return data, err
	}

	return data, nil
}
