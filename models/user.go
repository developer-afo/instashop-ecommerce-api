package models

import (
	"github.com/google/uuid"

	"github.com/developer-afo/instashop-ecommerce-api/lib/database"
)

type User struct {
	database.BaseModel

	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Email           string `json:"email"`
	IsEmailVerified bool   `json:"is_email_verified"`
	Password        string `json:"password"`
	Role            string `json:"role"`
}

type VerificationCode struct {
	database.BaseModel

	UserID  uuid.UUID `json:"user_id"`
	Code    string    `json:"code"`
	Purpose string    `json:"purpose"`
}
