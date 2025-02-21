package routes

import (
	"github.com/Hello256World/shop-api/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRouter(server *gin.Engine, db *gorm.DB) {
	authHandler := NewAuthHandler(db)
	usersHandler := NewUserHandler(db)
	orderHandler := NewOrderHandler(db)
	adminHandler := NewAdminHandler(db)
	productHandler := NewProductHandler(db)
	categoryHandler := NewCategoryHandler(db)
	superAdminHandler := NewSuperAdminHandler(db)
	imageProductHandler := NewImageProductHandler(db)
	specificationHandler := NewSpecificationHandler(db)
	compareProductHandler := NewCompareProductHandler(db)

	mainGroup := server.Group("/v1")

	publicGroup := mainGroup.Group("/public")
	publicGroup.POST("/super-admin-token", superAdminHandler.signin)
	publicGroup.POST("/admin-token", adminHandler.signin)
	publicGroup.POST("/signup", authHandler.signup)
	publicGroup.POST("/otp", authHandler.otp)
	publicGroup.POST("/signin", authHandler.signin)
	publicGroup.GET("/customers", usersHandler.getMe)
	publicGroup.GET("/categories",categoryHandler.getAllActive)
	publicGroup.GET("/products",productHandler.getAllActive)

	/// Super Admin
	superAdminGroup := mainGroup.Group("/limited/")
	superAdminGroup.Use(middleware.SuperAdminAccess)
	superAdminGroup.POST("admins", adminHandler.create)
	superAdminGroup.GET("admins", adminHandler.getAll)
	superAdminGroup.PUT("admins/:id", adminHandler.update)
	superAdminGroup.DELETE("admins/:id", adminHandler.delete)

	adminGroup := mainGroup.Group("/limited/")
	adminGroup.Use(middleware.AdminAccess)
	adminGroup.GET("categories", categoryHandler.getAll)
	adminGroup.POST("categories", categoryHandler.create)

	/// Image Products
	adminGroup.GET("products/:productId/image-product", imageProductHandler.getAll)
	adminGroup.POST("products/:productId/image-product", imageProductHandler.create)
	adminGroup.GET("products/:productId/image-product/:id", imageProductHandler.getById)
	adminGroup.PUT("products/:productId/image-product/:id", imageProductHandler.update)
	adminGroup.DELETE("products/:productId/image-product/:id", imageProductHandler.delete)

	/// Specifications
	adminGroup.GET("products/:productId/specifications", specificationHandler.getAll)
	adminGroup.POST("products/:productId/specifications", specificationHandler.create)
	adminGroup.PUT("products/:productId/specifications/:id", specificationHandler.update)
	adminGroup.DELETE("products/:productId/specifications/:id", specificationHandler.delete)

	/// Orders
	adminGroup.GET("orders", orderHandler.getAll)
	adminGroup.PUT("orders/:id", orderHandler.update)

	/// Compare Products
	adminGroup.GET("products/:productId/compare-products", compareProductHandler.getAll)
	adminGroup.GET("products/:productId/compare-products/:id", compareProductHandler.getById)
	adminGroup.POST("products/:productId/compare-products", compareProductHandler.create)
	adminGroup.PUT("products/:productId/compare-products/:id", compareProductHandler.update)
	adminGroup.DELETE("products/:productId/compare-products/:id", compareProductHandler.delete)

	subAdminGroup := adminGroup.Group("categories/:categoryId")
	subAdminGroup.PUT("", categoryHandler.update)
	subAdminGroup.GET("", categoryHandler.getById)
	subAdminGroup.DELETE("", categoryHandler.delete)

	/// Products
	subAdminGroup.GET("/products", productHandler.getAll)
	subAdminGroup.POST("/products", productHandler.create)
	subAdminGroup.PUT("/products/:productId", productHandler.update)
	subAdminGroup.DELETE("/products/:productId", productHandler.delete)
}
