package handler

import (
	"bloock-managed-api/internal/service"
	"bloock-managed-api/internal/service/create/request"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"net/http"
)

func PostSignData(signatureService service.SignService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		kid := ctx.Query("kid")
		commonName := ctx.Query("cn")
		mr, err := ctx.Request.MultipartReader()
		if err != nil {
			serverAPIError := NewInternalServerAPIError(err.Error())
			ctx.JSON(serverAPIError.Status, serverAPIError)
			return
		}
		keyID, err := uuid.Parse(kid)
		if err != nil {
			badRequestAPIError := NewBadRequestAPIError(err.Error())
			ctx.JSON(badRequestAPIError.Status, badRequestAPIError)
			return
		}
		if mr == nil {
			badRequestAPIError := NewBadRequestAPIError(err.Error())
			ctx.JSON(badRequestAPIError.Status, badRequestAPIError)
			return
		}

		p, err := mr.NextPart()

		if p.FormName() != "payload" {
			badRequestAPIError := NewBadRequestAPIError(err.Error())
			ctx.JSON(badRequestAPIError.Status, badRequestAPIError)
			return
		}

		data, err := io.ReadAll(p)
		if err != nil {
			serverAPIError := NewInternalServerAPIError(err.Error())
			ctx.JSON(serverAPIError.Status, serverAPIError)
		}

		signResponse, err := signatureService.Sign(ctx, *request.NewSignRequest(keyID, &commonName, data))
		if err != nil {
			serverAPIError := NewInternalServerAPIError(err.Error())
			ctx.JSON(serverAPIError.Status, serverAPIError)
		}

		ctx.JSON(http.StatusCreated, signResponse.Record())
	}
}
