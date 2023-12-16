package controllers

import (
	"dbo-management-app/helpers"
	"dbo-management-app/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderController struct {
	db *gorm.DB
}

func NewOrderController(db *gorm.DB) *OrderController {
	return &OrderController{
		db: db,
	}
}

func (o *OrderController) CheckoutProduct(ctx *gin.Context) {
	var reqOrder models.Order

	// Bind JSON
	if err := ctx.ShouldBindJSON(&reqOrder); err != nil {
		helpers.BadRequestResponse(ctx, err.Error())
		return
	}

	// Validate User and Product IDs
	if reqOrder.ProductID == 0 || reqOrder.Amount <= 0 {
		helpers.BadRequestResponse(ctx, "Product ID, and valid Amount are required")
		return
	}

	// Check if Product exist
	var product models.Product

	productResult := o.db.First(&product, reqOrder.ProductID)

	if productResult.Error != nil {
		helpers.NotFoundResponse(ctx, "Product not found")
		return
	}

	// Check if there is enough stock
	if product.Stock < reqOrder.Amount {
		helpers.BadRequestResponse(ctx, "Insufficient stock for the requested amount")
		return
	}

	userID, ok := ctx.Get("user_id")
	if !ok {
		helpers.InternalServerErrorResponse(ctx, "User information not found in the context")
		return
	}

	newOrder := models.Order{
		UserID:    userID.(uint),
		ProductID: reqOrder.ProductID,
		Amount:    reqOrder.Amount,
	}

	// Begin transaction
	tx := o.db.Begin()

	// Create order
	if err := tx.Create(&newOrder).Error; err != nil {
		tx.Rollback()
		log.Println("Error creating order:", err)
		helpers.InternalServerErrorResponse(ctx, "Error creating order")
		return
	}

	// Update product stock
	if err := tx.Model(&product).Update("stock", gorm.Expr("stock - ?", reqOrder.Amount)).Error; err != nil {
		tx.Rollback()
		helpers.InternalServerErrorResponse(ctx, "Error updating product stock")
		return
	}

	tx.Commit()

	helpers.WriteJsonResponse(ctx, http.StatusCreated, gin.H{
		"success": true,
		"message": "Order success",
		"data":    newOrder,
	})
}

func (o *OrderController) GetListCheckoutProducts(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		helpers.UnauthorizedResponse(ctx)
		return
	}

	// Get Params
	page := helpers.GetQueryInt(ctx, "page", 1)
	limit := helpers.GetQueryInt(ctx, "limit", 10)
	search := ctx.Query("search")
	sort := ctx.Query("sort")

	var orders []models.Order
	query := o.db.Preload("Product").Preload("User").Where("user_id = ?", userID)

	if search != "" {
		query = query.Joins("JOIN products ON orders.product_id = products.id").
			Where("products.name LIKE ?", "%"+search+"%")
	}

	if sort == "asc" {
		query = query.Order("orders.created_at ASC")
	} else {
		query = query.Order("orders.created_at DESC")
	}

	offset := (page - 1) * limit
	query = query.Offset(offset).Limit(limit)

	result := query.Find(&orders)
	if result.Error != nil {
		helpers.InternalServerErrorResponse(ctx, "Error retrieving orders")
		return
	}

	var totalPage int64
	o.db.Model(&models.Order{}).Where("user_id = ?", userID).Count(&totalPage)
	hasNextPage := (page * limit) < int(totalPage)

	response := gin.H{
		"success": true,
		"data":    orders,
		"pagination": gin.H{
			"current_page": page,
			"total_page":   (totalPage + int64(limit) - 1) / int64(limit),
			"limit":        limit,
			"offset":       offset,
			"next_page":    hasNextPage,
		},
	}

	helpers.WriteJsonResponse(ctx, http.StatusOK, response)
}

func (o *OrderController) GetListCheckoutAdmin(ctx *gin.Context) {
	// Get Params
	page := helpers.GetQueryInt(ctx, "page", 1)
	limit := helpers.GetQueryInt(ctx, "limit", 10)
	search := ctx.Query("search")
	sort := ctx.Query("sort")

	var orders []models.Order
	query := o.db.Preload("Product").Preload("User")

	if search != "" {
		query = query.Joins("JOIN products ON orders.product_id = products.id").
			Where("products.name LIKE ?", "%"+search+"%")
	}

	if sort == "asc" {
		query = query.Order("orders.created_at ASC")
	} else {
		query = query.Order("orders.created_at DESC")
	}

	offset := (page - 1) * limit
	query = query.Offset(offset).Limit(limit)

	result := query.Find(&orders)
	if result.Error != nil {
		helpers.InternalServerErrorResponse(ctx, "Error retrieving orders")
		return
	}

	var totalPage int64
	o.db.Model(&models.Order{}).Count(&totalPage)
	hasNextPage := (page * limit) < int(totalPage)

	response := gin.H{
		"success": true,
		"data":    orders,
		"pagination": gin.H{
			"current_page": page,
			"total_page":   (totalPage + int64(limit) - 1) / int64(limit),
			"limit":        limit,
			"offset":       offset,
			"next_page":    hasNextPage,
		},
	}

	helpers.WriteJsonResponse(ctx, http.StatusOK, response)
}

func (o *OrderController) Detail(ctx *gin.Context) {
	orderID := ctx.Param("orderId")
	var order models.Order

	// Check if the order ID is valid
	if err := o.db.Where("id = ?", orderID).Preload("Product").Preload("User").First(&order).Error; err != nil {
		helpers.NotFoundResponse(ctx, "Order not found")
		return
	}

	response := gin.H{
		"success": true,
		"data":    order,
	}

	helpers.WriteJsonResponse(ctx, http.StatusOK, response)
}
