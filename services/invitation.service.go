package services

import (
	"quest/models"

	"gorm.io/gorm"
)

type InvitationService struct {
}

func NewInvitationService() *InvitationService {
	return &InvitationService{}
}

func (service *InvitationService) CreateInvitation(data models.InvitationCreateDTO, tx *gorm.DB) (*models.Invitation, error) {
	invitation := models.Invitation{
		Email: data.Email,
	}
	if err := tx.Create(&invitation).Error; err != nil {
		return nil, err
	}

	return &invitation, nil
}
