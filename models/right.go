package models

type Right struct {
	ID   uint   `json:"id" gorm:"primary_key"`
	Name string `json:"name" gorm:"unique;not null"`
}

type Role struct {
	ID     uint    `json:"id" gorm:"primary_key"`
	Name   string  `json:"name" gorm:"unique;not null"`
	Rights []Right `json:"rights" gorm:"many2many:right_roles;"`
}

type RightRole struct {
	ID      uint  `json:"id" gorm:"primary_key"`
	RoleID  uint  `json:"role_id"`
	Role    Role  `json:"role" gorm:"foreignKey:RoleID"`
	Right   Right `json:"right" gorm:"foreignKey:RightID"`
	RightID uint  `json:"right_id"`
}

type RoleUser struct {
	ID     uint `json:"id" gorm:"primary_key"`
	RoleID uint `json:"role_id"`
	UserID uint `json:"user_id"`
}

type RightCreateDTO struct {
	Name string `json:"name" binding:"required"`
}

type RoleCreateDTO struct {
	Name   string `json:"name" binding:"required"`
	Rights []uint `json:"rights"`
}

type RoleUpdateDTO struct {
	ID     uint   `json:"id" binding:"required"`
	Name   string `json:"name" binding:"required"`
	Rights []uint `json:"rights"`
}
