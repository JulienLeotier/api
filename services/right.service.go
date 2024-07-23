package services

import (
	"errors"
	"quest/models"
	"strconv"

	"gorm.io/gorm"
)

type RightService struct {
}

func NewRightService() *RightService {
	return &RightService{}
}

func (s *RightService) CreateRight(right models.RightCreateDTO, tx *gorm.DB) (*models.Right, error) {
	newRight := &models.Right{
		Name: right.Name,
	}

	if err := tx.Create(newRight).Error; err != nil {
		return nil, err
	}

	return newRight, nil

}

func (s *RightService) GetRight(id string, tx *gorm.DB) (*models.Right, error) {
	var right models.Right
	if err := tx.Where("id = ?", id).First(&right).Error; err != nil {
		return nil, err
	}
	return &right, nil
}

func (s *RightService) GetAllRights(tx *gorm.DB) ([]models.Right, error) {
	var rights []models.Right
	if err :=
		tx.Find(&rights).Error; err != nil {
		return nil, err
	}
	return rights, nil
}

// CreateRole creates a new role
func (s *RightService) CreateRole(role models.RoleCreateDTO, tx *gorm.DB) (*models.Role, error) {
	newRole := &models.Role{
		Name: role.Name,
	}

	if err := tx.Create(newRole).Error; err != nil {
		return nil, err
	}
	if len(role.Rights) > 0 {
		for _, rightID := range role.Rights {
			right := models.Right{}
			if err := tx.Where("id = ?", rightID).First(&right).Error; err != nil {
				return nil, err
			}
			if err := tx.Model(&newRole).Association("Rights").Append(&right); err != nil {
				return nil, err
			}
			newRole.Rights = append(newRole.Rights, right)
		}
	}

	return newRole, nil
}
func (s *RightService) RemoveRightFromRole(roleID, rightID string, tx *gorm.DB) error {
	role := models.Role{}
	right := models.Right{}
	if err := tx.Where("id = ?", roleID).First(&role).Error; err != nil {
		return err
	}

	if err := tx.Where("id = ?", rightID).First(&right).Error; err != nil {
		return err
	}

	if err := tx.Model(&role).Association("Rights").Delete(&right); err != nil {
		return err
	}

	return nil
}

func (s *RightService) UpdateRole(roleUpdateDTO models.RoleUpdateDTO, tx *gorm.DB) (*models.Role, error) {
	var role models.Role
	if err := tx.Preload("Rights").First(&role, roleUpdateDTO.ID).Error; err != nil {
		return nil, errors.New("role not found")
	}

	// Update the role's name
	role.Name = roleUpdateDTO.Name

	// Update the associated rights
	var rights []models.Right
	if len(roleUpdateDTO.Rights) > 0 {
		if err := tx.Where("id IN ?", roleUpdateDTO.Rights).Find(&rights).Error; err != nil {
			return nil, errors.New("some rights not found")
		}
	}

	// Update the role with the new rights
	if err := tx.Model(&role).Association("Rights").Replace(rights); err != nil {
		return nil, errors.New("failed to update rights")
	}

	// Save the updated role
	if err := tx.Save(&role).Error; err != nil {
		return nil, errors.New("failed to save role")
	}

	return &role, nil
}

func (s *RightService) GetRole(id string, tx *gorm.DB) (*models.Role, error) {
	var role models.Role
	if err := tx.Where("id = ?", id).First(&role).Error; err != nil {
		return nil, err
	}
	if err := tx.Model(&role).Association("Rights").Find(&role.Rights); err != nil {
		return nil, err
	}
	return &role, nil
}

func (s *RightService) GetAllRoles(tx *gorm.DB) ([]models.Role, error) {
	var roles []models.Role
	if err := tx.Preload("Rights").Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (s *RightService) AssignRightToRole(roleID, rightID string, tx *gorm.DB) error {
	role := models.Role{}
	right := models.Right{}
	if err := tx.Where("id = ?", roleID).First(&role).Error; err != nil {
		return err
	}

	if err := tx.Where("id = ?", rightID).First(&right).Error; err != nil {
		return err
	}

	if err := tx.Model(&role).Association("Rights").Append(&right); err != nil {
		return err
	}

	return nil
}

func (s *RightService) AssignRoleToUser(roleID, userID string, tx *gorm.DB) error {
	role := models.Role{}
	if err := tx.Where("id = ?", roleID).First(&role).Error; err != nil {
		return err
	}

	userIDUint, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		return err
	}

	if err := tx.Model(&role).Association("Users").Append(&models.User{ID: uint(userIDUint)}); err != nil {
		return err
	}

	return nil
}

func (s *RightService) GetRoleForUser(userID string, tx *gorm.DB) ([]models.Role, error) {
	var roles []models.Role
	if err := tx.Where("user_id = ?", userID).Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (s *RightService) RemoveRoleFromUser(roleID, userID string, tx *gorm.DB) error {
	role := models.Role{}
	if err := tx.Where("id = ?", roleID).First(&role).Error; err != nil {
		return err
	}

	userIDUint, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		return err
	}

	if err := tx.Model(&role).Association("Users").Delete(&models.User{ID: uint(userIDUint)}); err != nil {
		return err
	}

	return nil
}

func (s *RightService) RemoveRole(roleID string, tx *gorm.DB) error {
	role := models.Role{}
	if err := tx.Where("id = ?", roleID).First(&role).Error; err != nil {
		return err
	}
	if err := tx.Model(&role).Association("Rights").Clear(); err != nil {
		return err
	}

	if err := tx.Delete(&role).Error; err != nil {
		return err
	}

	return nil
}

func (s *RightService) RemoveRight(rightID string, tx *gorm.DB) error {
	right := models.Right{}
	if err := tx.Where("id = ?", rightID).First(&right).Error; err != nil {
		return err
	}
	if err := tx.Model(&models.RightRole{}).Where("right_id = ?", rightID).Delete(&models.RightRole{}).Error; err != nil {
		return err
	}

	if err := tx.Delete(&right).Error; err != nil {
		return err
	}

	return nil
}
