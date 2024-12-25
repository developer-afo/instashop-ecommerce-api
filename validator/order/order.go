package order_validator

import (
	validation "github.com/go-ozzo/ozzo-validation"

	"github.com/developer-afo/instashop-ecommerce-api/payload/request"
	"github.com/developer-afo/instashop-ecommerce-api/validator"
)

type OrderValidator struct {
	validator.Validator[request.CreateOrderRequest]
}

func (validator *OrderValidator) CreateOrderValidate(req request.CreateOrderRequest) (map[string]interface{}, error) {
	err := validation.ValidateStruct(&req,
		validation.Field(&req.PaymentMethod, validation.Required),
		validation.Field(&req.Items, validation.Required, validation.Each(validation.Required)),
	)

	if err != nil {
		return validator.ValidateErr(err)
	}

	return nil, nil
}
