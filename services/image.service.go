package services

import (
	"api/models"
	"api/repositories"
	"strconv"
)

type ImageDTO struct {
	ID       uint   `json:"id"`
	FilePath string `json:"file_path"`
}

type ImageService interface {
	UploadImage(filePath string) (*ImageDTO, error)
	GetImage(id string) (*ImageDTO, error)
	GetImages() ([]ImageDTO, error)
	RemoveImage(id models.Image) error
}

type imageService struct {
	repository repositories.ImageRepository
}

func NewImageService(repo repositories.ImageRepository) ImageService {
	return &imageService{repository: repo}
}

func (s *imageService) UploadImage(filePath string) (*ImageDTO, error) {
	image := &models.Image{FilePath: filePath}
	savedImage, err := s.repository.Save(image)
	if err != nil {
		return nil, err
	}
	return &ImageDTO{ID: savedImage.ID, FilePath: savedImage.FilePath}, nil
}

func (s *imageService) GetImage(id string) (*ImageDTO, error) {
	image, err := s.repository.Get(id)
	if err != nil {
		return nil, err
	}
	return &ImageDTO{ID: image.ID, FilePath: image.FilePath}, nil
}

func (s *imageService) GetImages() ([]ImageDTO, error) {
	images, err := s.repository.GetAll()
	if err != nil {
		return nil, err
	}
	var result []ImageDTO
	for _, img := range images {
		result = append(result, ImageDTO{ID: img.ID, FilePath: img.FilePath})
	}
	return result, nil
}

func (s *imageService) RemoveImage(image models.Image) error {
	id := strconv.Itoa(int(image.ID))
	return s.repository.Remove(id)
}
