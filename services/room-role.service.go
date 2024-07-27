package services

import (
	"geniale/models"
	"geniale/repositories"

	"gorm.io/gorm"
)

type RoomRoleService interface {
	CreateRoomRole(validationData models.RoomRole, tx *gorm.DB) (*models.RoomRole, error)
	GetRoomRoleByID(roomRoleID string, tx *gorm.DB) (*models.RoomRole, error)
	GetRoomRoleByUserID(userID string, tx *gorm.DB) ([]models.RoomRole, error)
	GetAllRoomRoles(tx *gorm.DB) ([]models.RoomRole, error)
	DeleteRoomRole(roomRoleID string, tx *gorm.DB) error
	UpdateRoomRole(roomRoleID string, validationData models.RoomRole, tx *gorm.DB) (*models.RoomRole, error)
}

type roomRoleService struct {
	repository repositories.RoomRoleRepository
}

func NewRoomRoleService(repo repositories.RoomRoleRepository) RoomRoleService {
	return &roomRoleService{repository: repo}
}

func (s *roomRoleService) CreateRoomRole(validationData models.RoomRole, tx *gorm.DB) (*models.RoomRole, error) {
	roomRole := &models.RoomRole{
		RoomID:   validationData.RoomID,
		UserID:   validationData.UserID,
		Role:     validationData.Role,
		Name:     validationData.Name,
		Job:      validationData.Job,
		Region:   validationData.Region,
		Passion:  validationData.Passion,
		Anecdote: validationData.Anecdote,
	}
	roomRoleSave, err := s.repository.CreateRoomRole(roomRole, tx)
	if err != nil {
		return nil, err
	}
	return roomRoleSave, nil
}

func (s *roomRoleService) GetRoomRoleByID(roomRoleID string, tx *gorm.DB) (*models.RoomRole, error) {
	return s.repository.GetRoomRoleByID(roomRoleID, tx)
}

func (s *roomRoleService) GetRoomRoleByUserID(userID string, tx *gorm.DB) ([]models.RoomRole, error) {
	return s.repository.GetRoomRoleByUserID(userID, tx)
}

func (s *roomRoleService) GetAllRoomRoles(tx *gorm.DB) ([]models.RoomRole, error) {
	return s.repository.GetAllRoomRoles(tx)
}

func (s *roomRoleService) DeleteRoomRole(roomRoleID string, tx *gorm.DB) error {
	return s.repository.DeleteRoomRole(roomRoleID, tx)
}

func (s *roomRoleService) UpdateRoomRole(roomRoleID string, validationData models.RoomRole, tx *gorm.DB) (*models.RoomRole, error) {
	roomRole, err := s.repository.GetRoomRoleByID(roomRoleID, tx)
	if err != nil {
		return nil, err
	}

	roomRole.RoomID = validationData.RoomID
	roomRole.UserID = validationData.UserID
	roomRole.Role = validationData.Role
	roomRole.Name = validationData.Name
	roomRole.Job = validationData.Job
	roomRole.Region = validationData.Region
	roomRole.Passion = validationData.Passion
	roomRole.Anecdote = validationData.Anecdote

	return s.repository.UpdateRoomRole(roomRole, tx)
}
