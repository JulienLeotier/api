package models

type Image struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	FilePath string `json:"file_path"`
}

type ImageCreateDTO struct {
	FilePath string `json:"file" validate:"required"`
}
