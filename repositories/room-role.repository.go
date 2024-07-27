package repositories

import (
	"geniale/models"

	"gorm.io/gorm"
)

type RoomRoleRepository interface {
	CreateRoomRole(roomRole *models.RoomRole, tx *gorm.DB) (*models.RoomRole, error)
	GetRoomRoleByID(roomRoleID string, tx *gorm.DB) (*models.RoomRole, error)
	GetRoomRoleByUserID(userID string, tx *gorm.DB) ([]models.RoomRole, error)
	GetAllRoomRoles(tx *gorm.DB) ([]models.RoomRole, error)
	DeleteRoomRole(roomRoleID string, tx *gorm.DB) error
	UpdateRoomRole(roomRole *models.RoomRole, tx *gorm.DB) (*models.RoomRole, error)
}

type roomRoleRepository struct {
	DB *gorm.DB
}

func NewRoomRoleRepository(db *gorm.DB) RoomRoleRepository {
	return &roomRoleRepository{DB: db}
}

func (r *roomRoleRepository) CreateRoomRole(roomRole *models.RoomRole, tx *gorm.DB) (*models.RoomRole, error) {
	if err := tx.Create(roomRole).Error; err != nil {
		return nil, err
	}
	return roomRole, nil
}

func (r *roomRoleRepository) GetRoomRoleByID(roomRoleID string, tx *gorm.DB) (*models.RoomRole, error) {
	var roomRole models.RoomRole
	if err := tx.First(&roomRole, roomRoleID).Error; err != nil {
		return nil, err
	}
	return &roomRole, nil
}

func (r *roomRoleRepository) GetRoomRoleByUserID(userID string, tx *gorm.DB) ([]models.RoomRole, error) {
	var roomRoles []models.RoomRole
	if err := tx.Where("user_id = ?", userID).Find(&roomRoles).Error; err != nil {
		return nil, err
	}
	return roomRoles, nil
}

func (r *roomRoleRepository) GetAllRoomRoles(tx *gorm.DB) ([]models.RoomRole, error) {
	var roomRoles []models.RoomRole
	if err := tx.Find(&roomRoles).Error; err != nil {
		return nil, err
	}
	return roomRoles, nil
}

func (r *roomRoleRepository) DeleteRoomRole(roomRoleID string, tx *gorm.DB) error {
	if err := tx.Delete(&models.RoomRole{}, roomRoleID).Error; err != nil {
		return err
	}
	return nil
}

func (r *roomRoleRepository) UpdateRoomRole(roomRole *models.RoomRole, tx *gorm.DB) (*models.RoomRole, error) {
	if err := tx.Save(roomRole).Error; err != nil {
		return nil, err
	}
	return roomRole, nil
}
