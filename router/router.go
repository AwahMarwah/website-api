package router

import (
	"website-api/controller/auth"
	"website-api/controller/brand"
	"website-api/controller/category"
	content_page "website-api/controller/content-page"
	health_check "website-api/controller/health-check"
	"website-api/controller/product"
	"website-api/controller/role"
	"website-api/controller/user"
	"website-api/database"
	"website-api/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func Run(db database.DB, redis *redis.Client) (err error) {
	router := gin.Default()
	router.Use(middleware.NgrokSkipWarning())
	router.Use(cors.New(corsConfig))

	// PUBLIC API
	healthController := health_check.NewController(db.SqlDb)
	router.GET("/health", healthController.Check)

	authController := auth.NewController(db.GormDb)
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/forgot-password", authController.ForgotPassword)
		authGroup.POST("/reset-password", authController.ResetPassword)
		authGroup.POST("/resend-verification", authController.ResendVerification)
	}

	contentPageController := content_page.NewController(db.GormDb, redis)
	contentPageGroup := router.Group("/content-page")
	{
		// Public
		contentPageGroup.GET("/pages/:slug", contentPageController.GetBySlug)
		contentPageGroup.GET("faqs", contentPageController.GetFaq)
	}

	userController := user.NewController(db.GormDb)
	userGroup := router.Group("/user")
	{
		// PUBLIC
		userGroup.POST("/sign-up", userController.SignUp)
		userGroup.GET("/verify-email", userController.VerifyEmailFromLink)
		userGroup.POST("/verify-email", userController.VerifyEmail)
		userGroup.POST("/sign-in", userController.SignIn)

		// PRIVATE
		userGroup.DELETE("/sign-out", middleware.AuthMiddleware(db.GormDb), userController.SignOut)
		userGroup.GET("", middleware.AuthMiddleware(db.GormDb), userController.List)
		userGroup.PUT("/:id", middleware.AuthMiddleware(db.GormDb), userController.Update)
		userGroup.GET("/:id", middleware.AuthMiddleware(db.GormDb), userController.Detail)
	}

	// PRIVATE
	roleController := role.NewController(db.GormDb)
	roleGroup := router.Group("/role")
	{
		roleGroup.GET("", middleware.AuthMiddleware(db.GormDb), roleController.Find)
		roleGroup.POST("", middleware.AuthMiddleware(db.GormDb), roleController.Create)
		roleGroup.GET("/:id", middleware.AuthMiddleware(db.GormDb), roleController.Detail)
		roleGroup.PUT("/:id", middleware.AuthMiddleware(db.GormDb), roleController.Update)
		roleGroup.DELETE("/:id", middleware.AuthMiddleware(db.GormDb), roleController.Delete)
	}

	productController := product.NewController(db.GormDb, redis)
	productGroup := router.Group("/product")
	{
		// PUBLIC
		productGroup.GET("", productController.GetProduct)
	}

	brandController := brand.NewController(db.GormDb)
	brandGroup := router.Group("/brand")
	{
		// PUBLIC
		brandGroup.GET("", brandController.GetBrand)
		brandGroup.GET(":slug", brandController.GetBrandBySlug)
	}

	categoryController := category.NewController(db.GormDb)
	categoryGroup := router.Group("/category")
	{
		// PUBLIC
		categoryGroup.GET("", categoryController.GetCategory)
		categoryGroup.GET(":slug")
	}

	return router.Run()
}
