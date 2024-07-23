package controllers

import (
	"net/http"
	"quest/models"
	"quest/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RightController struct {
	RightService *services.RightService
}

func NewRightController(rightService *services.RightService) *RightController {
	return &RightController{
		RightService: rightService,
	}
}

func (ctrl *RightController) CreateRight(c *gin.Context) {
	validatedData := c.MustGet("dto").(*models.RightCreateDTO)
	tx := c.MustGet("tx").(*gorm.DB)

	right, err := ctrl.RightService.CreateRight(*validatedData, tx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Right created successfully", "right": right})
}

func (ctrl *RightController) GetRight(c *gin.Context) {
	id := c.Param("id")
	tx := c.MustGet("tx").(*gorm.DB)

	right, err := ctrl.RightService.GetRight(id, tx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"right": right})
}

func (ctrl *RightController) GetAllRights(c *gin.Context) {
	tx := c.MustGet("tx").(*gorm.DB)

	rights, err := ctrl.RightService.GetAllRights(tx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"rights": rights})
}

func (ctrl *RightController) CreateRole(c *gin.Context) {
	validatedData := c.MustGet("dto").(*models.RoleCreateDTO)
	tx := c.MustGet("tx").(*gorm.DB)

	role, err := ctrl.RightService.CreateRole(*validatedData, tx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role created successfully", "role": role})
}

func (ctrl *RightController) UpdateRole(c *gin.Context) {
	validatedData := c.MustGet("dto").(*models.RoleUpdateDTO)
	tx := c.MustGet("tx").(*gorm.DB)

	role, err := ctrl.RightService.UpdateRole(*validatedData, tx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role updated successfully", "role": role})
}

func (ctrl *RightController) GetRole(c *gin.Context) {
	id := c.Param("id")
	tx := c.MustGet("tx").(*gorm.DB)

	role, err := ctrl.RightService.GetRole(id, tx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"role": role})
}

func (ctrl *RightController) GetAllRoles(c *gin.Context) {
	tx := c.MustGet("tx").(*gorm.DB)

	roles, err := ctrl.RightService.GetAllRoles(tx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"roles": roles})
}

func (ctrl *RightController) AssignRightToRole(c *gin.Context) {
	roleID := c.Param("roleID")
	rightID := c.Param("rightID")
	tx := c.MustGet("tx").(*gorm.DB)

	err := ctrl.RightService.AssignRightToRole(roleID, rightID, tx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Right assigned to role successfully"})
}

func (ctrl *RightController) GetRoleForUser(c *gin.Context) {
	userID := c.Param("userID")
	tx := c.MustGet("tx").(*gorm.DB)

	roles, err := ctrl.RightService.GetRoleForUser(userID, tx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"roles": roles})
}

func (ctrl *RightController) AssignRoleToUser(c *gin.Context) {
	roleID := c.Param("roleID")
	userID := c.Param("userID")
	tx := c.MustGet("tx").(*gorm.DB)

	err := ctrl.RightService.AssignRoleToUser(roleID, userID, tx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role assigned to user successfully"})
}

func (ctrl *RightController) RemoveRoleFromUser(c *gin.Context) {
	roleID := c.Param("roleID")
	userID := c.Param("userID")
	tx := c.MustGet("tx").(*gorm.DB)

	err := ctrl.RightService.RemoveRoleFromUser(roleID, userID, tx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role removed from user successfully"})
}

func (ctrl *RightController) RemoveRole(c *gin.Context) {
	id := c.Param("id")
	tx := c.MustGet("tx").(*gorm.DB)

	err := ctrl.RightService.RemoveRole(id, tx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role removed successfully"})
}

func (ctrl *RightController) RemoveRight(c *gin.Context) {
	id := c.Param("id")
	tx := c.MustGet("tx").(*gorm.DB)

	err := ctrl.RightService.RemoveRight(id, tx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Right removed successfully"})
}
