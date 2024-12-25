package payment_gateway_service

import (
	payment_gateway_dto "github.com/developer-afo/instashop-ecommerce-api/dto/payment_gateway"
	"github.com/developer-afo/instashop-ecommerce-api/lib/constants"
	"github.com/developer-afo/instashop-ecommerce-api/service"
)

var (
	FlutterwaveStatusAbandoned  = "abandoned"
	FlutterwaveStatusFailed     = "failed"
	FlutterwaveStatusOngoing    = "ongoing"
	FlutterwaveStatusPending    = "pending"
	FlutterwaveStatusProcessing = "processing"
	FlutterwaveStatusQueued     = "queued"
	FlutterwaveStatusReversed   = "reversed"
	FlutterwaveStatusSuccess    = "successful"
)

type FlutterwaveServiceInterface interface {
	InitializePayment(paymentDto payment_gateway_dto.InitializeFlutterwaveRequest) (payment_gateway_dto.InitializeFlutterwaveResponse, error)
	VerifyPayment(reference string) (payment_gateway_dto.VerifyFlutterwaveResponse, error)
}

type flutterwaveService struct {
	httpService service.HttpServiceInterface
	secretKey   string
	callbackURL string
	baseURL     string
}

func NewFlutterwaveService(httpService service.HttpServiceInterface, env constants.Env) FlutterwaveServiceInterface {
	return &flutterwaveService{
		httpService: httpService,
		secretKey:   env.FLUTTERWAVE_SECRET_KEY,
		callbackURL: env.PAYMENT_CALLBACK_URL,
		baseURL:     "https://api.flutterwave.com/v3",
	}
}

func (p *flutterwaveService) InitializePayment(paymentDto payment_gateway_dto.InitializeFlutterwaveRequest) (data payment_gateway_dto.InitializeFlutterwaveResponse, err error) {
	url := p.baseURL + "/payments"

	headers := map[string]string{
		"Authorization": "Bearer " + p.secretKey,
	}

	paymentDto.RedirectURL = p.callbackURL + paymentDto.TxRef
	paymentDto.Customizations.Title = "MaziMart"

	resp, err := p.httpService.Post(url, headers, paymentDto)

	if err != nil {
		return data, err
	}

	defer resp.Body.Close()

	if err = p.httpService.BodyToDTO(resp.Body, &data); err != nil {
		return data, err
	}

	return data, nil
}

func (p *flutterwaveService) VerifyPayment(reference string) (data payment_gateway_dto.VerifyFlutterwaveResponse, err error) {
	url := p.baseURL + "/transactions/verify_by_reference?tx_ref=" + reference

	headers := map[string]string{
		"Authorization": "Bearer " + p.secretKey,
	}

	resp, err := p.httpService.Get(url, headers)

	if err != nil {
		return payment_gateway_dto.VerifyFlutterwaveResponse{}, err
	}

	defer resp.Body.Close()

	if err = p.httpService.BodyToDTO(resp.Body, &data); err != nil {
		return data, err
	}

	return data, nil
}
