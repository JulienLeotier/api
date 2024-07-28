package utils

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func JSONValidationMiddleware(dto any) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := c.ShouldBindJSON(dto); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		if err := validate.Struct(dto); err != nil {
			var errors []string
			for _, err := range err.(validator.ValidationErrors) {
				errors = append(errors, fmt.Sprintf("%s failed on the '%s' tag", err.Field(), err.Tag()))
			}
			c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
			c.Abort()
			return
		}

		c.Set("dto", dto)
		c.Next()
	}
}

func MultipartFormValidationMiddleware(dto any) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := c.ShouldBind(dto); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		if err := validate.Struct(dto); err != nil {
			var errors []string
			for _, err := range err.(validator.ValidationErrors) {
				errors = append(errors, fmt.Sprintf("%s failed on the '%s' tag", err.Field(), err.Tag()))
			}
			c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
			c.Abort()
			return
		}

		c.Set("dto", dto)
		c.Next()
	}
}
