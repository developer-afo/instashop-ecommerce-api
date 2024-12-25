package core_handler

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/developer-afo/instashop-ecommerce-api/dto"
	"github.com/developer-afo/instashop-ecommerce-api/handler"
	"github.com/developer-afo/instashop-ecommerce-api/lib/constants"
	"github.com/developer-afo/instashop-ecommerce-api/payload/request"
	"github.com/developer-afo/instashop-ecommerce-api/payload/response"
	coreRepository "github.com/developer-afo/instashop-ecommerce-api/repository/core"
	core_service "github.com/developer-afo/instashop-ecommerce-api/service/core"
	core_validator "github.com/developer-afo/instashop-ecommerce-api/validator/core"
)

type productHandler struct {
	productService core_service.ProductServiceInterface
	imageService   core_service.ImageServiceInterface
	validator      core_validator.ProductValidator
}

type ProductHandlerInterface interface {
	FindAllProducts(c *fiber.Ctx) error
	FindProduct(c *fiber.Ctx) error
	CreateProduct(c *fiber.Ctx) error
	UpdateProduct(c *fiber.Ctx) error
	DeleteProduct(c *fiber.Ctx) error
	FindImagesByProductId(c *fiber.Ctx) error
	DeleteImage(c *fiber.Ctx) error
	CreateImage(c *fiber.Ctx) error
}

func NewProductHandler(
	productService core_service.ProductServiceInterface,
	imageService core_service.ImageServiceInterface,
) ProductHandlerInterface {
	return &productHandler{
		productService: productService,
		imageService:   imageService,
	}
}

func (h *productHandler) GeneratePageable(c *fiber.Ctx) (pageable coreRepository.ProductPageable) {

	basePageable := handler.GeneratePageable(c)

	pageable.Page = basePageable.Page
	pageable.Size = basePageable.Size
	pageable.SortBy = basePageable.SortBy
	pageable.SortDirection = basePageable.SortDirection
	pageable.Search = basePageable.Search

	return pageable
}

func (handler *productHandler) ConvertToProductResponse(productDto dto.ProductDTO) response.ProductResponse {
	var productResp response.ProductResponse

	productResp.UUID = productDto.ID
	productResp.Slug = productDto.Slug
	productResp.Name = productDto.Name
	productResp.Description = productDto.Description
	productResp.Specification = productDto.Specification
	productResp.Price = productDto.Price
	productResp.SlashPrice = productDto.SlashPrice
	productResp.Stock = productDto.Stock
	productResp.Sales = productDto.Sales
	productResp.CreatedAt = productDto.CreatedAt

	for _, image := range productDto.Images {
		productResp.Images = append(productResp.Images, response.ImageResponse{
			Key: image.Key,
		})
	}

	return productResp
}

