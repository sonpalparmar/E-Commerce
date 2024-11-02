package models

import "gorm.io/gorm"

type UserType string

const (
	Seller UserType = "Seller"
	Buyer  UserType = "Buyer"
)

type BaseUser struct {
	gorm.Model
	Name         string `gorm:"not null"`
	Email        string `gorm:"uniqueIndex; not null"`
	Address      string
	UserType     UserType `gorm:"type:varchar(10);not null"`
	PasswordHash string   `gorm:"not null"`
}

type SellerUser struct {
	BaseUser `gorm:"embedded"`
	PanCard  string
}

type BuyerUser struct {
	BaseUser `gorm:"embedded"`
}

func (BaseUser) TableName() string {
	return "users"
}
func (SellerUser) TableName() string {
	return "users"
}
func (BuyerUser) TableName() string {
	return "users"
}
