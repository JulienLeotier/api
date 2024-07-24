package repositories

import (
	"geniale/models"

	"gorm.io/gorm"
)

type ImageRepository interface {
	Save(image *models.Image) (*models.Image, error)
	Get(id string) (*models.Image, error)
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
