package finance_service

import (
	"github.com/google/uuid"

	"github.com/developer-afo/instashop-ecommerce-api/dto"
	"github.com/developer-afo/instashop-ecommerce-api/lib/helper"
	"github.com/developer-afo/instashop-ecommerce-api/models"
	"github.com/developer-afo/instashop-ecommerce-api/repository"
	finance_repository "github.com/developer-afo/instashop-ecommerce-api/repository/finance"
)

var (
	TransactionTypeCredit     = "credit"
	TransactionTypeDebit      = "debit"
	TransactionStatusSuccess  = "success"
	TransactionStatusPending  = "pending"
	TransactionStatusFailed   = "failed"
	TransactionMethodWallet   = "wallet"
	TransactionMethodGateway  = "gateway"
	TransactionMethodTransfer = "transfer"
	TransactionVendorPayStack = "paystack"
	TransactionVendorMazimart = "instashop"
)

type TransactionServiceInterface interface {
	FindTransactionByUUID(TransactionID string) (dto.TransactionDTO, error)
	FindTransactionByReference(reference string) (dto.TransactionDTO, error)
	FindAllTransactions(pageable finance_repository.TransactionPageable) ([]dto.TransactionDTO, repository.Pagination, error)
	CreateTransaction(transaction dto.TransactionDTO) (dto.TransactionDTO, error)
	UpdateTransaction(transaction dto.TransactionDTO) (dto.TransactionDTO, error)
	ConfirmTransaction(transactionId string) (dto.TransactionDTO, error)
	FailTransaction(transactionId string) (dto.TransactionDTO, error)
	ConvertToDTO(transaction models.Transaction) dto.TransactionDTO
}

type transactionService struct {
	transactionRepository finance_repository.TransactionRepositoryInterface
}

func NewTransactionService(transactionRepository finance_repository.TransactionRepositoryInterface) TransactionServiceInterface {
	return &transactionService{transactionRepository: transactionRepository}
}

func (t *transactionService) ConvertToDTO(transaction models.Transaction) (transactionDto dto.TransactionDTO) {

	transactionDto.ID = transaction.ID
	transactionDto.UserID = transaction.UserID
	transactionDto.Amount = transaction.Amount
	transactionDto.Type = transaction.Type
	transactionDto.Reference = transaction.Reference
	transactionDto.Description = transaction.Description
	transactionDto.ShortDesc = transaction.ShortDesc
	transactionDto.Status = transaction.Status
	transactionDto.Method = transaction.Method
	transactionDto.Vendor = transaction.Vendor
	transactionDto.CreatedAt = transaction.CreatedAt
	transactionDto.UpdatedAt = transaction.UpdatedAt
	transactionDto.DeletedAt = transaction.DeletedAt.Time

	return transactionDto
}

func (t *transactionService) ConvertToModel(transactionDto dto.TransactionDTO) (transaction models.Transaction) {

	transaction.ID = transactionDto.ID
	transaction.UserID = transactionDto.UserID
	transaction.Amount = transactionDto.Amount
	transaction.Type = transactionDto.Type
	transaction.Reference = transactionDto.Reference
	transaction.Description = transactionDto.Description
	transaction.ShortDesc = transactionDto.ShortDesc
	transaction.Status = transactionDto.Status
	transaction.Method = transactionDto.Method
	transaction.Vendor = transactionDto.Vendor
	transaction.CreatedAt = transactionDto.CreatedAt
	transaction.UpdatedAt = transactionDto.UpdatedAt
	transaction.DeletedAt.Time = transactionDto.DeletedAt

	return transaction
}

// FindTransactionByUUID implements TransactionServiceInterface.
func (t *transactionService) FindTransactionByUUID(TransactionID string) (dto.TransactionDTO, error) {

	_uuid, err := uuid.Parse(TransactionID)

	if err != nil {
		return dto.TransactionDTO{}, err
	}

	transaction, err := t.transactionRepository.FindTransactionByUUID(_uuid)
	if err != nil {
		return dto.TransactionDTO{}, err
	}

	return t.ConvertToDTO(transaction), nil
}

// FindTransactionByReference implements TransactionServiceInterface.
func (t *transactionService) FindTransactionByReference(reference string) (dto.TransactionDTO, error) {

	transaction, err := t.transactionRepository.FindTransactionByReference(reference)
	if err != nil {
		return dto.TransactionDTO{}, err
	}

	return t.ConvertToDTO(transaction), nil
}

// FindAllTransactions implements TransactionServiceInterface.
func (t *transactionService) FindAllTransactions(pageable finance_repository.TransactionPageable) ([]dto.TransactionDTO, repository.Pagination, error) {
	transactions := []dto.TransactionDTO{}

	_transactions, pagination, err := t.transactionRepository.FindAllTransactions(pageable)
	if err != nil {
		return nil, pagination, err
	}

	for _, transaction := range _transactions {
		transactions = append(transactions, t.ConvertToDTO(transaction))
	}

	return transactions, pagination, nil
}

// CreateTransaction implements TransactionServiceInterface.
func (t *transactionService) CreateTransaction(transaction dto.TransactionDTO) (dto.TransactionDTO, error) {

	reference, err := helper.GenerateSnowflakeID()

	if err != nil {
		return dto.TransactionDTO{}, err
	}

	transaction.Reference = helper.Int64ToString(reference)

	transactionModel := t.ConvertToModel(transaction)
	transactionModel, err = t.transactionRepository.CreateTransaction(transactionModel)

	return t.ConvertToDTO(transactionModel), err
}

// UpdateTransaction implements TransactionServiceInterface.
func (t *transactionService) UpdateTransaction(transaction dto.TransactionDTO) (dto.TransactionDTO, error) {

	transactionModel := t.ConvertToModel(transaction)
	transactionModel, err := t.transactionRepository.UpdateTransaction(transactionModel)

	return t.ConvertToDTO(transactionModel), err
}

// ConfirmTransaction implements TransactionServiceInterface.
func (t *transactionService) ConfirmTransaction(transactionId string) (dto.TransactionDTO, error) {

	transaction, err := t.FindTransactionByUUID(transactionId)
	if err != nil {
		return dto.TransactionDTO{}, err
	}

	transaction.Status = TransactionStatusSuccess

	return t.UpdateTransaction(transaction)
}

// FailTransaction implements TransactionServiceInterface.
func (t *transactionService) FailTransaction(transactionId string) (dto.TransactionDTO, error) {

	transaction, err := t.FindTransactionByUUID(transactionId)
	if err != nil {
		return dto.TransactionDTO{}, err
	}

	transaction.Status = TransactionStatusFailed

	return t.UpdateTransaction(transaction)
}
