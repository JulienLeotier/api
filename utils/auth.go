package utils

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"os"
	"strconv"

	models "geniale/models"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func BeginAuthHandler(c *gin.Context) {
	q := c.Request.URL.Query()
	q.Add("provider", c.Param("provider"))
	c.Request.URL.RawQuery = q.Encode()
	url, err := gothic.GetAuthURL(c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get auth URL"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"url": url})
}

func generateUsername(email string) string {
	return email
}

func CompleteAuthHandler(c *gin.Context) {
	tx := c.MustGet("tx").(*gorm.DB)
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		println(gin.H{"error": "Failed to complete user auth"})
		c.Redirect(http.StatusTemporaryRedirect, os.Getenv("FRONT")+"/")
	}

	// Save the user to the database if they don't already exist
	var existingUser models.User
	if err := DB.Where("email = ?", user.Email).First(&existingUser).Error; err == gorm.ErrRecordNotFound {
		// Generate a username if not provided
		username := user.Name
		if username == "" {
			username = generateUsername(user.Email)
		}
		// add here the email of the user you want to allow to connect
		if user.Email != "julien.leotier@gmail.com" {
			println(gin.H{"error": "Failed to complete user auth"})
			c.Redirect(http.StatusTemporaryRedirect, os.Getenv("FRONT")+"/")
		}
		newUser := models.User{
			Username:     username,
			Email:        user.Email,
			Password:     "", // No password needed for OAuth users
			IsGoogleUser: true,
			EmailChecked: true,
			Type:         models.EnumUserTypeEmployee,
		}
		if err := DB.Create(&newUser).Error; err != nil {
			println(gin.H{"error": "Failed to create user"})
			tx.Rollback()
			c.Redirect(http.StatusTemporaryRedirect, os.Getenv("FRONT")+"/")
		}
		existingUser = newUser
	}

	// Generate JWT token for the user
	idToString := strconv.Itoa(int(existingUser.ID))
	token, err := GenerateToken(existingUser.Email, "user", idToString)
	if err != nil {
		println(gin.H{"error": "Failed to generate token"})
		tx.Rollback()
		c.Redirect(http.StatusTemporaryRedirect, os.Getenv("FRONT")+"/")
	}
	frontEndURL := os.Getenv("FRONT") + "/home?token=" + token
	println(frontEndURL)
	c.Redirect(http.StatusTemporaryRedirect, frontEndURL)
}

func GeneratePasswordResetToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
