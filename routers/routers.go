package moduleRouter

import (
	"quest/controllers"
	"quest/models"
	"quest/repositories"
	"quest/services"
	"quest/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupModuleRoutes(r *gin.Engine, db *gorm.DB) {
	userRepo := repositories.NewUserRepository(db)
	groupRepo := repositories.NewGroupRepository(db)
	groupUserRepo := repositories.NewGroupUserRepository(db)

	userService := services.NewUserService(userRepo, groupRepo, groupUserRepo)

	userController := controllers.NewUserController(userService)

	rightService := services.NewRightService()
	rightController := controllers.NewRightController(rightService)

	transcation := r.Group("/")
	transcation.Use(utils.TransactionMiddleware(db))
	{
		transcation.POST("/users/login", utils.ValidationMiddleware(new(models.LoginRequestDTO)), userController.LoginUser)
		transcation.POST("/users/register", utils.ValidationMiddleware(new(models.UserCreateDTO)), userController.CreateUser)
		transcation.POST("/users/request-password-reset", utils.ValidationMiddleware(new(models.PasswordResetDTO)), userController.RequestPasswordReset)
		transcation.POST("/users/reset-password", utils.ValidationMiddleware(new(models.UpdatePasswordDTO)), userController.ResetPassword)
		transcation.PUT("/users/:id", utils.AuthMiddleware, utils.UserAuthMiddleware, userController.UpdateUser)
		transcation.PUT("/users/:id/change-password", utils.AuthMiddleware, utils.UserAuthMiddleware, userController.ChangePassword)
		transcation.GET("/auth/:provider/callback", utils.CompleteAuthHandler)
		transcation.GET("/users/:id/send-email", utils.AuthMiddleware, utils.UserAuthMiddleware, userController.SendEmail)
		transcation.GET("/users/:id/check-code", utils.AuthMiddleware, utils.UserAuthMiddleware, userController.CheckCode)
		transcation.POST("/rights", utils.AuthMiddleware, utils.ValidationMiddleware(new(models.RightCreateDTO)), rightController.CreateRight)
		transcation.GET("/rights/:id", utils.AuthMiddleware, utils.UserAuthMiddleware, rightController.GetRight)
		transcation.GET("/rights", utils.AuthMiddleware, rightController.GetAllRights)
		transcation.POST("/roles", utils.AuthMiddleware, utils.ValidationMiddleware(new(models.RoleCreateDTO)), rightController.CreateRole)
		transcation.GET("/roles/:id", utils.AuthMiddleware, rightController.GetRole)
		transcation.GET("/roles", utils.AuthMiddleware, rightController.GetAllRoles)
		transcation.PUT("/roles/:id", utils.AuthMiddleware, utils.GetIdMiddleware, utils.ValidationMiddleware(new(models.RoleUpdateDTO)), rightController.UpdateRole)
		transcation.DELETE("/roles/:id", utils.AuthMiddleware, utils.GetIdMiddleware, rightController.RemoveRole)
		transcation.DELETE("/rights/:id", utils.AuthMiddleware, utils.GetIdMiddleware, rightController.RemoveRight)
	}
	r.GET("/users/:id", utils.AuthMiddleware, utils.UserAuthMiddleware, userController.GetUser)
	r.GET("/whoiam", utils.AuthMiddleware, userController.WhoIam)
	r.GET("/auth/:provider", utils.BeginAuthHandler)
}
