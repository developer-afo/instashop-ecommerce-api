package payment_gateway_service

import (
	"errors"
	"fmt"
	"time"

	payment_gateway_dto "github.com/developer-afo/instashop-ecommerce-api/dto/payment_gateway"
	finance_service "github.com/developer-afo/instashop-ecommerce-api/service/finance"
)

var (
	PaystackPaymentGateway    = "paystack"
	FlutterwavePaymentGateway = "flutterwave"
	ErrPaymentInitialization  = errors.New("failed to initialize payment")
)

type PaymentGatewayServiceInterface interface {
	InitializePayment(paymentDto payment_gateway_dto.PaymentInitializationDTO) (payment_gateway_dto.PaymentInitializationResponseDTO, error)
	VerifyPayment(reference string, gateway string) (payment_gateway_dto.PaymentVerifyResponseDTO, error)
}

type paymentGatewayService struct {
	paystackService    PaystackServiceInterface
	flutterwaveService FlutterwaveServiceInterface
}

func NewPaymentGatewayService(
	paystackService PaystackServiceInterface,
	flutterwaveService FlutterwaveServiceInterface,
) PaymentGatewayServiceInterface {
	return &paymentGatewayService{
		paystackService:    paystackService,
		flutterwaveService: flutterwaveService,
	}
}

func (p *paymentGatewayService) SetLogger(status bool, reference string, message string, gateway string) {
	currentTime := time.Now()

	logMessage := fmt.Sprintf("TimeStamp: %s, Status: %t, Reference: %v, Message: %v, Gateway: %v",
		currentTime.Format("2006-01-02 15:04:05"),
		status, reference, message, gateway)

	fmt.Println(logMessage)
}

func (p *paymentGatewayService) InitializePayment(paymentDto payment_gateway_dto.PaymentInitializationDTO) (payment_gateway_dto.PaymentInitializationResponseDTO, error) {
	switch paymentDto.Gateway {
	case PaystackPaymentGateway:
		return p.InitializePaystack(paymentDto)
	case FlutterwavePaymentGateway:
		return payment_gateway_dto.PaymentInitializationResponseDTO{}, errors.New("flutterwave not supported yet")
	default:
		return payment_gateway_dto.PaymentInitializationResponseDTO{}, errors.New("payment gateway must be paystack or flutterwave")
	}
}

func (p *paymentGatewayService) VerifyPayment(reference string, gateway string) (payment_gateway_dto.PaymentVerifyResponseDTO, error) {
	switch gateway {
	case PaystackPaymentGateway:
		return p.VerifyPaystack(reference)
	case FlutterwavePaymentGateway:
		return p.VerifyFlutterwave(reference)
	default:
		return payment_gateway_dto.PaymentVerifyResponseDTO{}, errors.New("payment gateway must be paystack or flutterwave")
	}
}

func (p *paymentGatewayService) InitializePaystack(paymentDto payment_gateway_dto.PaymentInitializationDTO) (payment_gateway_dto.PaymentInitializationResponseDTO, error) {
	paystackDto := payment_gateway_dto.Paystack{
		Amount:    paymentDto.Amount,
		Email:     paymentDto.Email,
		Reference: paymentDto.Reference,
	}
	initialize, err := p.paystackService.InitializePayment(paystackDto)

	p.SetLogger(initialize.Status, paymentDto.Reference, initialize.Message, "paystack")

	if err != nil {
		return payment_gateway_dto.PaymentInitializationResponseDTO{}, err
	}

	responseDto := payment_gateway_dto.PaymentInitializationResponseDTO{
		Status:     initialize.Status,
		PaymentURL: initialize.Data.AuthorizationURL,
	}

	return responseDto, nil
}

func (p *paymentGatewayService) InitializeFlutterwave(paymentDto payment_gateway_dto.PaymentInitializationDTO) (payment_gateway_dto.PaymentInitializationResponseDTO, error) {
	var flutterwaveDto payment_gateway_dto.InitializeFlutterwaveRequest
	var status bool

	flutterwaveDto.Amount = paymentDto.Amount
	flutterwaveDto.Customer.Email = paymentDto.Email
	flutterwaveDto.TxRef = paymentDto.Reference

	initialize, err := p.flutterwaveService.InitializePayment(flutterwaveDto)

	if initialize.Status == "success" {
		status = true
	} else {
		status = false
	}

	p.SetLogger(status, paymentDto.Reference, initialize.Message, "flutterwave")

	if err != nil {
		return payment_gateway_dto.PaymentInitializationResponseDTO{}, err
	}

	responseDto := payment_gateway_dto.PaymentInitializationResponseDTO{
		Status:     status,
		PaymentURL: initialize.Data.Link,
	}

	return responseDto, nil
}

func (p *paymentGatewayService) VerifyPaystack(reference string) (resp payment_gateway_dto.PaymentVerifyResponseDTO, err error) {

	verify, err := p.paystackService.VerifyPayment(reference)

	p.SetLogger(verify.Status, verify.Data.Reference, verify.Message, "paystack")

	if err != nil {
		return payment_gateway_dto.PaymentVerifyResponseDTO{}, err
	}

	resp.Status = verify.Status
	resp.Message = verify.Data.Message
	switch verify.Data.Status {
	case PaystackStatusSuccess:
		resp.PaymentStatus = finance_service.TransactionStatusSuccess
	case PaystackStatusFailed:
		resp.PaymentStatus = finance_service.TransactionStatusFailed
	case PaystackStatusAbandoned:
		resp.PaymentStatus = finance_service.TransactionStatusFailed
	case PaystackStatusReversed:
		resp.PaymentStatus = finance_service.TransactionStatusFailed
	default:
		resp.PaymentStatus = finance_service.TransactionStatusPending
	}

	return resp, nil
}

func (p *paymentGatewayService) VerifyFlutterwave(reference string) (resp payment_gateway_dto.PaymentVerifyResponseDTO, err error) {
	var status bool

	verify, err := p.flutterwaveService.VerifyPayment(reference)

	if verify.Status == "success" {
		status = true
	} else {
		status = false
	}

	p.SetLogger(status, verify.Data.TxRef, verify.Message, "flutterwave")

	if err != nil {
		return payment_gateway_dto.PaymentVerifyResponseDTO{}, err
	}

	resp.Status = status
	resp.Message = verify.Data.ProcessorResponse
	switch verify.Data.Status {
	case FlutterwaveStatusSuccess:
		resp.PaymentStatus = finance_service.TransactionStatusSuccess
	case FlutterwaveStatusFailed:
		resp.PaymentStatus = finance_service.TransactionStatusFailed
	case FlutterwaveStatusAbandoned:
		resp.PaymentStatus = finance_service.TransactionStatusFailed
	case FlutterwaveStatusReversed:
		resp.PaymentStatus = finance_service.TransactionStatusFailed
	default:
		resp.PaymentStatus = finance_service.TransactionStatusPending
	}

	return resp, nil
}
