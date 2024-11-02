package models

import "gorm.io/gorm"

type Products struct {
	gorm.Model
	UserID         uint   `gorm:"not null"`
	Title          string `gorm:"not null"`
	Discription    string
	Price          float64 `gorm:"not null"`
	CurrencyID     string  `gorm:"default:INR"`
	CurrencyFormat string  `gorm:"default:₹"`
	ProductImgs    []byte
	Count          int  `gorm:"default:0"`
	Availablity    bool `gorm:"default:false"`
}
