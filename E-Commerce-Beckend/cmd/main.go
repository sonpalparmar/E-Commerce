package main

import (
	"e-commerce/internal/config"
	"e-commerce/internal/database"
	"e-commerce/internal/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	db := database.DatabaseConn(cfg)

	// setup gin router

	router := gin.Default()

	routes.SetupRoutes(router, db, cfg)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

}
