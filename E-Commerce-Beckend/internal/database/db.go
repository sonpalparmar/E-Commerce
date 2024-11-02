package database

import (
	"e-commerce/internal/config"
	"e-commerce/internal/models"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DatabaseConn(cfg config.Config) *gorm.DB {
	dsn := "host=" + cfg.DBHost +
		" user=" + cfg.DBUser +
		" password=" + cfg.DBPassword +
		" dbname=" + cfg.DBName +
		" port=" + cfg.DBPort +
		" sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect Database")
	}
	fmt.Println("connection established")

	err = db.AutoMigrate(&models.SellerUser{}, &models.Products{})
	if err != nil {
		log.Fatal("Failed to migrate database")
	}

	fmt.Println("successfully migrate")

	return db
}
