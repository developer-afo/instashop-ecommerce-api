package userResponse

import (
	"github.com/developer-afo/instashop-ecommerce-api/dto"
	"github.com/developer-afo/instashop-ecommerce-api/payload/response"
)

type LoginResponse struct {
	response.Response

	Data dto.LoginResponseDTO `json:"data"`
}
