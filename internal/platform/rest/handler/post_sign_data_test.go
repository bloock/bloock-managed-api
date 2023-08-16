package handler_test

import (
	"bloock-managed-api/internal/platform/rest"
	"bloock-managed-api/internal/service/create"
	"bloock-managed-api/internal/service/create/request"
	"bloock-managed-api/internal/service/create/response"
	"bloock-managed-api/internal/service/get"
	mock_service "bloock-managed-api/internal/service/mock"
	"bloock-managed-api/internal/service/update"
	"bytes"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostSignData(t *testing.T) {
	ctrl := gomock.NewController(t)
	signApplicationService := mock_service.NewMockSignService(ctrl)
	server, err := rest.NewServer(
		"localhost",
		"8085",
		get.LocalKeys{},
		create.ManagedKey{},
		create.LocalKey{},
		signApplicationService,
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

	t.Run("given a valid request it should return 201 when no error occurs", func(t *testing.T) {
		kid := "768e7955-9690-4ba5-8ff9-23206d14ceb8"
		cn := "test"
		data := []byte("Hello World")
		keyID, err := uuid.Parse(kid)
		require.NoError(t, err)
		signApplicationService.EXPECT().Sign(*request.NewSignRequest(keyID, &cn, data)).Return(response.NewSignResponse(data), nil)
		buf := new(bytes.Buffer)
		w := multipart.NewWriter(buf)
		part, err := w.CreateFormFile("payload", uuid.New().String())
		assert.NoError(t, err)
		part.Write(data)
		w.Close()
		rec := httptest.NewRecorder()
		require.NoError(t, err)
		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/v1/sign?kid=%s&cn=%s", kid, cn), buf)
		req.Header.Add("Content-Type", w.FormDataContentType())
		engine.ServeHTTP(rec, req)

		res := rec.Result()
		assert.Equal(t, http.StatusCreated, res.StatusCode)
	})
}
