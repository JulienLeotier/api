package utils

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		c.Abort()
		return
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := ValidateToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	c.Set("id", token.ID)
	c.Next()
}

func UserAuthMiddleware(c *gin.Context) {
	paramID := c.Param("id")

	if paramID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	c.Set("userID", paramID)
	c.Next()
}

func GetIdMiddleware(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		c.Abort()
		return
	}
	c.Set("id", id)
	c.Next()
}
