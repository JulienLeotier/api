package controllers

import (
	"geniale/services"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type ImageController struct {
	service services.ImageService
}

func NewImageController(service services.ImageService) *ImageController {
	return &ImageController{service: service}
}

func (c *ImageController) UploadImage(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No file is received"})
		return
	}

	filename := filepath.Base(file.Filename)
	if err := ctx.SaveUploadedFile(file, "./uploads/"+filename); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save the file"})
		return
	}

	image, err := c.service.UploadImage("uploads/" + filename)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"image": image})
}

func (c *ImageController) GetImage(ctx *gin.Context) {
	id := ctx.Param("id")
	image, err := c.service.GetImage(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"image": image})
}