func (handler *productHandler) FindAllProducts(c *fiber.Ctx) error {
	var resp response.Response
	pageable := handler.GeneratePageable(c)

	products, pagination, err := handler.productService.FindAllProducts(pageable)
	if err != nil {
		resp.Status = http.StatusBadRequest
		resp.Message = err.Error()
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	productsResp := []response.ProductResponse{}

	for _, product := range products {
		productsResp = append(productsResp, handler.ConvertToProductResponse(product))
	}

	resp.Status = http.StatusOK
	resp.Message = "All Products Fetched Successfully"
	resp.Data = map[string]interface{}{"results": productsResp, "pagination": pagination}

	return c.Status(http.StatusOK).JSON(resp)
}

func (handler *productHandler) FindProduct(c *fiber.Ctx) error {
	var resp response.Response
	var productResp response.ProductResponse
	productSlug := c.Params("slug")

	product, err := handler.productService.FindProductBySlug(productSlug)

	if err != nil {
		resp.Status = http.StatusBadRequest
		resp.Message = "Product not found"
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	productResp = handler.ConvertToProductResponse(product)

	resp.Status = http.StatusOK
	resp.Message = "Product Fetched Successfully"
	resp.Data = map[string]interface{}{"product": productResp}

	return c.Status(http.StatusOK).JSON(resp)
}

func (handler *productHandler) CreateProduct(c *fiber.Ctx) error {
	var resp response.Response
	var createProductRequest request.CreateProductRequest

	if err := c.BodyParser(&createProductRequest); err != nil {
		resp.Status = constants.ClientErrorBadRequest
		resp.Message = "Invalid request payload"
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	if validation, err := handler.validator.CreateProductValidate(createProductRequest); err != nil {
		resp.Status = constants.ClientUnProcessableEntity
		resp.Message = err.Error()
		resp.Data = validation
		return c.Status(http.StatusUnprocessableEntity).JSON(resp)
	}

	// check if slash price is greater than price
	if createProductRequest.SlashPrice > createProductRequest.Price {
		resp.Status = constants.ClientUnProcessableEntity
		resp.Message = "Slash price must be less than price"

		return c.Status(http.StatusUnprocessableEntity).JSON(resp)
	}

	_, err := handler.productService.CreateProduct(createProductRequest)

	if err != nil {
		resp.Status = http.StatusBadRequest
		resp.Message = err.Error()
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	resp.Status = http.StatusOK
	resp.Message = http.StatusText(http.StatusOK)

	return c.Status(http.StatusOK).JSON(resp)
}

func (handler *productHandler) UpdateProduct(c *fiber.Ctx) error {
	var resp response.Response
	var updateProductRequest request.UpdateProductRequest
	var productDto dto.ProductDTO

	productid, err := uuid.Parse(c.Params("id"))
	if err != nil {
		resp.Status = constants.ClientUnProcessableEntity
		resp.Message = "Product ID is not a valid UUID format"

		return c.Status(http.StatusUnprocessableEntity).JSON(resp)
	}

	if err := c.BodyParser(&updateProductRequest); err != nil {
		resp.Status = constants.ClientErrorBadRequest
		resp.Message = "Invalid request payload"
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	if validation, err := handler.validator.UpdateProductValidate(updateProductRequest); err != nil {
		resp.Status = constants.ClientUnProcessableEntity
		resp.Message = err.Error()
		resp.Data = validation
		return c.Status(http.StatusUnprocessableEntity).JSON(resp)
	}

	productDto.ID = productid
	productDto.Name = updateProductRequest.Name
	productDto.Description = updateProductRequest.Description
	productDto.Specification = updateProductRequest.Specification
	productDto.Price = float64(updateProductRequest.Price)
	productDto.SlashPrice = float64(updateProductRequest.SlashPrice)
	productDto.Stock = updateProductRequest.Stock

	_, err = handler.productService.UpdateProduct(productDto)

	if err != nil {
		resp.Status = http.StatusBadRequest
		resp.Message = err.Error()
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	resp.Status = http.StatusOK
	resp.Message = http.StatusText(http.StatusOK)

	return c.Status(http.StatusOK).JSON(resp)
}

func (h *productHandler) DeleteProduct(c *fiber.Ctx) error {
	var resp response.Response

	productUUID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		fmt.Println(err)
		resp.Status = constants.ClientUnProcessableEntity
		resp.Message = "Product ID is not a valid UUID format"

		return c.Status(http.StatusUnprocessableEntity).JSON(resp)
	}

	err = h.productService.DeleteProduct(productUUID)
	if err != nil {
		resp.Status = http.StatusBadRequest
		resp.Message = "Product Not Found"
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	resp.Status = http.StatusNoContent
	resp.Message = "Product deleted successfully"

	return c.Status(http.StatusOK).JSON(resp)
}

func (h *productHandler) FindImagesByProductId(c *fiber.Ctx) error {
	var resp response.Response
	var imagesResp []response.ImageResponse
	productID := c.Params("product_id")

	images, err := h.imageService.FindImagesByProductId(productID)

	if err != nil {
		resp.Status = http.StatusBadRequest
		resp.Message = err.Error()
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	for _, image := range images {
		imagesResp = append(imagesResp, response.ImageResponse{
			Key: image.Key,
		})
	}

	resp.Status = http.StatusOK
	resp.Message = http.StatusText(http.StatusOK)
	resp.Data = map[string]interface{}{"results": imagesResp}

	return c.Status(http.StatusOK).JSON(resp)
}
func (h *productHandler) DeleteImage(c *fiber.Ctx) error {
	var resp response.Response
	key := c.Params("key")

	// TODO: validate if image is for product

	err := h.imageService.DeleteImage(key)

	if err != nil {
		resp.Status = http.StatusBadRequest
		resp.Message = err.Error()
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	resp.Status = http.StatusOK
	resp.Message = http.StatusText(http.StatusOK)

	return c.Status(http.StatusOK).JSON(resp)
}

func (h *productHandler) CreateImage(c *fiber.Ctx) error {
	var resp response.Response
	var createImageRequest request.ImageRequest
	var imageDto dto.ImageDTO

	productId, err := uuid.Parse(c.Params("product_id"))

	if err != nil {
		resp.Status = constants.ClientUnProcessableEntity
		resp.Message = "Product ID is not a valid UUID format"

		return c.Status(http.StatusUnprocessableEntity).JSON(resp)
	}

	if err := c.BodyParser(&createImageRequest); err != nil {
		resp.Status = constants.ClientErrorBadRequest
		resp.Message = "Invalid request payload"
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	if validation, err := h.validator.CreateImageValidate(createImageRequest); err != nil {
		resp.Status = constants.ClientUnProcessableEntity
		resp.Message = err.Error()
		resp.Data = validation

		return c.Status(http.StatusUnprocessableEntity).JSON(resp)
	}

	imageDto.ProductUUID = productId
	imageDto.Key = createImageRequest.Key

	_, err = h.imageService.CreateImage(imageDto)

	if err != nil {
		resp.Status = http.StatusBadRequest
		resp.Message = err.Error()
		return c.Status(http.StatusBadRequest).JSON(resp)
	}

	resp.Status = http.StatusOK
	resp.Message = http.StatusText(http.StatusOK)

	return c.Status(http.StatusOK).JSON(resp)
}
