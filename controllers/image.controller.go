package controllers

import (
	"geniale/models"
	"geniale/services"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type ImageController struct {
	service services.ImageService
}

func NewImageController(service services.ImageService) *ImageController {
	return &ImageController{service: service}
}

func (c *ImageController) UploadImages(ctx *gin.Context) {
	file, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No files are received"})
		return
	}

	files := file.File["file"]
	if len(files) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No files are received"})
		return
	}

	var images []services.ImageDTO

	for _, file := range files {
		filename := filepath.Base(file.Filename)
		filePath := "./uploads/" + filename
		if err := ctx.SaveUploadedFile(file, filePath); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save the file"})
			return
		}

		image, err := c.service.UploadImage(filePath)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		images = append(images, *image)
	}

	ctx.JSON(http.StatusOK, gin.H{"images": images})
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

func (c *ImageController) GetImages(ctx *gin.Context) {
	images, err := c.service.GetImages()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"images": images})
}

func (c *ImageController) RemoveImage(ctx *gin.Context) {
	id := ctx.Param("id")
	image, err := c.service.GetImage(id)
	if err != nil {
		if err.Error() == "record not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := os.Remove("./" + image.FilePath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove image"})
		return
	}
	if err := c.service.RemoveImage(models.Image(*image)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Image removed successfully"})
}
