package router

import (
	"website-api/controller/auth"
	"website-api/controller/brand"
	"website-api/controller/cart"
	"website-api/controller/category"
	contentPage "website-api/controller/content-page"
	healthCheck "website-api/controller/health-check"
	"website-api/controller/master"
	menuController "website-api/controller/menu"
	"website-api/controller/order"
	permissionController "website-api/controller/permission"
	"website-api/controller/product"
	"website-api/controller/role"
	"website-api/controller/user"
	userAddressController "website-api/controller/user_address"
	"website-api/database"
	"website-api/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	_ "website-api/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Run(db database.DB, redis *redis.Client) (err error) {
	router := gin.Default()
	router.Use(middleware.NgrokSkipWarning())
	router.Use(cors.New(corsConfig))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// PUBLIC API
	healthController := healthCheck.NewController(db.SqlDb)
	router.GET("/health", healthController.Check)

	authController := auth.NewController(db.GormDb)
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/forgot-password", authController.ForgotPassword)
		authGroup.POST("/reset-password", authController.ResetPassword)
		authGroup.POST("/resend-verification", authController.ResendVerification)
	}

	contentPageController := contentPage.NewController(db.GormDb, redis)
	contentPageGroup := router.Group("/content-page")
	{
		// Public
		contentPageGroup.GET("/pages/:slug", contentPageController.GetBySlug)
		contentPageGroup.GET("/faqs", contentPageController.GetFaq)
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
		userGroup.GET("", middleware.AuthMiddleware(db.GormDb), middleware.RequireRole("super_admin", "admin"), userController.List)
		userGroup.PUT("/:id", middleware.AuthMiddleware(db.GormDb), userController.Update)
		userGroup.GET("/:id", middleware.AuthMiddleware(db.GormDb), userController.Detail)
	}

	userAddressController := userAddressController.NewController(db.GormDb, redis)
	userAddressGroup := router.Group("/user-address")
	{
		userAddressGroup.GET("", middleware.AuthMiddleware(db.GormDb), userAddressController.List)
		userAddressGroup.POST("", middleware.AuthMiddleware(db.GormDb), userAddressController.Create)
		userAddressGroup.GET("/:id", middleware.AuthMiddleware(db.GormDb), userAddressController.Detail)
		userAddressGroup.DELETE("/:id", middleware.AuthMiddleware(db.GormDb), userAddressController.Delete)
		userAddressGroup.PUT("/:id", middleware.AuthMiddleware(db.GormDb), userAddressController.Update)
		userAddressGroup.PATCH("/:id/primary", middleware.AuthMiddleware(db.GormDb), userAddressController.ChangePrimary)
	}

	// PRIVATE
	roleController := role.NewController(db.GormDb)
	roleGroup := router.Group("/role", middleware.AuthMiddleware(db.GormDb), middleware.RequireRole("super_admin", "admin"))
	{
		roleGroup.GET("", roleController.Find)
		roleGroup.POST("", roleController.Create)
		roleGroup.GET("/:id", roleController.Detail)
		roleGroup.PUT("/:id", roleController.Update)
		roleGroup.DELETE("/:id", roleController.Delete)
	}

	productController := product.NewController(db.GormDb, redis)
	productGroup := router.Group("/product")
	{
		// PUBLIC
		productGroup.GET("", productController.GetProduct)
		productGroup.GET("/:id", productController.GetProductDetail)
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
	}

	cartController := cart.NewController(db.GormDb)
	cartGroup := router.Group("/cart")
	{
		cartGroup.GET("", middleware.AuthMiddleware(db.GormDb), cartController.GetCart)
		cartGroup.POST("", middleware.AuthMiddleware(db.GormDb), cartController.Create)
		cartGroup.DELETE("", middleware.AuthMiddleware(db.GormDb), cartController.Delete)
	}

	orderController := order.NewController(db.GormDb)
	orderGroup := router.Group("/order")
	{
		orderGroup.POST("", middleware.AuthMiddleware(db.GormDb), orderController.Checkout)
		orderGroup.GET("", middleware.AuthMiddleware(db.GormDb), orderController.List)
		orderGroup.GET("/:id", middleware.AuthMiddleware(db.GormDb), orderController.Detail)
		orderGroup.POST("/:id/payment-link", middleware.AuthMiddleware(db.GormDb), orderController.CreatePaymentLink)
		// PUBLIC - webhook dari Midtrans (tanpa JWT)
		orderGroup.POST("/notification", orderController.HandleNotification)
	}

	// Admin orders (perlu role admin/super_admin - di-guard di Fase 2 via Permission)
	adminOrderController := order.NewController(db.GormDb)
	adminOrderGroup := router.Group("/admin/orders")
	{
		adminOrderGroup.GET("", middleware.AuthMiddleware(db.GormDb), middleware.RequireRole("super_admin", "admin"), adminOrderController.ListAdmin)
	}

	// MENU (RBAC) - hanya super_admin/admin
	menuCtl := menuController.NewController(db.GormDb)
	menuGroup := router.Group("/menu", middleware.AuthMiddleware(db.GormDb), middleware.RequireRole("super_admin", "admin"))
	{
		menuGroup.GET("", menuCtl.List)
		menuGroup.GET("/tree", menuCtl.Tree)
		menuGroup.POST("", menuCtl.Create)
		menuGroup.PUT("/:id", menuCtl.Update)
		menuGroup.DELETE("/:id", menuCtl.Delete)
	}
	// Menu untuk user (navigasi, auth saja)
	router.GET("/menu/my", middleware.AuthMiddleware(db.GormDb), menuCtl.GetMyMenus)

	// Role menu assignment - hanya super_admin/admin
	roleMenuGroup := router.Group("/role", middleware.AuthMiddleware(db.GormDb), middleware.RequireRole("super_admin", "admin"))
	{
		roleMenuGroup.GET("/:id/menus", menuCtl.GetRoleMenus)
		roleMenuGroup.PUT("/:id/menus", menuCtl.AssignMenus)
	}

	// Permissions list - hanya super_admin/admin
	permCtl := permissionController.NewController(db.GormDb)
	router.GET("/permission", middleware.AuthMiddleware(db.GormDb), middleware.RequireRole("super_admin", "admin"), permCtl.List)

	masterController := master.NewController(db.GormDb)
	masterGroup := router.Group("/master")
	{
		masterGroup.GET("/provincies", masterController.GetProvince)
		masterGroup.GET("/cities", masterController.GetCities)
		masterGroup.GET("/district", masterController.GetDistrict)
		masterGroup.GET("/subdistricts", masterController.GetSubdistrict)
	}

	return router.Run()
}
