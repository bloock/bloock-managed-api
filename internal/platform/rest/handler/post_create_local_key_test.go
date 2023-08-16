package handler_test

import (
	"bloock-managed-api/internal/domain"
	"bloock-managed-api/internal/platform/rest"
	"bloock-managed-api/internal/platform/rest/handler"
	"bloock-managed-api/internal/service/create"
	"bloock-managed-api/internal/service/get"
	mock_service "bloock-managed-api/internal/service/mock"
	"bloock-managed-api/internal/service/update"
	"encoding/json"
	"errors"
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

func TestPostCreateLocalKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	localKeyCreateService := mock_service.NewMockLocalKeyCreateService(ctrl)
	server, err := rest.NewServer(
		"localhost",
		"8085",
		get.LocalKeys{},
		create.ManagedKey{},
		localKeyCreateService,
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

	endpoint := "/v1/key/local?kty=%s"
	t.Run("given a valid request it should return 201 when no error occurs", func(t *testing.T) {
		keyType := key.EcP256k
		localKey := domain.NewLocalKey(key.LocalKey{}, keyType, uuid.New())
		localKeyCreateService.EXPECT().Create(keyType).Return(*localKey, nil)
		rec := httptest.NewRecorder()
		kty := "EcP256k"
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf(endpoint, kty), nil)
		engine.ServeHTTP(rec, req)

		res := rec.Result()
		assert.Equal(t, http.StatusCreated, res.StatusCode)
		var createLocalKeyResponse handler.CreateLocalKeyResponse
		bytes, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		err = json.Unmarshal(bytes, &createLocalKeyResponse)
		require.NoError(t, err)
	})

	t.Run("given a invalid keytype it should return 400", func(t *testing.T) {
		rec := httptest.NewRecorder()
		kty := "ecp256k"
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf(endpoint, kty), nil)
		engine.ServeHTTP(rec, req)

		res := rec.Result()
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	t.Run("given a valid request it should return 500 when error occurs", func(t *testing.T) {
		localKeyCreateService.EXPECT().Create(key.EcP256k).Return(domain.LocalKey{}, errors.New("some err"))
		rec := httptest.NewRecorder()
		kty := "EcP256k"
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf(endpoint, kty), nil)
		engine.ServeHTTP(rec, req)

		res := rec.Result()
		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})
}
