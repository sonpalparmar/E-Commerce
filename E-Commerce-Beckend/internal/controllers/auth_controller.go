package controllers

import (
	"e-commerce/internal/config"
	"e-commerce/internal/models"
	"e-commerce/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthController struct {
	DB  *gorm.DB
	Cfg config.Config
}

func NewAuthController(db *gorm.DB, cfg config.Config) *AuthController {
	return &AuthController{DB: db, Cfg: cfg}
}

type SignUpInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	UserType string `json:"usertype" binding:"required,oneof=Seller Buyer"`
	Address  string `json:"address"`
	PanCard  string `json:"pancard" binding:"required_if=UserType Seller"`
}

func (ac *AuthController) SignUp(c *gin.Context) {
	var input SignUpInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	baseuser := models.BaseUser{
		Name:         input.Name,
		Email:        input.Email,
		UserType:     models.UserType(input.UserType),
		PasswordHash: hashedPassword,
		Address:      input.Address,
	}
	// start a transaction
	tx := ac.DB.Begin()

	if input.UserType == "Buyer" {
		user := models.BuyerUser{
			BaseUser: baseuser,
		}
		if err := tx.Create(&user).Error; err != nil {
			tx.Rollback()
			utils.RespondWithError(c, http.StatusBadRequest, "Email already exists")
			return
		}
	} else {
		user := models.SellerUser{
			BaseUser: baseuser,
			PanCard:  input.PanCard,
		}
		if err := tx.Create(&user).Error; err != nil {
			tx.Rollback()
			utils.RespondWithError(c, http.StatusBadRequest, "Email already exists")
			return
		}
	}

	tx.Commit()

	utils.RespondWithSuccess(c, http.StatusCreated, gin.H{"message": "user registered successfully"})
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (ac *AuthController) Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	var user models.BaseUser
	if err := ac.DB.Table("users").Where("email = ?", input.Email).First(&user).Error; err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "invalid email or password")
		return
	}

	if !utils.CheckPassword(input.Password, user.PasswordHash) {
		utils.RespondWithError(c, http.StatusBadRequest, "invalid email or password")
		return
	}

	token, err := utils.GenerateToken(user.ID, string(user.UserType), ac.Cfg.JWTSecret)
	if err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "failed to generate token")
		return
	}
	utils.RespondWithSuccess(c, http.StatusOK, gin.H{"token": token})
}
