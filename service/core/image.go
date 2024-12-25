package core_service

import (
	"github.com/developer-afo/instashop-ecommerce-api/dto"
	"github.com/developer-afo/instashop-ecommerce-api/models"
	core_repository "github.com/developer-afo/instashop-ecommerce-api/repository/core"
	"github.com/google/uuid"
)

type ImageServiceInterface interface {
	CreateImage(dto dto.ImageDTO) (dto.ImageDTO, error)
	DeleteImage(key string) error
	BatchCreateImages(images []dto.ImageDTO) error
	FindImagesByProductId(productUUID string) ([]dto.ImageDTO, error)
	ConvertToDTO(image models.Image) dto.ImageDTO
}

type imageService struct {
	imageRepository core_repository.ImageRepositoryInterface
}

func NewImageService(imageRepository core_repository.ImageRepositoryInterface) ImageServiceInterface {
	return &imageService{imageRepository: imageRepository}
}

func (service *imageService) ConvertToDTO(image models.Image) (imageDto dto.ImageDTO) {

	imageDto.ID = image.ID
	imageDto.ProductUUID = image.ProductID
	imageDto.Key = image.Key
	imageDto.CreatedAt = image.CreatedAt
	imageDto.UpdatedAt = image.UpdatedAt
	imageDto.DeletedAt = image.DeletedAt.Time

	return imageDto
}

func (service *imageService) ConvertToModel(imageDto dto.ImageDTO) (image models.Image) {

	image.ID = imageDto.ID
	image.ProductID = imageDto.ProductUUID
	image.Key = imageDto.Key
	image.CreatedAt = imageDto.CreatedAt
	image.UpdatedAt = imageDto.UpdatedAt
	image.DeletedAt.Time = imageDto.DeletedAt

	return image
}

// CreateImage implements ImageServiceInterface.
func (service *imageService) CreateImage(imageDtoArg dto.ImageDTO) (dto.ImageDTO, error) {

	imageModel := service.ConvertToModel(imageDtoArg)
	imageModel, err := service.imageRepository.CreateImage(imageModel)
	if err != nil {
		return dto.ImageDTO{}, err
	}

	return service.ConvertToDTO(imageModel), nil
}

// DeleteImage implements ImageServiceInterface.
func (service *imageService) DeleteImage(key string) error {

	err := service.imageRepository.DeleteImageByKey(key)
	if err != nil {
		return err
	}

	return nil
}

// BatchCreateImages implements ImageServiceInterface.
func (service *imageService) BatchCreateImages(imageDtos []dto.ImageDTO) error {

	var images []models.Image
	for _, imageDto := range imageDtos {
		imageDto.ID, _ = uuid.NewV7()
		images = append(images, service.ConvertToModel(imageDto))
	}

	err := service.imageRepository.BatchCreateImages(images)
	if err != nil {
		return err
	}

	return nil
}

// GetAllProductImages implements ImageServiceInterface.
func (service *imageService) FindImagesByProductId(productUUID string) ([]dto.ImageDTO, error) {

	productUuid, err := uuid.Parse(productUUID)

	if err != nil {
		return []dto.ImageDTO{}, err
	}

	images, err := service.imageRepository.FindImagesByProductId(productUuid)
	if err != nil {
		return []dto.ImageDTO{}, err
	}

	var imageDtos []dto.ImageDTO
	for _, image := range images {
		imageDtos = append(imageDtos, service.ConvertToDTO(image))
	}

	return imageDtos, nil
}
