package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"geniale/models"
	"geniale/services"
	"geniale/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserController struct {
	UserService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		UserService: userService,
	}
}
func (ctrl *UserController) WhoIam(c *gin.Context) {
	id := c.MustGet("id").(string)
	user, err := ctrl.UserService.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to find user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (ctrl *UserController) CreateUser(c *gin.Context) {
	validatedData := c.MustGet("dto").(*models.UserCreateDTO)
	tx := c.MustGet("tx").(*gorm.DB)

	user, err := ctrl.UserService.CreateUser(*validatedData, tx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully", "user": user})
}

func (ctrl *UserController) LoginUser(c *gin.Context) {
	validatedData := c.MustGet("dto").(*models.LoginRequestDTO)

	token, err := ctrl.UserService.LoginUser(*validatedData)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User logged in successfully", "token": token})
}

func (ctrl *UserController) GetUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid ID format: %v", err)})
		return
	}

	user, err := ctrl.UserService.GetUser(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (ctrl *UserController) RequestPasswordReset(c *gin.Context) {
	tx := c.MustGet("tx").(*gorm.DB)
	dto := c.MustGet("dto").(*models.PasswordResetDTO)
	user, err := ctrl.UserService.FindByEmail(dto.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user"})
		return
	}

	token, err := utils.GeneratePasswordResetToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	user.PasswordResetToken = token
	if err := ctrl.UserService.UpdateUser(user, tx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save token"})
		return
	}
	templatePath := os.Getenv("TEMPLATE_DIR") + os.Getenv("TEMPLATE_RESET_PASSWORD")
	emailData := utils.EmailData{
		Name:  user.Username,
		Token: token,
		Url:   os.Getenv("FRONT") + "/reset-password?token=",
	}
	htmlContent, err := utils.RenderTemplate(templatePath, emailData)
	if err != nil {
		log.Printf("Could not render template: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to render email template"})
		return
	}
	subject := "Password reset"
	err = utils.SendMail(user.Email, subject, htmlContent)
	if err != nil {
		log.Printf("Could not send email: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset email sent"})
}

func (ctrl *UserController) ResetPassword(c *gin.Context) {
	dto := c.MustGet("dto").(*models.UpdatePasswordDTO)
	tx := c.MustGet("tx").(*gorm.DB)

	var user models.User
	if err := ctrl.UserService.FindByPasswordResetToken(dto.Token, &user); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid token"})
		return
	}

	hashedPassword, err := utils.HashPassword(dto.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
	}
	user.Password = hashedPassword
	user.PasswordResetToken = ""
	if err := ctrl.UserService.UpdateUser(&user, tx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset password"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}

func (ctrl *UserController) UpdateUser(c *gin.Context) {
	tx := c.MustGet("tx").(*gorm.DB)
	id := c.GetString("userID")
	idUser := c.GetString("id")

	if id != idUser {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only modify your own information"})
		return
	}

	var userUpdate models.UserUpdateDTO
	if err := c.ShouldBindJSON(&userUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := ctrl.UserService.GetUserByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user"})
		return
	}

	if userUpdate.Email != "" {
		if !user.IsGoogleUser {
			user.Email = userUpdate.Email
			user.EmailChecked = false
		}
	}

	if userUpdate.Username != "" {
		user.Username = userUpdate.Username
	}
	if userUpdate.Phone != "" {
		user.Phone = userUpdate.Phone
	}
	if err := ctrl.UserService.ExistingUser(*user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email or username already exists"})
		return
	}

	if err := ctrl.UserService.UpdateUser(user, tx); err != nil {
		if err == gorm.ErrDuplicatedKey {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email or username already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func (ctrl *UserController) ChangePassword(c *gin.Context) {
	tx := c.MustGet("tx").(*gorm.DB)
	id := c.GetString("userID")
	idUser := c.GetString("id")

	if id != idUser {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only modify your own information"})
		return
	}

	var pwdChange models.ChangePasswordDTO
	if err := c.ShouldBindJSON(&pwdChange); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(id)
	user, err := ctrl.UserService.GetUserByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pwdChange.CurrentPassword)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Current password is incorrect"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pwdChange.NewPassword), 14)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user.Password = string(hashedPassword)
	if err := ctrl.UserService.UpdateUser(user, tx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

func (ctrl *UserController) SendEmail(c *gin.Context) {
	tx := c.MustGet("tx").(*gorm.DB)
	idUser := c.GetString("userID")
	id := c.GetString("id")

	if id != idUser {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only modify your own information"})
		return
	}

	var emailData models.EmailDTO

	user, err := ctrl.UserService.GetUserByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user"})
		return
	}
	if user.IsGoogleUser {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Google users cannot send emails"})
		return
	}
	codeData, err := ctrl.UserService.GetCode(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get code"})
		return

	}
	emailData.Name = user.Username
	if codeData.Code == "" {
		code := utils.GenerateRandomCode()
		emailData.Code = code
		err = ctrl.UserService.SaveUserCode(user, code, tx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save code"})
			return
		}
	} else {
		emailData.Code = codeData.Code
	}

	templatePath := os.Getenv("TEMPLATE_DIR") + os.Getenv("TEMPLATE_EMAIL_CODE")
	htmlContent, err := utils.RenderTemplate(templatePath, emailData)
	if err != nil {
		log.Printf("Could not render template: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to render email template"})
		return
	}
	subject := "Verification code"
	err = utils.SendMail(user.Email, subject, htmlContent)
	if err != nil {
		log.Printf("Could not send email: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email sent successfully"})
}

func (ctrl *UserController) CheckCode(c *gin.Context) {
	tx := c.MustGet("tx").(*gorm.DB)
	idUser := c.GetString("userID")
	id := c.GetString("id")
	code := c.Query("code")

	if id != idUser {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only check code for your own account"})
		return
	}
	user, err := ctrl.UserService.GetUserByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user"})
		return
	}
	if user.IsGoogleUser {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Google users don't have verification codes"})
		return
	}
	err = ctrl.UserService.CheckCode(user, code, tx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid code"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Code is valid"})
}

func (ctrl *UserController) GetAllUsers(c *gin.Context) {
	users, err := ctrl.UserService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}
