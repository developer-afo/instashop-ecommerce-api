package core_repository

import (
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/developer-afo/instashop-ecommerce-api/lib/database"
	"github.com/developer-afo/instashop-ecommerce-api/models"
	"github.com/developer-afo/instashop-ecommerce-api/repository"
)

type ProductPageable struct {
	repository.Pageable
}

// ProductRepositoryInterface is a contract that defines the methods to be implemented by ProductRepository.
type ProductRepositoryInterface interface {
	FindAllProducts(pageable ProductPageable) ([]models.Product, repository.Pagination, error)
	CreateProduct(product models.Product) (models.Product, error)
	FindProductByUUID(uuid uuid.UUID) (models.Product, error)
	FindProductBySlug(slug string) (models.Product, error)
	UpdateProduct(product models.Product) (models.Product, error)
	DeleteProduct(uuid uuid.UUID) error
}

// productRepository is a struct that defines the database connection.
type productRepository struct {
	database database.DatabaseInterface
}

// NewProductRepository is a function that returns a new instance of ProductRepository.
func NewProductRepository(database database.DatabaseInterface) ProductRepositoryInterface {
	return &productRepository{database: database}
}

// FindAllProducts is a method that returns all products.
func (p *productRepository) FindAllProducts(pageable ProductPageable) ([]models.Product, repository.Pagination, error) {
	var products []models.Product
	var product models.Product
	var pagination repository.Pagination
	var result *gorm.DB
	var errCount error

	pagination.CurrentPage = int64(pageable.Page)
	pagination.TotalItems = 0
	pagination.TotalPages = 1

	offset := (pageable.Page - 1) * pageable.Size
	model := p.database.Connection().
		Model(&product).
		Preload("Images").
		Select("products.*, COALESCE(SUM(order_items.quantity), 0) as sales").
		Joins("LEFT JOIN order_items ON order_items.product_id = products.id").
		Group("products.id").
		Order("CASE WHEN products.stock = 0 THEN 1 ELSE 0 END ASC")

	if len(strings.TrimSpace(pageable.Search)) > 0 {
		model = model.Where("LOWER(products.name) LIKE ?", "%"+strings.ToLower(pageable.Search)+"%") // Specify the table name 'products' for the 'name' column
	}

	errCount = model.Count(&pagination.TotalItems).Error
	paginatedQuery := model.Offset(int(offset)).Limit(int(pageable.Size))
	result = paginatedQuery.Model(&models.Product{}).Where(product).Find(&products)

	if result.Error != nil {
		return nil, pagination, result.Error
	}

	if errCount != nil {
		return nil, pagination, errCount
	}

	pagination.TotalPages = (pagination.TotalItems + int64(pageable.Size) - 1) / int64(pageable.Size)

	return products, pagination, nil
}

// CreateProduct is a method that creates a new product.
func (p *productRepository) CreateProduct(product models.Product) (models.Product, error) {
	product.Prepare()

	err := p.database.Connection().Create(&product).Error

	return product, err
}

// FindProductByUUID is a method that returns a product by its ID.
func (p *productRepository) FindProductByUUID(uuid uuid.UUID) (models.Product, error) {
	var product models.Product

	err := p.database.Connection().
		Model(&models.Product{}).
		Where("id = ?", uuid).
		First(&product).Error

	return product, err
}

// FindProductBySlug is a method that returns a product by its slug.
func (p *productRepository) FindProductBySlug(slug string) (models.Product, error) {
	var product models.Product

	err := p.database.Connection().
		Model(&models.Product{}).
		Where("products.slug = ?", slug).
		Preload("Images").
		Select("products.*, COALESCE(SUM(order_items.quantity), 0) as sales").
		Joins("LEFT JOIN order_items ON order_items.product_id = products.id").
		Group("products.id").
		First(&product).Error

	return product, err
}

// UpdateProduct is a method that updates a product.
func (p *productRepository) UpdateProduct(product models.Product) (models.Product, error) {

	err := p.database.Connection().
		Model(&models.Product{}).
		Where("id = ?", product.ID).
		Select("category_id", "name", "description", "specification", "price", "slash_price", "stock", "brand").
		Updates(&product).Error

	return product, err
}

// DeleteProduct is a method that deletes a product.
func (p *productRepository) DeleteProduct(uuid uuid.UUID) error {

	product, err := p.FindProductByUUID(uuid)

	if err != nil {
		return err
	}

	err = p.database.Connection().Delete(&product).Error

	return err
}
