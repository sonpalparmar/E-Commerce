package routes

import (
	"e-commerce/internal/config"
	"e-commerce/internal/controllers"
	"e-commerce/internal/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB, cfg config.Config) {
	authController := controllers.NewAuthController(db, cfg)
	sellerController := controllers.NewSellerController(db)

	router.POST("/signup", authController.SignUp)
	router.POST("/login", authController.Login)

	protected := router.Group("/")
	protected.Use(middlewares.AuthMiddleware(cfg.JWTSecret))

	seller := protected.Group("/seller")
	seller.Use(middlewares.RoleMiddleware("Seller"))
	{
		seller.POST("/addproducts", sellerController.AddProducts)
		seller.GET("/getallproducts", sellerController.GetAllProducts)
		seller.GET("/updateproducts/:id", sellerController.UpdateProducts)
	}

}
