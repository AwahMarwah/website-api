package router

import (
	"website-api/controller/auth"
	health_check "website-api/controller/health-check"
	"website-api/controller/role"
	"website-api/controller/user"
	"website-api/database"
	"website-api/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Run(db database.DB) (err error) {
	router := gin.Default()
	router.Use(cors.New(corsConfig))
	healthController := health_check.NewController(db.SqlDb)
	router.GET("/health", healthController.Check)

	authController := auth.NewController(db.GormDb)
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/forgot-password", authController.ForgotPassword)
		authGroup.POST("/reset-password", authController.ResetPassword)
		authGroup.POST("/resend-verification", authController.ResendVerification)
	}

	userController := user.NewController(db.GormDb)
	userGroup := router.Group("/user")
	{
		userGroup.POST("/sign-up", userController.SignUp)
		userGroup.GET("/verify-email", userController.VerifyEmailFromLink)
		userGroup.POST("/verify-email", userController.VerifyEmail)
		userGroup.POST("/sign-in", userController.SignIn)
		userGroup.GET("", middleware.AuthMiddleware(), userController.List)
		userGroup.GET("/:id", middleware.AuthMiddleware(), userController.Detail)
	}

	roleController := role.NewController(db.GormDb)
	roleGroup := router.Group("/role")
	{
		roleGroup.GET("", middleware.AuthMiddleware(), roleController.Find)
		roleGroup.POST("", middleware.AuthMiddleware(), roleController.Create)
		roleGroup.GET("/:id", middleware.AuthMiddleware(), roleController.Detail)
		roleGroup.PUT("/:id", middleware.AuthMiddleware(), roleController.Update)
		roleGroup.DELETE("/:id", middleware.AuthMiddleware(), roleController.Delete)
	}

	return router.Run()
}
