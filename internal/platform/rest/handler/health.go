package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Health() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"success": true,
		})
	}
}
