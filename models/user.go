package models

import (
	"time"
)

type UserType int

const (
	EnumUserTypeEmployee  = 0
	EnumUserTypeUser      = 1
	EnumUserTypeHealthPro = 2
	EnumUserTypePatient   = 3
)

type User struct {
	ID                 uint       `json:"id" gorm:"primary_key"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	DeletedAt          *time.Time `json:"deleted_at" sql:"index"`
	Username           string     `json:"username" gorm:"unique;not null"`
	Email              string     `json:"email" gorm:"unique;not null"`
	Password           string     `json:"-"`
	Phone              string     `json:"phone"`
	PasswordResetToken string     `json:"-"`
	IsGoogleUser       bool       `json:"is_google_user" gorm:"default:false"`
	EmailChecked       bool       `json:"email_checked" gorm:"default:false"`
	Type               UserType   `sql:"type:ENUM('0', '1', '2', '3');default:'1'" json:"type"`
}
type UserCode struct {
	ID        uint       `json:"id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" sql:"index"`
	UserID    uint       `json:"user_id"`
	Code      string     `json:"code"`
}
type LoginRequestDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

type UserCreateDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=100"`
	Username string `json:"username" binding:"required"`
}

type UserCreateTmpDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=100"`
	Username string `json:"username" binding:"required"`
}

type PasswordResetDTO struct {
	Email string `json:"email" validate:"required,email"`
}

type UpdatePasswordDTO struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8,max=100"`
}

type GroupUser struct {
	GroupID uint `json:"group_id" gorm:"primary_key"`
	UserID  uint `json:"user_id" gorm:"primary_key"`
}

type Group struct {
	ID        uint       `json:"id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" sql:"index"`
	Name      string     `json:"name" gorm:"unique;not null"`
}

type LoginUserResponseDTO struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}

type CreateUserResponseDTO struct {
	Message   string    `json:"message"`
	User      User      `json:"user"`
	Group     Group     `json:"group"`
	GroupUser GroupUser `json:"group_user"`
}

type UserUpdateDTO struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

type EmailDTO struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type ChangePasswordDTO struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8,max=100"`
}
