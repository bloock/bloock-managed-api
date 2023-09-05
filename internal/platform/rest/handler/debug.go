package handler

import (
	"bloock-managed-api/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Debug(processService service.BaseProcessService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var formData postProcessForm
		err := ctx.Bind(&formData)
		if err != nil {
			badRequestAPIError := NewBadRequestAPIError("error binding form")
			ctx.JSON(badRequestAPIError.Status, badRequestAPIError)
			return
		}

		ctx.JSON(http.StatusAccepted, true)
	}
}
