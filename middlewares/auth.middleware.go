package middlewares

import (
	"dbo-management-app/helpers"
	"dbo-management-app/service"
	"strings"

	"github.com/gin-gonic/gin"
)

func AdminAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")

		// Remove Bearer
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		if tokenString == "" {
			helpers.UnauthorizedResponse(ctx)
			ctx.Abort()
			return
		}

		// Verify Token
		claims, err := service.VerifyToken(tokenString)
		if err != nil {
			helpers.UnauthorizedResponse(ctx)
			ctx.Abort()
			return
		}

		// Check Role
		role := claims.Role
		if role != "admin" {
			helpers.ForbiddenResponse(ctx)
			ctx.Abort()
			return
		}

		// Set information
		ctx.Set("user_id", claims.ID)
		ctx.Set("name", claims.Name)
		ctx.Set("email", claims.Email)
		ctx.Set("role", role)

		ctx.Next()
	}
}

func UserAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")

		// Remove Bearer
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		if tokenString == "" {
			helpers.UnauthorizedResponse(ctx)
			ctx.Abort()
			return
		}

		// Verify Token
		claims, err := service.VerifyToken(tokenString)
		if err != nil {
			helpers.UnauthorizedResponse(ctx)
			ctx.Abort()
			return
		}

		// Set information
		ctx.Set("user_id", claims.ID)
		ctx.Set("name", claims.Name)
		ctx.Set("email", claims.Email)
		ctx.Set("role", claims.Role)

		ctx.Next()
	}
}
