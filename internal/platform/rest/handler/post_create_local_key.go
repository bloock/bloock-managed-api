package handler

import (
	"bloock-managed-api/internal/domain"
	"bloock-managed-api/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func PostCreateLocalKey(localKeyCreateService service.LocalKeyCreateService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		kty := ctx.Query("kty")
		keyType, err := domain.ValidateKeyType(kty)
		if err != nil {
			badRequestAPIError := NewBadRequestAPIError(err.Error())
			ctx.JSON(badRequestAPIError.Status, badRequestAPIError)
			return
		}
		create, err := localKeyCreateService.Create(nil, keyType)
		if err != nil {
			serverAPIError := NewInternalServerAPIError(err.Error())
			ctx.JSON(serverAPIError.Status, serverAPIError)
		}

		ctx.JSON(http.StatusCreated, CreateLocalKeyResponse{ID: create.Id().String()})
	}
}

type CreateLocalKeyResponse struct {
	ID string `json:"id"`
}
