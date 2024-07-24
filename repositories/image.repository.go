package repositories

import (
	"geniale/models"

	"gorm.io/gorm"
)

type ImageRepository interface {
	Save(image *models.Image) (*models.Image, error)
	Get(id string) (*models.Image, error)
	GetAll() ([]models.Image, error)
	Remove(id string) error
}

type imageRepository struct {
	DB *gorm.DB
}

func NewImageRepository(db *gorm.DB) ImageRepository {
	return &imageRepository{DB: db}
}

func (r *imageRepository) Save(image *models.Image) (*models.Image, error) {
	if err := r.DB.Create(image).Error; err != nil {
		return nil, err
	}
	return image, nil
}

func (r *imageRepository) Get(id string) (*models.Image, error) {
	var image models.Image
	result := r.DB.First(&image, id)
	return &image, result.Error
}

func (r *imageRepository) GetAll() ([]models.Image, error) {
	var images []models.Image
	result := r.DB.Find(&images)
	return images, result.Error
}

func (r *imageRepository) Remove(id string) error {
	return r.DB.Delete(&models.Image{}, id).Error
}
