package core_repository

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/developer-afo/instashop-ecommerce-api/lib/database"
	"github.com/developer-afo/instashop-ecommerce-api/models"
)

type ImageRepositoryInterface interface {
	CreateImage(image models.Image) (models.Image, error)
	BatchCreateImages(images []models.Image) error
	FindImagesByProductId(productId uuid.UUID) ([]models.Image, error)
	DeleteImageByID(id uuid.UUID) error
	DeleteImageByKey(key string) error
}

type imageRepository struct {
	database database.DatabaseInterface
}

func NewImageRepository(database database.DatabaseInterface) ImageRepositoryInterface {
	return &imageRepository{database: database}
}

// CreateImage implements ImageRepositoryInterface.
func (i *imageRepository) CreateImage(image models.Image) (models.Image, error) {
	var productCount int64
	image.Prepare()

	// check if product exists
	err := i.database.Connection().Model(&models.Product{}).Where("id = ?", image.ProductID).Count(&productCount).Error

	if err != nil {
		return models.Image{}, err
	}

	if productCount == 0 {
		return models.Image{}, fmt.Errorf("product with id %s not found", image.ProductID)
	}

	err = i.database.Connection().Create(&image).Error

	return image, err
}

// BatchCreateImages implements ImageRepositoryInterface.
func (i *imageRepository) BatchCreateImages(images []models.Image) error {

	err := i.database.Connection().Create(&images).Error

	return err
}

// FindImagesByProductId implements ImageRepositoryInterface.
func (i *imageRepository) FindImagesByProductId(productId uuid.UUID) (images []models.Image, err error) {

	err = i.database.Connection().Model(&models.Image{}).Where("product_id = ?", productId).Find(&images).Error

	return images, err
}

// FindImageById implements ImageRepositoryInterface.
func (i *imageRepository) FindImageById(id uuid.UUID) (image models.Image, err error) {

	err = i.database.Connection().Model(&models.Image{}).Where("id = ?", id).First(&image).Error

	return image, err
}

// FindImageByKey implements ImageRepositoryInterface.
func (i *imageRepository) FindImageByKey(key string) (image models.Image, err error) {

	err = i.database.Connection().Model(&models.Image{}).Where("key = ?", key).First(&image).Error

	return image, err
}

// DeleteImageByID implements ImageRepositoryInterface.
func (i *imageRepository) DeleteImageByID(id uuid.UUID) error {

	image, imgErr := i.FindImageById(id)

	if imgErr != nil {
		return imgErr
	}

	err := i.database.Connection().Delete(&image).Error

	return err
}

// DeleteImageByKey implements ImageRepositoryInterface.
func (i *imageRepository) DeleteImageByKey(key string) error {

	image, imgErr := i.FindImageByKey(key)

	if imgErr != nil {
		return imgErr
	}

	err := i.database.Connection().
		Where("id = ?", image.ID).
		Delete(&image).Error

	return err
}
