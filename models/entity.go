package models

type Entity struct {
	ID      uint   `json:"id" gorm:"primary_key"`
	Name    string `json:"name" gorm:"unique;not null"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Website string `json:"website"`
	User    []User `gorm:"many2many:entity_users;"`
}

type EntityCreateDTO struct {
	Name    string `json:"name" validate:"required"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Website string `json:"website"`
}

type EntityUpdateDTO struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Website string `json:"website"`
}

type EntityResponseDTO struct {
	Message string `json:"message"`
	Entity  Entity `json:"entity"`
}

type EntityListResponseDTO struct {
	Message  string   `json:"message"`
	Entities []Entity `json:"entities"`
}
