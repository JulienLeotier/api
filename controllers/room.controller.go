package controllers

import (
	"geniale/models"
	"geniale/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RoomController struct {
	services services.RoomService
}

func NewRoomController(roomService services.RoomService) *RoomController {
	return &RoomController{
		services: roomService,
	}
}

func (ctrl *RoomController) CreateRoom(ctx *gin.Context) {
	var validatedData models.RoomCreateDTO

	if err := ctx.ShouldBind(&validatedData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := ctx.MustGet("tx").(*gorm.DB)

	room, err := ctrl.services.CreateRoom(validatedData, tx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"data": room})
}

func (ctrl *RoomController) GetRoom(ctx *gin.Context) {
	roomID := ctx.MustGet("id").(string)

	tx := ctx.MustGet("tx").(*gorm.DB)

	room, err := ctrl.services.GetRoomByID(roomID, tx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": room})
}

func (ctrl *RoomController) GetAllRooms(ctx *gin.Context) {
	tx := ctx.MustGet("tx").(*gorm.DB)

	rooms, err := ctrl.services.GetAllRooms(tx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": rooms})
}

func (ctrl *RoomController) DeleteRoom(ctx *gin.Context) {
	roomID := ctx.MustGet("id").(string)

	tx := ctx.MustGet("tx").(*gorm.DB)

	if err := ctrl.services.DeleteRoom(roomID, tx); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Room deleted successfully"})
}

func (ctrl *RoomController) UpdateRoom(ctx *gin.Context) {
	roomID := ctx.MustGet("id").(string)

	var validatedData models.RoomCreateDTO
	if err := ctx.ShouldBind(&validatedData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := ctx.MustGet("tx").(*gorm.DB)

	room, err := ctrl.services.UpdateRoom(roomID, validatedData, tx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": room})
}
