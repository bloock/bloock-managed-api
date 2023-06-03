package handler

import (
	"bloock-managed-api/internal/service/update"
	"encoding/json"
	"github.com/bloock/bloock-sdk-go/v2/client"
	"github.com/bloock/bloock-sdk-go/v2/entity/integrity"
	"github.com/gin-gonic/gin"
	"net/http"
)

func PostReceiveConfirmation(certification update.CertificationAnchor, secretKey string, enforceTolerance bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var webhookRequest integrity.Anchor
		if err := ctx.BindJSON(&webhookRequest); err != nil {
			ctx.JSON(http.StatusBadRequest, "invalid json body")
			return
		}

		bloockSignature := ctx.GetHeader("Bloock-Signature")

		webhookClient := client.NewWebhookClient()
		bodyBytes, err := json.Marshal(webhookRequest)
		if err != nil {

			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		ok, err := webhookClient.VerifyWebhookSignature(bodyBytes, bloockSignature, secretKey, enforceTolerance)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		if !ok {
			ctx.JSON(http.StatusBadRequest, "invalid signature")
			return
		}

		if err := certification.UpdateAnchor(ctx, webhookRequest); err != nil {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.Status(http.StatusAccepted)
	}
}
