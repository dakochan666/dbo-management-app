package router

import (
	"dbo-management-app/controllers"
	"dbo-management-app/database"
	"dbo-management-app/middlewares"
	"dbo-management-app/seeder"

	"github.com/gin-gonic/gin"
)

func StartApp() *gin.Engine {
	router := gin.Default()

	db := database.ConnectDB()

	// Seeding admin account
	seeder := seeder.NewSeeder(db)
	seeder.SeedAdmin()

	userController := controllers.NewUserController(db)
	productController := controllers.NewProductController(db)
	orderController := controllers.NewOrderController(db)

	userGroup := router.Group("/users")
	{
		userGroup.POST("/register", userController.Register)
		userGroup.POST("/login", userController.Login)
	}

	productGroup := router.Group("/products")
	{
		productGroup.GET("/", productController.GetList)
		productGroup.GET("/:productId", productController.GetDetail)
	}
	// Admin
	productGroup.Use(middlewares.AdminAuth())
	{
		productGroup.POST("/", productController.Create)
		productGroup.PUT("/:productId", productController.Update)
		productGroup.DELETE("/:productId", productController.Delete)
	}

	orderGroup := router.Group("/orders")
	// Customer
	orderGroup.Use(middlewares.UserAuth())
	{
		orderGroup.POST("/", orderController.CheckoutProduct)
		orderGroup.GET("/:orderId", orderController.Detail)
		orderGroup.GET("/customer", orderController.GetListCheckoutProducts)
	}
	// Admin
	orderGroup.Use(middlewares.AdminAuth())
	{
		orderGroup.GET("/admin", orderController.GetListCheckoutAdmin)
	}

	return router
}
