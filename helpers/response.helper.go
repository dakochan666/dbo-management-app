package helpers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func BadRequestResponse(ctx *gin.Context, payload interface{}) {
	WriteJsonResponse(ctx, http.StatusBadRequest, gin.H{
		"success": false,
		"error":   payload,
	})
}

func UnauthorizedResponse(ctx *gin.Context) {
	WriteJsonResponse(ctx, http.StatusUnauthorized, gin.H{
		"success": false,
		"error":   "Invalid token",
	})
}

func ForbiddenResponse(ctx *gin.Context) {
	WriteJsonResponse(ctx, http.StatusForbidden, gin.H{
		"success": false,
		"error":   "Permission denied",
	})
}

func InternalServerErrorResponse(ctx *gin.Context, payload interface{}) {
	WriteJsonResponse(ctx, http.StatusInternalServerError, gin.H{
		"success": false,
		"error":   payload,
	})
}

func NotFoundResponse(ctx *gin.Context, payload interface{}) {
	WriteJsonResponse(ctx, http.StatusNotFound, gin.H{
		"success": false,
		"error":   payload,
	})
}

func WriteJsonResponse(ctx *gin.Context, status int, payload interface{}) {
	ctx.JSON(status, payload)
}
