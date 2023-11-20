package webhook

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/bloock/bloock-managed-api/internal/config"
	api_error "github.com/bloock/bloock-managed-api/internal/platform/rest/error"
	"github.com/bloock/bloock-managed-api/internal/service/notify"
	"github.com/rs/zerolog"

	"github.com/bloock/bloock-sdk-go/v2/client"
	"github.com/gin-gonic/gin"
)

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

type WebhookResponse struct {
	Success bool `json:"success"`
}

func PostReceiveWebhook(l zerolog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		buf, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			webhookErr := api_error.NewInternalServerAPIError("error while reading webhook body")
			ctx.JSON(webhookErr.Status, webhookErr)
			return
		}

		var webhookRequest WebhookRequest
		if err = json.Unmarshal(buf, &webhookRequest); err != nil {
			webhookErr := api_error.NewBadRequestAPIError("invalid json body")
			ctx.JSON(webhookErr.Status, webhookErr)
			return
		}
		bloockSignature := ctx.GetHeader("Bloock-Signature")

		webhookClient := client.NewWebhookClient()
		ok, err := webhookClient.VerifyWebhookSignature(buf, bloockSignature, config.Configuration.Bloock.WebhookSecretKey, false)
		if err != nil {
			serverAPIError := api_error.NewInternalServerAPIError(err.Error())
			ctx.JSON(serverAPIError.Status, serverAPIError)
			return
		}
		if !ok {
			badRequestAPIError := api_error.NewBadRequestAPIError("invalid bloock webhook signature")
			ctx.JSON(badRequestAPIError.Status, badRequestAPIError)
			return
		}

		notifyService := notify.NewNotifyService(ctx, l)
		if err = notifyService.Notify(ctx, webhookRequest.Data.Id); err != nil {
			serverAPIError := api_error.NewInternalServerAPIError(err.Error())
			ctx.JSON(serverAPIError.Status, serverAPIError)
			return
		}

		ctx.JSON(http.StatusOK, WebhookResponse{Success: true})
	}
}
