package utils

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// JSONValidationMiddleware validates the JSON body and checks for missing or invalid data
func JSONValidationMiddleware(dto any) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Try to bind the JSON
		if err := c.ShouldBindJSON(dto); err != nil {
			if err.Error() == "EOF" {
				// Detect if the body is empty
				c.JSON(http.StatusBadRequest, gin.H{"error": "Request body cannot be empty"})
			} else {
				// General JSON binding error
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format: " + err.Error()})
			}
			c.Abort()
			return
		}

		// Validate the struct
		if err := validate.Struct(dto); err != nil {
			var errors []string
			for _, err := range err.(validator.ValidationErrors) {
				// Better error message formatting for validation errors
				errors = append(errors, fmt.Sprintf("%s failed validation on the '%s' tag", err.Field(), err.Tag()))
			}
			c.JSON(http.StatusBadRequest, gin.H{"errors": strings.Join(errors, ", ")})
			c.Abort()
			return
		}

		// If validation passes, store the validated DTO in the context
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
