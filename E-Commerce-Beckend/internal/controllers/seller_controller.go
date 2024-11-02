package controllers

import (
	"e-commerce/internal/models"
	"e-commerce/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SellerController struct {
	DB *gorm.DB
}

func NewSellerController(db *gorm.DB) *SellerController {
	return &SellerController{DB: db}
}

type InputProduct struct {
	Title       string  `json:"title" binding:"required"`
	Discription string  `json:"discription"`
	Price       float64 `json:"price" binding:"required"`
	Count       int     `json:"count" binding:"required"`
}

// add products

func (sc *SellerController) AddProducts(c *gin.Context) {
	var input InputProduct
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}
	var availablity bool
	if input.Count > 0 {
		availablity = true
	}

	userid, exists := c.Get("userID")
	if !exists {
		utils.RespondWithError(c, http.StatusNotFound, "user Id not found")
		return
	}

	product := models.Products{
		UserID:      userid.(uint),
		Title:       input.Title,
		Discription: input.Discription,
		Price:       input.Price,
		Count:       input.Count,
		Availablity: availablity,
	}

	tx := sc.DB.Begin()

	if err := tx.Create(&product).Error; err != nil {
		tx.Rollback()
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	tx.Commit()

	utils.RespondWithSuccess(c, http.StatusCreated, gin.H{"message": "successfully products added"})

}

//get all product

func (sc *SellerController) GetAllProducts(c *gin.Context) {
	userid, exists := c.Get("userID")
	if !exists {
		utils.RespondWithError(c, http.StatusNotFound, "user Id not found")
		return
	}

	var allProducts []models.Products

	if err := sc.DB.Table("products").Where("user_id = ?", userid.(uint)).Find(&allProducts).Error; err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to retrieve products")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{"products": allProducts})

}

// update products

func (sc *SellerController) UpdateProducts(c *gin.Context) {
	var product models.Products

	if err := sc.DB.Table("products").Where("id = ?", c.Param("id")).First(&product).Error; err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	userid, exists := c.Get("userID")
	if !exists {
		utils.RespondWithError(c, http.StatusNotFound, "user Id not found")
		return
	}

	if userid.(uint) != product.UserID {
		utils.RespondWithError(c, http.StatusForbidden, "You do not have permission to update this product")
		return
	}

	var input InputProduct

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	product.Title = input.Title
	product.Price = input.Price
	product.Count = input.Count
	if input.Discription != "" {
		product.Discription = input.Discription
	}

	if err := sc.DB.Save(&product).Error; err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, "failed to save updated details")
		return
	}

	utils.RespondWithSuccess(c, http.StatusAccepted, gin.H{"message": "successfuly updated"})
}

func (sc *SellerController) GetSingleProducts(c *gin.Context) {
	var product models.Products

	if err := sc.DB.Table("products").Where("id = ?", c.Param("id")).First(&product).Error; err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	userid, exists := c.Get("userID")
	if !exists {
		utils.RespondWithError(c, http.StatusNotFound, "user Id not found")
		return
	}

	if userid.(uint) != product.UserID {
		utils.RespondWithError(c, http.StatusForbidden, "You do not have permission to get details this product")
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{"product": product})
}

func (sc *SellerController) DeleteProducts(c *gin.Context) {
	var product models.Products

	if err := sc.DB.Table("products").Where("id = ?", c.Param("id")).First(&product).Error; err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	userid, exists := c.Get("userID")
	if !exists {
		utils.RespondWithError(c, http.StatusNotFound, "user Id not found")
		return
	}

	if userid.(uint) != product.UserID {
		utils.RespondWithError(c, http.StatusForbidden, "You do not have permission to delete this product")
		return
	}

	if err := sc.DB.Delete(&product).Error; err != nil {
		utils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithSuccess(c, http.StatusOK, gin.H{"message": "product deleted successfully"})
}
