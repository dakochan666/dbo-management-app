package controllers

import (
	"dbo-management-app/helpers"
	"dbo-management-app/models"
	"dbo-management-app/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct {
	db *gorm.DB
}

func NewUserController(db *gorm.DB) *UserController {
	return &UserController{
		db: db,
	}
}

func (u *UserController) Register(ctx *gin.Context) {
	var newReqUser models.ReqUser

	// Bind JSON
	err := ctx.ShouldBindJSON(&newReqUser)
	if err != nil {
		helpers.BadRequestResponse(ctx, err.Error())
		return
	}

	// Encrypt password
	encryptedPassword, err := service.HashPassword(newReqUser.Password)
	if err != nil {
		helpers.InternalServerErrorResponse(ctx, err.Error())
		return
	}

	newUser := &models.User{
		Name:     newReqUser.Name,
		Email:    newReqUser.Email,
		Password: string(encryptedPassword),
		Role:     "customer",
	}

	result := u.db.Create(&newUser)
	if result.Error != nil {
		helpers.BadRequestResponse(ctx, "Wrong payload")
		return
	}

	helpers.WriteJsonResponse(ctx, http.StatusCreated, gin.H{
		"success": true,
		"message": "Register success",
	})
}

func (u *UserController) Login(ctx *gin.Context) {
	var reqUser models.ReqUser

	// Bind JSON
	if err := ctx.ShouldBindJSON(&reqUser); err != nil {
		helpers.BadRequestResponse(ctx, err.Error())
		return
	}

	// Check request payload
	if reqUser.Email == "" || reqUser.Password == "" {
		helpers.BadRequestResponse(ctx, "Please fill your email and password")
		return
	}

	// Find User by Email
	var user models.User
	result := u.db.Where("email = ?", reqUser.Email).First(&user)
	if result.Error != nil {
		helpers.BadRequestResponse(ctx, "Wrong email or password")
		return
	}

	// Password Verification
	err := service.ComparePassword([]byte(user.Password), []byte(reqUser.Password))
	if err != nil {
		helpers.BadRequestResponse(ctx, "Wrong email or password")
		return
	}

	// Create Token
	token, err := service.SignToken(user)
	if err != nil {
		helpers.InternalServerErrorResponse(ctx, "Failed to sign token")
		return
	}

	helpers.WriteJsonResponse(ctx, http.StatusOK, gin.H{
		"success": true,
		"token":   token,
	})
}
