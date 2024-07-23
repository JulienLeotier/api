package controllers

import (
	"net/http"
	"quest/models"
	"quest/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type InvitationController struct {
	InvitationService *services.InvitationService
	UserService       *services.UserService
}

func NewInvitationController(invitationService *services.InvitationService) *InvitationController {
	return &InvitationController{
		InvitationService: invitationService,
	}
}

func (ctrl *InvitationController) CreateInvitation(c *gin.Context) {
	validatedData := c.MustGet("dto").(*models.InvitationCreateDTO)
	tx := c.MustGet("tx").(*gorm.DB)

	invitation, err := ctrl.InvitationService.CreateInvitation(*validatedData, tx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := models.UserCreateTmpDTO{
		Email:    validatedData.Email,
		Password: "",
		Username: validatedData.Email,
	}

	tmpUser, err := ctrl.UserService.CreateTmpUser(user, tx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	invitationResponse := models.InvitationResponseDTO{
		Message: "Invitation created successfully",
		Invitation: models.InvitationCreateResponseDTO{
			Invitation: *invitation,
			User:       *tmpUser,
		},
	}
	c.JSON(http.StatusCreated, invitationResponse)
}
