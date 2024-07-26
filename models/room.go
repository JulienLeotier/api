package models

import (
	"mime/multipart"

	"gorm.io/gorm"
)

type Room struct {
	gorm.Model
	Name        string `json:"name" gorm:"unique;not null"`
	Description string `json:"description"`
	Files       []File `json:"files" gorm:"many2many:room_files;"`
}

type File struct {
	gorm.Model
	URL   string `json:"url" gorm:"not null"`
	Rooms []Room `json:"-" gorm:"many2many:room_files;"`
}

type RoomCreateDTO struct {
	ID          uint                    `json:"id"`
	Name        string                  `form:"name" binding:"required"`
	Description string                  `form:"description"`
	Files       []*multipart.FileHeader `form:"file" binding:"required"`
}

type RoomResponseDTO struct {
	ID          uint     `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Files       []string `json:"files"`
}
