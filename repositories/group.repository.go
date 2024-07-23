package repositories

import (
	models "geniale/models"

	"gorm.io/gorm"
)

type GroupRepository struct {
	DB *gorm.DB
}

func NewGroupRepository(db *gorm.DB) *GroupRepository {
	return &GroupRepository{DB: db}
}

func (r *GroupRepository) FindByName(name string) (*models.Group, error) {
	var group models.Group
	result := r.DB.Where("name = ?", name).First(&group)
	return &group, result.Error
}

func (r *GroupRepository) FindGroupNameByUserID(userID uint) (string, error) {
	var groupName string
	result := r.DB.Table("groups").Select("name").Joins("JOIN group_users ON groups.id = group_users.group_id").Where("group_users.user_id = ?", userID).Scan(&groupName)
	return groupName, result.Error
}
