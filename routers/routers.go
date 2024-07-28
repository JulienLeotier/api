package moduleRouter

import (
	"geniale/controllers"
	"geniale/models"
	"geniale/repositories"
	"geniale/services"
	"geniale/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupModuleRoutes(r *gin.Engine, db *gorm.DB) {
	userRepo := repositories.NewUserRepository(db)

	userService := services.NewUserService(userRepo)

	userController := controllers.NewUserController(userService)

	rightService := services.NewRightService()
	rightController := controllers.NewRightController(rightService)

	imageRepository := repositories.NewImageRepository(db)
	imageService := services.NewImageService(imageRepository)
	imageController := controllers.NewImageController(imageService)

	roomRepository := repositories.NewRoomRepository(db)
	roomService := services.NewRoomService(roomRepository, *userRepo)
	roomController := controllers.NewRoomController(roomService)

	roomRoleRepository := repositories.NewRoomRoleRepository(db)
	roomRoleService := services.NewRoomRoleService(roomRoleRepository)
	roomRoleController := controllers.NewRoomRoleController(roomRoleService)

	transcation := r.Group("/")
	transcation.Use(utils.TransactionMiddleware(db))
	{
		transcation.POST("/users/login", utils.JSONValidationMiddleware(new(models.LoginRequestDTO)), userController.LoginUser)
		transcation.POST("/users/register", utils.JSONValidationMiddleware(new(models.UserCreateDTO)), userController.CreateUser)
		transcation.POST("/users/request-password-reset", utils.JSONValidationMiddleware(new(models.PasswordResetDTO)), userController.RequestPasswordReset)
		transcation.POST("/users/reset-password", utils.JSONValidationMiddleware(new(models.UpdatePasswordDTO)), userController.ResetPassword)
		transcation.PUT("/users/:id", utils.AuthMiddleware, utils.UserAuthMiddleware, userController.UpdateUser)
		transcation.PUT("/users/:id/change-password", utils.AuthMiddleware, utils.UserAuthMiddleware, userController.ChangePassword)
		transcation.GET("/auth/:provider/callback", utils.CompleteAuthHandler)
		transcation.GET("/users/:id/send-email", utils.AuthMiddleware, utils.UserAuthMiddleware, userController.SendEmail)
		transcation.GET("/users/:id/check-code", utils.AuthMiddleware, utils.UserAuthMiddleware, userController.CheckCode)
		transcation.POST("/rights", utils.AuthMiddleware, utils.JSONValidationMiddleware(new(models.RightCreateDTO)), rightController.CreateRight)
		transcation.GET("/rights/:id", utils.AuthMiddleware, utils.UserAuthMiddleware, rightController.GetRight)
		transcation.GET("/rights", utils.AuthMiddleware, rightController.GetAllRights)
		transcation.POST("/roles", utils.AuthMiddleware, utils.JSONValidationMiddleware(new(models.RoleCreateDTO)), rightController.CreateRole)
		transcation.GET("/roles/:id", utils.AuthMiddleware, rightController.GetRole)
		transcation.GET("/roles", utils.AuthMiddleware, rightController.GetAllRoles)
		transcation.PUT("/roles/:id", utils.AuthMiddleware, utils.GetIdMiddleware, utils.JSONValidationMiddleware(new(models.RoleUpdateDTO)), rightController.UpdateRole)
		transcation.DELETE("/roles/:id", utils.AuthMiddleware, utils.GetIdMiddleware, rightController.RemoveRole)
		transcation.DELETE("/rights/:id", utils.AuthMiddleware, utils.GetIdMiddleware, rightController.RemoveRight)

		// Image routes
		transcation.POST("/images", utils.AuthMiddleware, imageController.UploadImages)
		transcation.GET("/images/:id", utils.AuthMiddleware, imageController.GetImage)
		transcation.GET("/images", utils.AuthMiddleware, imageController.GetImages)
		transcation.DELETE("/images/:id", utils.AuthMiddleware, imageController.RemoveImage)
		// Room routes
		transcation.POST("/rooms", utils.AuthMiddleware, utils.MultipartFormValidationMiddleware(new(models.RoomCreateDTO)), roomController.CreateRoom)
		transcation.GET("/rooms/:id", utils.AuthMiddleware, utils.GetIdMiddleware, roomController.GetRoom)
		transcation.DELETE("/rooms/:id", utils.AuthMiddleware, utils.GetIdMiddleware, roomController.DeleteRoom)
		transcation.PATCH("/rooms/:id", utils.AuthMiddleware, utils.MultipartFormValidationMiddleware(new(models.RoomCreateDTO)), utils.GetIdMiddleware, roomController.UpdateRoom)
		transcation.GET("/rooms", utils.AuthMiddleware, roomController.GetAllRooms)

		transcation.POST("/room_roles", utils.AuthMiddleware, utils.JSONValidationMiddleware(new(models.RoomRoleCreateDTO)), roomRoleController.CreateRoomRole)
		transcation.PUT("/room_roles/:id", utils.AuthMiddleware, utils.GetIdMiddleware, utils.JSONValidationMiddleware(new(models.RoomRoleCreateDTO)), roomRoleController.UpdateRoomRole)
		transcation.GET("/room_roles/:id", utils.AuthMiddleware, utils.GetIdMiddleware, roomRoleController.GetRoomRole)
		transcation.GET("/room_roles", utils.AuthMiddleware, roomRoleController.GetAllRoomRoles)
		transcation.DELETE("/room_roles/:id", utils.AuthMiddleware, utils.GetIdMiddleware, roomRoleController.DeleteRoomRole)
		transcation.GET("/room_roles/user/:user_id/room/:room_id", utils.AuthMiddleware, roomRoleController.GetRoomRoleByUserIDAndRoomID)

	}
	r.GET("/users/:id", utils.AuthMiddleware, utils.UserAuthMiddleware, userController.GetUser)
	r.GET("/whoiam", utils.AuthMiddleware, userController.WhoIam)
	r.GET("/auth/:provider", utils.BeginAuthHandler)
}
