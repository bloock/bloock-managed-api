package handler

import (
	"bloock-managed-api/internal/service/update"
	"bloock-managed-api/internal/service/update/request"
	"encoding/json"
	"github.com/bloock/bloock-sdk-go/v2/client"
	"github.com/gin-gonic/gin"
	"net/http"
)

func PostReceiveWebhook(certification update.CertificationAnchor, secretKey string, enforceTolerance bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var webhookRequest WebhookRequest
		if err := ctx.BindJSON(&webhookRequest); err != nil {
			NewBadRequestAPIError("invalid json body")
			return
		}
		bloockSignature := ctx.GetHeader("Bloock-Signature")

		webhookClient := client.NewWebhookClient()
		bodyBytes, err := json.Marshal(webhookRequest)
		if err != nil {

			ctx.JSON(http.StatusInternalServerError, NewInternalServerAPIError(err.Error()))
			return
		}
		ok, err := webhookClient.VerifyWebhookSignature(bodyBytes, bloockSignature, secretKey, enforceTolerance)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, NewInternalServerAPIError(err.Error()))
			return
		}
		if !ok {
			ctx.JSON(http.StatusBadRequest, NewBadRequestAPIError("invalid signature"))
			return
		}

		if err := certification.UpdateAnchor(ctx, request.UpdateCertificationAnchorRequest{
			AnchorId: webhookRequest.Data.Network.AnchorId,
			Payload:  webhookRequest,
		}); err != nil {
			ctx.JSON(http.StatusInternalServerError, NewInternalServerAPIError(err.Error()))
			return
		}

		ctx.Status(http.StatusAccepted)
	}
}

type WebhookRequest struct {
	WebhookId string `json:"webhook_id"`
	RequestId string `json:"request_id"`
	Type      string `json:"type"`
	CreatedAt int    `json:"created_at"`
	Data      struct {
		CreatedAt    int  `json:"created_at"`
		Finalized    bool `json:"finalized"`
		Id           int  `json:"id"`
		MessageCount int  `json:"message_count"`
		Network      struct {
			AnchorId  int    `json:"anchor_id"`
			CreatedAt int    `json:"created_at"`
			Name      string `json:"name"`
			Status    string `json:"status"`
			Test      bool   `json:"test"`
			TxHash    string `json:"tx_hash"`
		} `json:"network"`
		Root string `json:"root"`
		Test bool   `json:"test"`
	} `json:"data"`
}
