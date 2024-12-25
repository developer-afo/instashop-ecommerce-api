package core_service

import (
	"gorm.io/gorm"

	"github.com/developer-afo/instashop-ecommerce-api/dto"
	"github.com/developer-afo/instashop-ecommerce-api/lib/helper"
	"github.com/developer-afo/instashop-ecommerce-api/models"
	"github.com/developer-afo/instashop-ecommerce-api/payload/request"
	"github.com/developer-afo/instashop-ecommerce-api/repository"
	coreRepository "github.com/developer-afo/instashop-ecommerce-api/repository/core"
	"github.com/google/uuid"
)

type ProductServiceInterface interface {
	CreateProduct(dto request.CreateProductRequest) (dto.ProductDTO, error)
	FindAllProducts(pageable coreRepository.ProductPageable) ([]dto.ProductDTO, repository.Pagination, error)
	FindProductByUUID(id string) (dto.ProductDTO, error)
	FindProductBySlug(slug string) (dto.ProductDTO, error)
	UpdateProduct(dto dto.ProductDTO) (dto.ProductDTO, error)
	DeleteProduct(id uuid.UUID) error
	ConvertToDTO(product models.Product) dto.ProductDTO
}

type productService struct {
	productRepository coreRepository.ProductRepositoryInterface
	imageService      ImageServiceInterface
}

func NewProductService(
	productRepository coreRepository.ProductRepositoryInterface,
	imageService ImageServiceInterface,
) ProductServiceInterface {
	return &productService{
		productRepository: productRepository,
		imageService:      imageService,
	}
}

func (service *productService) ConvertToDTO(product models.Product) (productDto dto.ProductDTO) {

	productDto.ID = product.ID
	productDto.Name = product.Name
	productDto.Slug = product.Slug
	productDto.Description = product.Description
	productDto.Specification = product.Specification
	productDto.Price = product.Price
	productDto.SlashPrice = product.SlashPrice
	productDto.Stock = product.Stock
	productDto.Sales = product.Sales
	productDto.CreatedAt = product.CreatedAt
	productDto.UpdatedAt = product.UpdatedAt
	productDto.DeletedAt = product.DeletedAt.Time
	for _, image := range product.Images {
		productDto.Images = append(productDto.Images, service.imageService.ConvertToDTO(image))
	}
	return productDto
}

func (service *productService) ConvertToModel(productDto dto.ProductDTO) (product models.Product) {

	product.ID = productDto.ID
	product.Name = productDto.Name
	product.Slug = productDto.Slug
	product.Description = productDto.Description
	product.Specification = productDto.Specification
	product.Price = productDto.Price
	product.SlashPrice = productDto.SlashPrice
	product.Stock = productDto.Stock
	product.CreatedAt = productDto.CreatedAt
	product.UpdatedAt = productDto.UpdatedAt
	product.DeletedAt.Time = productDto.DeletedAt

	return product
}

// CreateProduct implements ProductServiceInterface.
func (service *productService) CreateProduct(createProduct request.CreateProductRequest) (dto.ProductDTO, error) {
	var productDto dto.ProductDTO
	var imageDtos []dto.ImageDTO
	slug := helper.GenerateSlug(createProduct.Name)

	// check if slug already exists
	_, err := service.productRepository.FindProductBySlug(slug)

	if err != nil && err != gorm.ErrRecordNotFound {
		return dto.ProductDTO{}, err
	}

	if err == gorm.ErrRecordNotFound {
		productDto.Slug = slug
	} else {
		// add timestamp to slug
		productDto.Slug = slug + "-" + helper.GenerateTimestamp()
	}

	productDto.Name = createProduct.Name
	productDto.Description = createProduct.Description
	productDto.Specification = createProduct.Specification
	productDto.Price = float64(createProduct.Price)
	productDto.SlashPrice = float64(createProduct.SlashPrice)
	productDto.Stock = createProduct.Stock

	product := service.ConvertToModel(productDto)
	newRecord, err := service.productRepository.CreateProduct(product)

	if err != nil {
		return dto.ProductDTO{}, err
	}

	for _, image := range createProduct.Images {
		imageDtos = append(imageDtos, dto.ImageDTO{
			ProductUUID: newRecord.ID,
			Key:         image,
		})
	}

	imageErr := service.imageService.BatchCreateImages(imageDtos)

	if imageErr != nil {
		return dto.ProductDTO{}, imageErr
	}

	return service.ConvertToDTO(newRecord), nil
}

// FindAllProducts implements ProductServiceInterface.
func (service *productService) FindAllProducts(pageable coreRepository.ProductPageable) ([]dto.ProductDTO, repository.Pagination, error) {

	products, pagination, err := service.productRepository.FindAllProducts(pageable)

	if err != nil {
		return nil, pagination, err
	}

	var productDtos []dto.ProductDTO
	for _, product := range products {
		eachProduct := service.ConvertToDTO(product)

		productDtos = append(productDtos, eachProduct)
	}

	return productDtos, pagination, nil
}

// FindProductByUUID implements ProductServiceInterface.
func (service *productService) FindProductByUUID(id string) (dto.ProductDTO, error) {
	uuid, err := uuid.Parse(id)

	if err != nil {
		return dto.ProductDTO{}, err
	}

	product, err := service.productRepository.FindProductByUUID(uuid)
	if err != nil {
		return dto.ProductDTO{}, err
	}

	return service.ConvertToDTO(product), nil
}

// FindProductBySlug implements ProductServiceInterface.
func (service *productService) FindProductBySlug(slug string) (dto.ProductDTO, error) {

	product, err := service.productRepository.FindProductBySlug(slug)
	if err != nil {
		return dto.ProductDTO{}, err
	}

	return service.ConvertToDTO(product), nil
}

// UpdateProduct implements ProductServiceInterface.
func (service *productService) UpdateProduct(productDtoArg dto.ProductDTO) (dto.ProductDTO, error) {

	product := service.ConvertToModel(productDtoArg)
	product, err := service.productRepository.UpdateProduct(product)
	if err != nil {
		return dto.ProductDTO{}, err
	}

	return service.ConvertToDTO(product), nil
}

func (s *productService) DeleteProduct(id uuid.UUID) error {
	return s.productRepository.DeleteProduct(id)
}
