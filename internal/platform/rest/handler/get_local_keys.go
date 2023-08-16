package handler

import (
	"bloock-managed-api/internal/domain"
	"bloock-managed-api/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetLocalKeys(getLocalKeys service.GetLocalKeysService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		localKeys, err := getLocalKeys.Get(ctx)
		if err != nil {
			serverAPIError := NewInternalServerAPIError(err.Error())
			ctx.JSON(serverAPIError.Status, serverAPIError)
			return
		}

		if len(localKeys) == 0 {
			ctx.JSON(http.StatusNoContent, nil)
			return
		}

		response, err := toGetLocalKeysResponse(localKeys)
		if err != nil {
			badRequestAPIError := NewBadRequestAPIError(err.Error())
			ctx.JSON(badRequestAPIError.Status, badRequestAPIError.Message)
			return
		}

		ctx.JSON(http.StatusOK, response)
	}
}

func toGetLocalKeysResponse(keys []domain.LocalKey) ([]GetLocalKeysResponse, error) {
	var response []GetLocalKeysResponse

	for _, localKey := range keys {

		response = append(response, GetLocalKeysResponse{
			ID:      localKey.Id().String(),
			KeyType: localKey.KeyTypeStr(),
		})
	}

	return response, nil
}

type GetLocalKeysResponse struct {
	ID      string
	KeyType string
}
