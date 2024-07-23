package models

import "time"

type Patient struct {
	ID          uint `gorm:"primaryKey"`
	UserID      uint `gorm:"not null"`
	User        User
	Name        string `gorm:"size:255;not null"`
	DOB         time.Time
	Gender      string   `gorm:"size:10"`
	ContactInfo string   `gorm:"type:jsonb"`
	Entity      []Entity `gorm:"many2many:entity_patients;"`
	Suspended   bool     `gorm:"default:false"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
