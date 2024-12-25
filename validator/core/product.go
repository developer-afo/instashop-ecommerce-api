package core_validator

import (
	validation "github.com/go-ozzo/ozzo-validation"

	"github.com/developer-afo/instashop-ecommerce-api/payload/request"
	"github.com/developer-afo/instashop-ecommerce-api/validator"
)

type ProductValidator struct {
	validator.Validator[request.CreateProductRequest]
}

func (validator *ProductValidator) CreateProductValidate(req request.CreateProductRequest) (map[string]interface{}, error) {
	err := validation.ValidateStruct(&req,
		validation.Field(&req.Name, validation.Required, validation.Length(3, 50)),
		validation.Field(&req.Description, validation.Required),
		validation.Field(&req.Specification, validation.Required),
		validation.Field(&req.Price, validation.Required, validation.Min(0)),
		validation.Field(&req.Stock, validation.Required, validation.Min(0)),
		validation.Field(&req.SlashPrice, validation.Max(req.Price)),
		validation.Field(&req.Images, validation.Required, validation.Each(validation.Required, validation.Length(3, 100))),
	)

	if err != nil {
		return validator.ValidateErr(err)
	}

	return nil, nil
}

func (validator *ProductValidator) UpdateProductValidate(req request.UpdateProductRequest) (map[string]interface{}, error) {
	err := validation.ValidateStruct(&req,
		validation.Field(&req.Name, validation.Required, validation.Length(3, 50)),
		validation.Field(&req.Description, validation.Required),
		validation.Field(&req.Specification, validation.Required),
		validation.Field(&req.Price, validation.Required, validation.Min(0)), // TODO: add is.Int validation on fields like this
		validation.Field(&req.Stock, validation.Min(0)),
		validation.Field(&req.SlashPrice, validation.Max(req.Price)),
	)

	if err != nil {
		return validator.ValidateErr(err)
	}

	return nil, nil
}

func (v *ProductValidator) CreateImageValidate(req request.ImageRequest) (map[string]interface{}, error) {
	err := validation.ValidateStruct(&req,
		validation.Field(&req.Key, validation.Required, validation.Length(3, 100)),
	)

	if err != nil {
		return v.ValidateErr(err)
	}

	return nil, nil
}
