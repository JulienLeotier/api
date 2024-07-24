package models

type Room struct {
	ID          uint   `json:"id" gorm:"primary_key"`
	Name        string `json:"name" gorm:"unique;not null"`
	Description string `json:"description"`
}

type RoomCreateDTO struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type RoomUpdateDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type RoomResponseDTO struct {
	Message string `json:"message"`
	Room    Room   `json:"room"`
}
