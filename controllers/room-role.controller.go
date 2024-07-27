package controllers

import (
	"geniale/models"
	"geniale/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RoomRoleController struct {
	services services.RoomRoleService
}

func NewRoomRoleController(roomRoleService services.RoomRoleService) *RoomRoleController {
	return &RoomRoleController{
		services: roomRoleService,
	}
}

func (ctrl *RoomRoleController) CreateRoomRole(ctx *gin.Context) {
	var validatedData models.RoomRole

	if err := ctx.ShouldBind(&validatedData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := ctx.MustGet("tx").(*gorm.DB)

	roomRole, err := ctrl.services.CreateRoomRole(validatedData, tx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"data": roomRole})
}

func (ctrl *RoomRoleController) GetRoomRole(ctx *gin.Context) {
	roomRoleID := ctx.MustGet("id").(string)

	tx := ctx.MustGet("tx").(*gorm.DB)

	roomRole, err := ctrl.services.GetRoomRoleByID(roomRoleID, tx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": roomRole})
}

func (ctrl *RoomRoleController) GetRoomRoleByUserID(ctx *gin.Context) {
	userID := ctx.MustGet("userID").(string)
	idUser := ctx.GetString("id")
	if userID != idUser {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You can only view your own information"})
		return
	}
	tx := ctx.MustGet("tx").(*gorm.DB)

	roomRoles, err := ctrl.services.GetRoomRoleByUserID(userID, tx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": roomRoles})
}

func (ctrl *RoomRoleController) GetAllRoomRoles(ctx *gin.Context) {
	tx := ctx.MustGet("tx").(*gorm.DB)

	roomRoles, err := ctrl.services.GetAllRoomRoles(tx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": roomRoles})
}

func (ctrl *RoomRoleController) DeleteRoomRole(ctx *gin.Context) {
	roomRoleID := ctx.MustGet("id").(string)

	tx := ctx.MustGet("tx").(*gorm.DB)

	if err := ctrl.services.DeleteRoomRole(roomRoleID, tx); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Room role deleted successfully"})
}

func (ctrl *RoomRoleController) UpdateRoomRole(ctx *gin.Context) {
	roomRoleID := ctx.MustGet("id").(string)

	var validatedData models.RoomRole
	if err := ctx.ShouldBind(&validatedData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := ctx.MustGet("tx").(*gorm.DB)

	roomRole, err := ctrl.services.UpdateRoomRole(roomRoleID, validatedData, tx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": roomRole})
}
