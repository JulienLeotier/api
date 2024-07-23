package repositories

import (
	models "geniale/models"

	"gorm.io/gorm"
)

type GroupUserRepository struct {
	DB *gorm.DB
}

func NewGroupUserRepository(db *gorm.DB) *GroupUserRepository {
	return &GroupUserRepository{DB: db}
}

func (r *GroupUserRepository) Create(groupUser *models.GroupUser) error {
	return r.DB.Create(groupUser).Error
}
