package handler

import (
	"bloock-managed-api/internal/service/update"
	"github.com/bloock/bloock-sdk-go/v2/entity/integrity"
	"github.com/gin-gonic/gin"
	"net/http"
)

func PostReceiveConfirmation(certification update.CertificationAnchor) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var anchor integrity.Anchor
		if err := ctx.BindJSON(&anchor); err != nil {
			ctx.JSON(http.StatusBadRequest, "invalid json body")
			return
		}

		if err := certification.UpdateAnchor(ctx, anchor); err != nil {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.Status(http.StatusAccepted)
	}
}
