package helpers

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetQueryInt(ctx *gin.Context, key string, defaultValue int) int {
	valueStr := ctx.Query(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}
