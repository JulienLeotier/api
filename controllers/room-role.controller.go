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
	validatedData := ctx.MustGet("dto").(*models.RoomRoleCreateDTO)

	tx := ctx.MustGet("tx").(*gorm.DB)

	roomRole, err := ctrl.services.CreateRoomRole(*validatedData, tx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"data": roomRole})
}

func (ctrl *RoomRoleController) GetRoomRole(ctx *gin.Context) {
	id := ctx.MustGet("id").(string)

	tx := ctx.MustGet("tx").(*gorm.DB)

	roomRole, err := ctrl.services.GetRoomRoleByID(id, tx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": roomRole})
}

func (ctrl *RoomRoleController) GetRoomRoleByUserID(ctx *gin.Context) {
	userID := ctx.Param("id")

	tx := ctx.MustGet("tx").(*gorm.DB)

	roomRoles, err := ctrl.services.GetRoomRoleByUserID(userID, tx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": roomRoles})
}

func (ctrl *RoomRoleController) GetRoomRoleByRoomID(ctx *gin.Context) {
	roomID := ctx.Param("id")

	tx := ctx.MustGet("tx").(*gorm.DB)

	roomRoles, err := ctrl.services.GetRoomRoleByRoomID(roomID, tx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": roomRoles})
}

func (ctrl *RoomRoleController) GetRoomRoleByUserIDAndRoomID(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	roomID := ctx.Param("room_id")

	tx := ctx.MustGet("tx").(*gorm.DB)

	roomRoles, err := ctrl.services.GetRoomRoleByUserIDAndRoomID(userID, roomID, tx)
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
	id := ctx.MustGet("id").(string)

	tx := ctx.MustGet("tx").(*gorm.DB)

	if err := ctrl.services.DeleteRoomRole(id, tx); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Room role deleted successfully"})
}

func (ctrl *RoomRoleController) UpdateRoomRole(ctx *gin.Context) {
	id := ctx.MustGet("id").(string)
	validatedData := ctx.MustGet("dto").(*models.RoomRoleCreateDTO)

	tx := ctx.MustGet("tx").(*gorm.DB)

	roomRole, err := ctrl.services.UpdateRoomRole(id, *validatedData, tx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": roomRole})
}
