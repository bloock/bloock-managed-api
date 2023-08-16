package handler_test

import (
	"bloock-managed-api/internal/domain"
	"bloock-managed-api/internal/platform/rest"
	"bloock-managed-api/internal/platform/rest/handler"
	"bloock-managed-api/internal/service/create"
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

func TestGetLocalKeys(t *testing.T) {
	ctrl := gomock.NewController(t)
	getLocalKeysService := mock_service.NewMockGetLocalKeysService(ctrl)
	server, err := rest.NewServer(
		"localhost",
		"8085",
		getLocalKeysService,
		create.ManagedKey{},
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

	t.Run("given results it should return 200 ok", func(t *testing.T) {
		localKeys := []domain.LocalKey{
			*domain.NewLocalKey(key.LocalKey{
				Key:        "k",
				PrivateKey: "p",
			}, key.EcP256k, uuid.New()),
			*domain.NewLocalKey(key.LocalKey{
				Key:        "k",
				PrivateKey: "p",
			}, key.Rsa4096, uuid.New()),
		}
		expectedResponse := []handler.GetLocalKeysResponse{
			{
				ID:      localKeys[0].Id().String(),
				KeyType: "EcP256k",
			},
			{
				ID:      localKeys[1].Id().String(),
				KeyType: "Rsa4096",
			},
		}
		getLocalKeysService.EXPECT().Get().Return(localKeys, nil)
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/key/local"), nil)

		engine.ServeHTTP(rec, req)

		res := rec.Result()
		assert.Equal(t, http.StatusOK, res.StatusCode)
		var currentResponse []handler.GetLocalKeysResponse
		bytes, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		err = json.Unmarshal(bytes, &currentResponse)
		require.NoError(t, err)
		assert.ElementsMatch(t, expectedResponse, currentResponse)
	})

	t.Run("given no results it should return 204 no content", func(t *testing.T) {

		getLocalKeysService.EXPECT().Get().Return([]domain.LocalKey{}, nil)
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/key/local"), nil)

		engine.ServeHTTP(rec, req)

		res := rec.Result()
		assert.Equal(t, http.StatusNoContent, res.StatusCode)

	})

	t.Run("given no error it should return 500 internal server error", func(t *testing.T) {

		getLocalKeysService.EXPECT().Get().Return([]domain.LocalKey{}, errors.New("some err"))
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/key/local"), nil)

		engine.ServeHTTP(rec, req)

		res := rec.Result()
		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)

	})
}
