package services

import (
	"geniale/models"
	"geniale/repositories"
)

type ImageService interface {
	UploadImage(filePath string) (*models.Image, error)
	GetImage(id string) (*models.Image, error)
}

type imageService struct {
	repository repositories.ImageRepository
}

func NewImageService(repo repositories.ImageRepository) ImageService {
	return &imageService{repository: repo}
}

func (s *imageService) UploadImage(filePath string) (*models.Image, error) {
	image := &models.Image{FilePath: filePath}
	return s.repository.Save(image)
}

func (s *imageService) GetImage(id string) (*models.Image, error) {
	return s.repository.Get(id)
}
