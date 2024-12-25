package finance_repository

import (
	"strings"

	"github.com/google/uuid"

	"github.com/developer-afo/instashop-ecommerce-api/lib/database"
	"github.com/developer-afo/instashop-ecommerce-api/models"
	"github.com/developer-afo/instashop-ecommerce-api/repository"
)

type TransactionPageable struct {
	repository.Pageable

	Method string
	Type   string
	Status string
	Vendor string
	UserID uuid.UUID
}

type TransactionRepositoryInterface interface {
	FindTransactionByUUID(uuid uuid.UUID) (models.Transaction, error)
	FindAllTransactions(pageable TransactionPageable) ([]models.Transaction, repository.Pagination, error)
	FindTransactionByReference(reference string) (models.Transaction, error)
	CreateTransaction(transaction models.Transaction) (models.Transaction, error)
	UpdateTransaction(transaction models.Transaction) (models.Transaction, error)
}

type transactionRepository struct {
	database database.DatabaseInterface
}

func NewTransactionRepository(database database.DatabaseInterface) TransactionRepositoryInterface {
	return &transactionRepository{database: database}
}

// FindAllTransactions is a method that returns all transactions.
func (t *transactionRepository) FindAllTransactions(pageable TransactionPageable) ([]models.Transaction, repository.Pagination, error) {
	var transactions []models.Transaction
	var transaction models.Transaction
	var pagination repository.Pagination
	var errCount error

	pagination.CurrentPage = int64(pageable.Page)
	pagination.TotalItems = 0
	pagination.TotalPages = 1

	offset := (pageable.Page - 1) * pageable.Size
	model := t.database.Connection().Model(&transaction)

	if len(strings.TrimSpace(pageable.Search)) > 0 {
		model.Where("transactions.reference LIKE ?", "%"+pageable.Search+"%")
	}

	if len(strings.TrimSpace(pageable.Method)) > 0 {
		model.Where("transactions.method = ?", pageable.Method)
	}

	if len(strings.TrimSpace(pageable.Type)) > 0 {
		model.Where("transactions.type = ?", pageable.Type)
	}

	if len(strings.TrimSpace(pageable.Status)) > 0 {
		model.Where("transactions.status = ?", pageable.Status)
	}

	if len(strings.TrimSpace(pageable.Vendor)) > 0 {
		model.Where("transactions.vendor = ?", pageable.Vendor)
	}

	if pageable.UserID != uuid.Nil {
		model.Where("transactions.user_id = ?", pageable.UserID)
	}

	errCount = model.Count(&pagination.TotalItems).Error
	paginatedQuery := model.Offset(int(offset)).Limit(int(pageable.Size)).Order(pageable.SortBy + " " + pageable.SortDirection)

	if err := paginatedQuery.Model(&models.Transaction{}).Where(transaction).Find(&transactions).Error; err != nil {
		return nil, pagination, err
	}

	if errCount != nil {
		return nil, pagination, errCount
	}

	pagination.TotalPages = pagination.TotalItems / int64(pageable.Size)

	if pagination.TotalPages == 0 {
		pagination.TotalPages = 1
	}

	return transactions, pagination, nil
}

// FindTransactionByUUID implements TransactionRepositoryInterface.
func (t *transactionRepository) FindTransactionByUUID(uuid uuid.UUID) (transaction models.Transaction, err error) {

	err = t.database.Connection().Model(&models.Transaction{}).Where("id = ?", uuid).First(&transaction).Error

	return transaction, err
}

// FindTransactionByReference implements TransactionRepositoryInterface.
func (t *transactionRepository) FindTransactionByReference(reference string) (transaction models.Transaction, err error) {

	err = t.database.Connection().Model(&models.Transaction{}).Where("reference = ?", reference).First(&transaction).Error

	return transaction, err
}

// CreateTransaction implements TransactionRepositoryInterface.
func (t *transactionRepository) CreateTransaction(transaction models.Transaction) (models.Transaction, error) {
	transaction.Prepare()

	err := t.database.Connection().Create(&transaction).Error

	return transaction, err
}

// UpdateTransaction implements TransactionRepositoryInterface.
func (t *transactionRepository) UpdateTransaction(transaction models.Transaction) (models.Transaction, error) {

	err := t.database.Connection().
		Model(&models.Transaction{}).
		Where("id = ?", transaction.ID).
		Updates(transaction).Error

	return transaction, err
}
