package controllers

import (
	"dbo-management-app/helpers"
	"dbo-management-app/models"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductController struct {
	db *gorm.DB
}

func NewProductController(db *gorm.DB) *ProductController {
	return &ProductController{
		db: db,
	}
}

func (p *ProductController) Create(ctx *gin.Context) {
	var reqProduct models.Product

	// Bind JSON
	if err := ctx.ShouldBindJSON(&reqProduct); err != nil {
		helpers.BadRequestResponse(ctx, err.Error())
		return
	}

	result := p.db.Create(&reqProduct)
	if result.Error != nil {
		helpers.BadRequestResponse(ctx, "Wrong payload")
		return
	}

	helpers.WriteJsonResponse(ctx, http.StatusCreated, gin.H{
		"success": true,
		"message": "Product created",
	})
}

func (p *ProductController) GetList(ctx *gin.Context) {
	var products []models.Product

	// Get Params
	page := helpers.GetQueryInt(ctx, "page", 1)
	perPage := helpers.GetQueryInt(ctx, "limit", 5)
	search := ctx.Query("search")

	offset := (page - 1) * perPage

	// Create search condition
	var whereClause string
	var whereArgs []interface{}
	if search != "" {
		whereClause = "name LIKE ?"
		whereArgs = append(whereArgs, "%"+search+"%")
	}

	result := p.db.Where(whereClause, whereArgs...).Offset(offset).Limit(perPage).Find(&products)

	if result.Error != nil {
		helpers.InternalServerErrorResponse(ctx, "Error retrieving products")
		return
	}

	// Get Total Page
	var totalPage int64
	p.db.Model(&models.Product{}).Count(&totalPage)

	hasNextPage := (page * perPage) < int(totalPage)

	response := gin.H{
		"success": true,
		"data":    products,
		"pagination": gin.H{
			"current_page": page,
			"total_page":   (totalPage + int64(perPage) - 1) / int64(perPage),
			"limit":        perPage,
			"offset":       offset,
			"next_page":    hasNextPage,
		},
	}

	helpers.WriteJsonResponse(ctx, http.StatusOK, response)
}

func (p *ProductController) GetDetail(ctx *gin.Context) {
	// Get ID
	productID := ctx.Param("productId")

	// Validate ID
	id, err := strconv.ParseUint(productID, 10, 64)
	if err != nil {
		helpers.BadRequestResponse(ctx, "Invalid product ID")
		return
	}

	var product models.Product
	result := p.db.First(&product, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		helpers.NotFoundResponse(ctx, "Product not found")
		return
	} else if result.Error != nil {
		helpers.InternalServerErrorResponse(ctx, "Error retrieving product details")
		return
	}

	response := gin.H{
		"success": true,
		"message": "Product details retrieved",
		"product": product,
	}

	helpers.WriteJsonResponse(ctx, http.StatusOK, response)
}

func (p *ProductController) Update(ctx *gin.Context) {
	// Get ID
	productID := ctx.Param("productId")

	// Validate ID
	id, err := strconv.ParseUint(productID, 10, 64)
	if err != nil {
		helpers.BadRequestResponse(ctx, "Invalid product ID")
		return
	}

	var existingProduct models.Product
	result := p.db.First(&existingProduct, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		helpers.NotFoundResponse(ctx, "Product not found")
		return
	} else if result.Error != nil {
		helpers.InternalServerErrorResponse(ctx, "Error retrieving product details")
		return
	}

	// Bind JSON
	var updatedProduct models.Product
	if err := ctx.ShouldBindJSON(&updatedProduct); err != nil {
		helpers.BadRequestResponse(ctx, err.Error())
		return
	}

	result = p.db.Model(&existingProduct).Updates(&updatedProduct)

	if result.Error != nil {
		helpers.InternalServerErrorResponse(ctx, "Error updating product details")
		return
	}

	response := gin.H{
		"success": true,
		"message": "Product details updated",
		"product": existingProduct,
	}

	helpers.WriteJsonResponse(ctx, http.StatusOK, response)
}

func (p *ProductController) Delete(ctx *gin.Context) {
	// Get ID
	productID := ctx.Param("productId")

	// Validate ID
	id, err := strconv.ParseUint(productID, 10, 64)
	if err != nil {
		helpers.BadRequestResponse(ctx, "Invalid product ID")
		return
	}

	// Find Product
	var existingProduct models.Product
	result := p.db.First(&existingProduct, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		helpers.NotFoundResponse(ctx, "Product not found")
		return
	} else if result.Error != nil {
		helpers.InternalServerErrorResponse(ctx, "Error retrieving product details")
		return
	}

	// Delete product
	result = p.db.Delete(&existingProduct)

	if result.Error != nil {
		helpers.InternalServerErrorResponse(ctx, "Error deleting product")
		return
	}

	helpers.WriteJsonResponse(ctx, http.StatusOK, gin.H{
		"success": true,
		"message": "Product deleted",
	})
}
