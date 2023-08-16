package handler

import (
	"bloock-managed-api/internal/domain"
	"bloock-managed-api/internal/service"
	"bloock-managed-api/internal/service/create/request"
	"github.com/bloock/bloock-sdk-go/v2/entity/key"
	"github.com/gin-gonic/gin"
	"net/http"
)

func PostCreateManagedKey(managedKeyCreateService service.ManagedKeyCreateService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var managedKeyRequest CreateManagedKeyRequest
		err := ctx.BindJSON(&managedKeyRequest)
		if err != nil {
			return
		}

		keyRequest, err := toCreateManagedKeyRequest(managedKeyRequest.Name,
			managedKeyRequest.KeyType,
			managedKeyRequest.Expiration,
			managedKeyRequest.ProtectionLevel)
		if err != nil {
			badRequestAPIError := NewBadRequestAPIError(err.Error())
			ctx.JSON(badRequestAPIError.Status, badRequestAPIError.Message)
			return
		}

		create, err := managedKeyCreateService.Create(
			keyRequest,
		)
		if err != nil {
			serverAPIError := NewInternalServerAPIError(err.Error())
			ctx.JSON(serverAPIError.Status, serverAPIError)
		}

		ctx.JSON(http.StatusCreated, CreateManagedKeyResponse{ID: create.ID})
	}
}

func toCreateManagedKeyRequest(name string, keyType string, expiration int, level int) (request.CreateManagedKeyRequest, error) {
	kty, err := domain.ValidateKeyType(keyType)
	if err != nil {
		return request.CreateManagedKeyRequest{}, err
	}
	switch level {
	case 0:
		return *request.NewCreateManagedKeyRequest(name, kty, expiration, key.KEY_PROTECTION_SOFTWARE), nil
	case 1:
		return *request.NewCreateManagedKeyRequest(name, kty, expiration, key.KEY_PROTECTION_HSM), nil
	}

	return request.CreateManagedKeyRequest{}, NewBadRequestAPIError("invalid request")
}

type CreateManagedKeyRequest struct {
	Name            string `json:"name"`
	KeyType         string `json:"key_type"`
	ProtectionLevel int    `json:"protection_level"`
	Expiration      int    `json:"expiration"`
}

type CreateManagedKeyResponse struct {
	ID string
}
