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
	merchantController "website-api/controller/merchant"
	"website-api/controller/order"
	permissionController "website-api/controller/permission"
	"website-api/controller/product"
	reviewController "website-api/controller/review"
	"website-api/controller/role"
	shippingController "website-api/controller/shipping"
	uploadController "website-api/controller/upload"
	"website-api/controller/user"
	userAddressController "website-api/controller/user_address"
	"website-api/database"
	"website-api/middleware"
	"website-api/third-party/provider/minio"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	_ "website-api/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Run(db database.DB, redis *redis.Client, minioProvider minio.Provider) (err error) {
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
	// Change password (auth required)
	router.POST("/auth/change-password", middleware.AuthMiddleware(db.GormDb), authController.ChangePassword)

	contentPageController := contentPage.NewController(db.GormDb, redis)
	contentPageGroup := router.Group("/content-page")
	{
		contentPageGroup.GET("/pages/:slug", contentPageController.GetBySlug)
		contentPageGroup.GET("/faqs", contentPageController.GetFaq)
	}

	// Admin FAQ CRUD
	adminCpGroup := router.Group("/admin/content-page", middleware.AuthMiddleware(db.GormDb), middleware.RequireRole("super_admin", "admin"))
	{
		adminCpGroup.GET("/faqs", contentPageController.GetFaqListAdmin)
		adminCpGroup.POST("/faqs", contentPageController.CreateFaq)
		adminCpGroup.PUT("/faqs/:id", contentPageController.UpdateFaq)
		adminCpGroup.DELETE("/faqs/:id", contentPageController.DeleteFaq)
	}

	userController := user.NewController(db.GormDb)
	userGroup := router.Group("/user")
	{
		userGroup.POST("/sign-up", userController.SignUp)
		userGroup.GET("/verify-email", userController.VerifyEmailFromLink)
		userGroup.POST("/verify-email", userController.VerifyEmail)
		userGroup.POST("/sign-in", userController.SignIn)
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

	roleController := role.NewController(db.GormDb)
	roleGroup := router.Group("/role", middleware.AuthMiddleware(db.GormDb), middleware.RequireRole("super_admin", "admin"))
	{
		roleGroup.GET("", roleController.Find)
		roleGroup.POST("", roleController.Create)
		roleGroup.GET("/:id", roleController.Detail)
		roleGroup.PUT("/:id", roleController.Update)
		roleGroup.DELETE("/:id", roleController.Delete)
	}

	// Upload controller (shared for avatar + product images)
	uploadCtl := uploadController.NewController(db.GormDb, minioProvider)

	productController := product.NewController(db.GormDb, redis, minioProvider)
	productGroup := router.Group("/product")
	{
		// PUBLIC
		productGroup.GET("", productController.GetProduct)
		productGroup.GET("/:id", productController.GetProductDetail)
	}

	// Admin product CRUD
	adminProductGroup := router.Group("/admin/product", middleware.AuthMiddleware(db.GormDb), middleware.RequireRole("super_admin", "admin"))
	{
		adminProductGroup.POST("", productController.CreateProduct)
		adminProductGroup.PUT("/:id", productController.UpdateProduct)
		adminProductGroup.DELETE("/:id", productController.DeleteProduct)
	}

	// Admin product image management
	adminProductImageGroup := router.Group("/product", middleware.AuthMiddleware(db.GormDb), middleware.RequireRole("super_admin", "admin"))
	{
		adminProductImageGroup.POST("/:id/image", uploadCtl.UploadProductImage)
		adminProductImageGroup.PUT("/:id/images/:imageId/primary", productController.SetPrimaryImage)
		adminProductImageGroup.DELETE("/:id/images/:imageId", productController.DeleteImage)
		adminProductImageGroup.PUT("/:id/images/reorder", productController.ReorderImages)
	}

	// Review produk
	reviewCtl := reviewController.NewController(db.GormDb)
	router.GET("/product/:id/reviews", reviewCtl.ListByProduct)
	router.POST("/product/:id/reviews", middleware.AuthMiddleware(db.GormDb), reviewCtl.Create)

brandController := brand.NewController(db.GormDb)
	brandGroup := router.Group("/brand")
	{
		// PUBLIC
		brandGroup.GET("", brandController.GetBrand)
		brandGroup.GET("/:slug", brandController.GetBrandBySlug)
	}

	// Admin brand CRUD
	adminBrandGroup := router.Group("/admin/brand", middleware.AuthMiddleware(db.GormDb), middleware.RequireRole("super_admin", "admin"))
	{
		adminBrandGroup.POST("", brandController.Create)
		adminBrandGroup.PUT("/:id", brandController.Update)
		adminBrandGroup.DELETE("/:id", brandController.Delete)
	}

categoryController := category.NewController(db.GormDb)
	categoryGroup := router.Group("/category")
	{
		// PUBLIC
		categoryGroup.GET("", categoryController.GetCategory)
		categoryGroup.GET("/:slug", categoryController.GetCategoryBySlug)
	}

	// Admin category CRUD
	adminCategoryGroup := router.Group("/admin/category", middleware.AuthMiddleware(db.GormDb), middleware.RequireRole("super_admin", "admin"))
	{
		adminCategoryGroup.POST("", categoryController.Create)
		adminCategoryGroup.PUT("/:id", categoryController.Update)
		adminCategoryGroup.DELETE("/:id", categoryController.Delete)
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
		orderGroup.PATCH("/:id/cancel", middleware.AuthMiddleware(db.GormDb), orderController.CancelOrder)
		orderGroup.POST("/notification", orderController.HandleNotification)
	}

	adminOrderController := order.NewController(db.GormDb)
	adminOrderGroup := router.Group("/admin/orders")
	{
		adminOrderGroup.GET("", middleware.AuthMiddleware(db.GormDb), middleware.RequireRole("super_admin", "admin"), adminOrderController.ListAdmin)
		adminOrderGroup.PATCH("/:id/status", middleware.AuthMiddleware(db.GormDb), middleware.RequireRole("super_admin", "admin"), adminOrderController.UpdateStatusAdmin)
	}

	menuCtl := menuController.NewController(db.GormDb)
	menuGroup := router.Group("/menu", middleware.AuthMiddleware(db.GormDb), middleware.RequireRole("super_admin", "admin"))
	{
		menuGroup.GET("", menuCtl.List)
		menuGroup.GET("/tree", menuCtl.Tree)
		menuGroup.POST("", menuCtl.Create)
		menuGroup.PUT("/:id", menuCtl.Update)
		menuGroup.DELETE("/:id", menuCtl.Delete)
	}
	router.GET("/menu/my", middleware.AuthMiddleware(db.GormDb), menuCtl.GetMyMenus)

	roleMenuGroup := router.Group("/role", middleware.AuthMiddleware(db.GormDb), middleware.RequireRole("super_admin", "admin"))
	{
		roleMenuGroup.GET("/:id/menus", menuCtl.GetRoleMenus)
		roleMenuGroup.PUT("/:id/menus", menuCtl.AssignMenus)
	}

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

	merchantCtl := merchantController.NewController(db.GormDb)
	merchantGroup := router.Group("/merchant")
	{
		// Public
		merchantGroup.GET("", merchantCtl.List)
	}

	// Merchant self-service (auth)
	merchantAuthGroup := router.Group("/merchant", middleware.AuthMiddleware(db.GormDb))
	{
		merchantAuthGroup.POST("/register", merchantCtl.Register)
		merchantAuthGroup.GET("/my", merchantCtl.GetMy)
		merchantAuthGroup.PUT("/my", middleware.RequireRole("super_admin", "admin", "merchant"), merchantCtl.UpdateMy)
	}

	// Admin approve merchant
	router.PATCH("/admin/merchant/:id/approve", middleware.AuthMiddleware(db.GormDb), middleware.RequireRole("super_admin", "admin"), merchantCtl.Approve)

	shippingCtl := shippingController.NewController(db.GormDb)
	router.POST("/shipping/cost", middleware.AuthMiddleware(db.GormDb), shippingCtl.Cost)

	// Upload: avatar (self) — definisi awal sudah di atas (uploadCtl)
	router.POST("/upload/image", middleware.AuthMiddleware(db.GormDb), uploadCtl.UploadImage)

	return router.Run()
}
