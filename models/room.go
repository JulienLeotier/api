package models

import (
	"mime/multipart"

	"gorm.io/gorm"
)

type Room struct {
	gorm.Model
	Name               string `json:"name" gorm:"unique;not null"`
	Description        string `json:"description"`
	EnigmaFiles        []File `json:"enigma_files" gorm:"many2many:room_enigma_files;"`
	RevelationFiles    []File `json:"revelation_files" gorm:"many2many:room_revelation_files;"`
	IntroFiles         []File `json:"intro_files" gorm:"many2many:room_intro_files;"`
	OutroFiles         []File `json:"outro_files" gorm:"many2many:room_outro_files;"`
	AmbianceMusicFiles []File `json:"ambiance_music_files" gorm:"many2many:room_ambiance_music_files;"`
	VariantID          uint   `json:"variant_id"`
	Variant            User   `json:"variant" gorm:"foreignKey:VariantID"`
	DetectiveID        uint   `json:"detective_id"`
	Detective          User   `json:"detective" gorm:"foreignKey:DetectiveID"`
}

type File struct {
	gorm.Model
	URL   string `json:"url" gorm:"not null"`
	Rooms []Room `json:"-" gorm:"many2many:room_files;"`
}

type RoomCreateDTO struct {
	ID                 uint                    `json:"id"`
	Name               string                  `form:"name" binding:"required"`
	Description        string                  `form:"description"`
	EnigmaFiles        []*multipart.FileHeader `form:"enigma_files"`
	RevelationFiles    []*multipart.FileHeader `form:"revelation_files"`
	IntroFiles         []*multipart.FileHeader `form:"intro_files"`
	OutroFiles         []*multipart.FileHeader `form:"outro_files"`
	AmbianceMusicFiles []*multipart.FileHeader `form:"ambiance_music_files"`
	Variant            uint                    `form:"variant" binding:"required"`
	Detective          uint                    `form:"detective" binding:"required"`
}

type RoomResponseDTO struct {
	ID                 uint     `json:"id"`
	Name               string   `json:"name"`
	Description        string   `json:"description"`
	EnigmaFiles        []string `json:"enigma_files"`
	RevelationFiles    []string `json:"revelation_files"`
	IntroFiles         []string `json:"intro_files"`
	OutroFiles         []string `json:"outro_files"`
	AmbianceMusicFiles []string `json:"ambiance_music_files"`
	Variant            User     `json:"variant"`
	Detective          User     `json:"detective"`
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
