package request

type CreateOrderRequest struct {
	PaymentMethod string                   `json:"payment_method"`
	Items         []CreateOrderRequestItem `json:"items"`
}

type CreateOrderRequestItem struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"` // TODO: quantity should be greater than 0
}

type CreateShippingTypeRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
}

type UpdateShippingTypeRequest struct {
	CreateShippingTypeRequest
}

type CreateShippingAddressRequest struct {
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Phone           string `json:"phone"`
	AlternatePhone  string `json:"alternate_phone"`
	Address         string `json:"address"`
	CityID          string `json:"city_id"`
	StateID         string `json:"state_id"`
	ClosestLandmark string `json:"closest_landmark"`
}

type UpdateShippingAddressRequest struct {
	CreateShippingAddressRequest
}

type CreateCountryRequest struct {
	Name string `json:"name"`
}

type UpdateCountryRequest struct {
	CreateCountryRequest
}

type CreateLocationRequest struct {
	CountryID string `json:"country_id"`
	States    []struct {
		Name   string `json:"name"`
		Cities []struct {
			Name  string  `json:"name"`
			Price float64 `json:"price"`
		} `json:"cities"`
	} `json:"states"`
}

type CreateStateRequest struct {
	Name string `json:"name"`
}

type CreateCityRequest struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type CreateCouponRequest struct {
	Type        string `json:"type"`
	Value       int    `json:"value"`
	Description string `json:"description"`
}
