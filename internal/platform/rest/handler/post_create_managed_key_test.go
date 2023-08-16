package handler_test

import (
	"bloock-managed-api/internal/platform/rest"
	"bloock-managed-api/internal/platform/rest/handler"
	"bloock-managed-api/internal/service/create"
	"bloock-managed-api/internal/service/create/request"
	"bloock-managed-api/internal/service/get"
	mock_service "bloock-managed-api/internal/service/mock"
	"bloock-managed-api/internal/service/update"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bloock/bloock-sdk-go/v2/entity/key"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostCreateManagedKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	managedKeyCreateService := mock_service.NewMockManagedKeyCreateService(ctrl)
	server, err := rest.NewServer(
		"localhost",
		"8085",
		get.LocalKeys{},
		managedKeyCreateService,
		create.LocalKey{},
		create.SignatureService{},
		create.Certification{},
		update.CertificationAnchor{},
		"",
		true,
		zerolog.Logger{},
		true,
	)
	require.NoError(t, err)
	go server.Start()
	require.NoError(t, err)
	engine := server.Engine()

	t.Run("given a valid request it should return 201 created", func(t *testing.T) {
		keyType := key.EcP256k
		name := "test"
		expiration := 1
		protectionLevel := key.KEY_PROTECTION_SOFTWARE
		expectedRequest := request.NewCreateManagedKeyRequest(name, keyType, expiration, protectionLevel)
		managedKey := key.ManagedKey{
			ID:         uuid.NewString(),
			Name:       name,
			Protection: protectionLevel,
			KeyType:    keyType,
			Expiration: int64(expiration),
			Key:        "1234",
		}
		managedKeyCreateService.EXPECT().Create(*expectedRequest).Return(managedKey, nil)
		createManagedKeyHttpRequest := handler.CreateManagedKeyRequest{
			Name:            name,
			KeyType:         "EcP256k",
			ProtectionLevel: 0,
			Expiration:      expiration,
		}
		jsonBody, err := json.Marshal(createManagedKeyHttpRequest)
		require.NoError(t, err)
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/v1/key/managed", bytes.NewReader(jsonBody))

		engine.ServeHTTP(rec, req)

		res := rec.Result()
		assert.Equal(t, http.StatusCreated, res.StatusCode)
		var managedKeyResponse handler.CreateManagedKeyResponse
		bytes, err := io.ReadAll(res.Body)
		err = json.Unmarshal(bytes, &managedKeyResponse)
		require.NoError(t, err)
		assert.Equal(t, managedKey.ID, managedKeyResponse.ID)
	})

	t.Run("given a invalid keytype it should return 400 bad request", func(t *testing.T) {
		name := "test"
		expiration := 1
		createManagedKeyHttpRequest := handler.CreateManagedKeyRequest{
			Name:            name,
			KeyType:         "ECP256K",
			ProtectionLevel: 0,
			Expiration:      expiration,
		}
		jsonBody, err := json.Marshal(createManagedKeyHttpRequest)
		require.NoError(t, err)
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/v1/key/managed", bytes.NewReader(jsonBody))

		engine.ServeHTTP(rec, req)

		res := rec.Result()
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	})

	t.Run("given a invalid key protection level it should return 400 bad request", func(t *testing.T) {
		name := "test"
		expiration := 1
		createManagedKeyHttpRequest := handler.CreateManagedKeyRequest{
			Name:            name,
			KeyType:         "EcP256k",
			ProtectionLevel: 2,
			Expiration:      expiration,
		}
		jsonBody, err := json.Marshal(createManagedKeyHttpRequest)
		require.NoError(t, err)
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/v1/key/managed", bytes.NewReader(jsonBody))

		engine.ServeHTTP(rec, req)

		res := rec.Result()
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	})

	t.Run("given error creating the key it should return 500 internal server error", func(t *testing.T) {

		name := "test"
		expiration := 1
		createManagedKeyHttpRequest := handler.CreateManagedKeyRequest{
			Name:            name,
			KeyType:         "EcP256k",
			ProtectionLevel: 0,
			Expiration:      expiration,
		}
		managedKeyCreateService.EXPECT().Create(gomock.Any()).Return(key.ManagedKey{}, fmt.Errorf("some err"))
		jsonBody, err := json.Marshal(createManagedKeyHttpRequest)
		require.NoError(t, err)
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/v1/key/managed", bytes.NewReader(jsonBody))

		engine.ServeHTTP(rec, req)

		res := rec.Result()
		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)

	})
}
