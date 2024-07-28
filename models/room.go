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
	VariantID   uint   `json:"variant_id"`
	Variant     User   `json:"variant" gorm:"foreignKey:VariantID"`
	DetectiveID uint   `json:"detective_id"`
	Detective   User   `json:"detective" gorm:"foreignKey:DetectiveID"`
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
	Variant     uint                    `form:"variant" binding:"required"`
	Detective   uint                    `form:"detective" binding:"required"`
}

type RoomResponseDTO struct {
	ID          uint     `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Files       []string `json:"files"`
	Variant     User     `json:"variant"`
	Detective   User     `json:"detective"`
}

type RoomRole struct {
	gorm.Model
	RoomID   uint   `json:"room_id" gorm:"foreignKey:RoomID"`
	Room     Room   `json:"room"`
	UserID   uint   `json:"user_id" gorm:"foreignKey:UserID"`
	User     User   `json:"user"`
	Role     string `json:"role"`
	Name     string `json:"name"`
	Job      string `json:"job"`
	Region   string `json:"region"`
	Passion  string `json:"passion"`
	Anecdote string `json:"anecdote"`
}

type RoomRoleCreateDTO struct {
	RoomID   uint   `json:"room_id" binding:"required"`
	UserID   uint   `json:"user_id" binding:"required"`
	Role     string `json:"role" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Job      string `json:"job" binding:"required"`
	Region   string `json:"region" binding:"required"`
	Passion  string `json:"passion" binding:"required"`
	Anecdote string `json:"anecdote" binding:"required"`
}
