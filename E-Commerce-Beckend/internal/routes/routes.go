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
		seller.PUT("/updateproducts/:id", sellerController.UpdateProducts)
		seller.GET("/getsingleproducts/:id", sellerController.GetSingleProducts)
		seller.DELETE("/deleteproducts/:id", sellerController.DeleteProducts)
	}

	buyer := protected.Group("/buyer")
	buyer.Use(middlewares.RoleMiddleware("Buyer"))
	{

	}

}
