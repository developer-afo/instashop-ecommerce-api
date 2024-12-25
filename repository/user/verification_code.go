package user_repository

import (
	"github.com/google/uuid"

	"github.com/developer-afo/instashop-ecommerce-api/lib/database"
	"github.com/developer-afo/instashop-ecommerce-api/models"
)

type verificationCodeRepository struct {
	database database.DatabaseInterface
}

type VerificationCodeRepositoryInterface interface {
	CreateVerificationCode(code models.VerificationCode) (models.VerificationCode, error)
	FindCodeAndUserId(code string, userId uuid.UUID) (models.VerificationCode, error)
	FindCodeByUserId(userId uuid.UUID) (models.VerificationCode, error)
	FindByCode(code string) (models.VerificationCode, error)
	DeleteVerificationCode(userId uuid.UUID) error
}

func NewVerificationCodeRepository(
	database database.DatabaseInterface,
) VerificationCodeRepositoryInterface {
	return &verificationCodeRepository{
		database: database,
	}
}

// CreateVerificationCode implements CodeRepositoryInterface.
func (c *verificationCodeRepository) CreateVerificationCode(code models.VerificationCode) (models.VerificationCode, error) {
	code.Prepare()

	err := c.database.Connection().Create(&code).Error

	if err != nil {

		return models.VerificationCode{}, err
	}

	return code, err
}

// DeleteVerificationCode implements CodeRepositoryInterface.
func (c *verificationCodeRepository) DeleteVerificationCode(userId uuid.UUID) error {

	code, err := c.FindCodeByUserId(userId)

	if err != nil {
		return err
	}

	err = c.database.Connection().Delete(&code).Error

	if err != nil {

		return err
	}

	return nil
}

// FindCodeAndUserId implements CodeRepositoryInterface.
func (c *verificationCodeRepository) FindCodeAndUserId(code string, userId uuid.UUID) (codeModel models.VerificationCode, err error) {

	err = c.database.Connection().Model(&models.VerificationCode{}).Where("code = ? AND user_id = ?", code, userId).First(&codeModel).Error

	return codeModel, err
}

// FindCodeByUserId implements CodeRepositoryInterface.
func (c *verificationCodeRepository) FindCodeByUserId(userId uuid.UUID) (code models.VerificationCode, err error) {

	err = c.database.Connection().Model(&models.VerificationCode{}).Where("user_id = ?", userId).First(&code).Error

	return code, err
}

// FindByCode implements CodeRepositoryInterface.
func (c *verificationCodeRepository) FindByCode(code string) (codeModel models.VerificationCode, err error) {

	err = c.database.Connection().Model(&models.VerificationCode{}).Where("code = ?", code).First(&codeModel).Error

	return codeModel, err
}
