package handler_test

import (
	"bloock-managed-api/internal/domain"
	"bloock-managed-api/internal/platform/rest"
	"bloock-managed-api/internal/platform/test_utils/fixtures"
	mock_service "bloock-managed-api/internal/service/mock"
	"bloock-managed-api/internal/service/process/request"
	"bloock-managed-api/internal/service/process/response"
	"bytes"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProcessServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	processService := mock_service.NewMockBaseProcessService(ctrl)
	availabilityService := mock_service.NewMockAvailabilityService(ctrl)
	notifyService := mock_service.NewMockNotifyService(ctrl)
	updateAnchor := mock_service.NewMockCertificateUpdateAnchorService(ctrl)
	server, err := rest.NewServer(
		"localhost",
		"8085",
		processService,
		availabilityService,
		updateAnchor,
		notifyService,
		"",
		zerolog.Logger{},
		true,
	)
	require.NoError(t, err)
	go server.Start()
	require.NoError(t, err)
	engine := server.Engine()

	t.Run("given service error it should return 500 status code", func(t *testing.T) {
		data := []byte("Hello World")
		integrityEnabled := true
		authenticityEnabled := true
		authenticityKeySource := domain.MANAGED_KEY.String()
		authenticityKeyType := "EcP256k"
		authenticityKid := "768e7955-9690-4ba5-8ff9-23206d14ceb8"
		encryptionEnabled := true
		encryptionKeySource := domain.MANAGED_KEY.String()
		encryptionKeyType := "EcP256k"
		encryptionKid := "768e7955-9690-4ba5-8ff9-23206d14ceb8"
		authenticityUseEnsResolution := true
		availabilityType := domain.HOSTED.String()

		processRequest, err := request.NewProcessRequest(data, integrityEnabled, authenticityEnabled, authenticityKeySource, authenticityKeyType, authenticityKid, authenticityUseEnsResolution, encryptionEnabled, encryptionKeySource, encryptionKeyType, encryptionKid, availabilityType)
		require.NoError(t, err)
		processService.EXPECT().Process(gomock.Any(), *processRequest).Return(nil, errors.New("some error"))
		buf := new(bytes.Buffer)
		writer := multipart.NewWriter(buf)
		part, err := writer.CreateFormFile("file", uuid.New().String())
		require.NoError(t, err)
		_, err = part.Write(data)
		require.NoError(t, err)
		_ = writer.WriteField("integrity.enabled", strconv.FormatBool(integrityEnabled))
		_ = writer.WriteField("authenticity.enabled", strconv.FormatBool(authenticityEnabled))
		_ = writer.WriteField("authenticity.keySource", authenticityKeySource)
		_ = writer.WriteField("authenticity.keyType", authenticityKeyType)
		_ = writer.WriteField("authenticity.key", authenticityKid)
		_ = writer.WriteField("authenticity.useEnsResolution", strconv.FormatBool(authenticityUseEnsResolution))
		_ = writer.WriteField("encryption.enabled", strconv.FormatBool(encryptionEnabled))
		_ = writer.WriteField("encryption.keySource", encryptionKeySource)
		_ = writer.WriteField("encryption.keyType", encryptionKeyType)
		_ = writer.WriteField("encryption.key", encryptionKid)
		_ = writer.WriteField("availability.type", availabilityType)
		_ = writer.Close()
		rec := httptest.NewRecorder()
		require.NoError(t, err)
		url := "/v1/process"
		req := httptest.NewRequest(http.MethodPost, url, buf)
		req.Header.Add("Content-Type", writer.FormDataContentType())
		engine.ServeHTTP(rec, req)

		res := rec.Result()
		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})
}

func TestPostProcessMultipart(t *testing.T) {
	ctrl := gomock.NewController(t)
	processService := mock_service.NewMockBaseProcessService(ctrl)
	availabilityService := mock_service.NewMockAvailabilityService(ctrl)
	notifyService := mock_service.NewMockNotifyService(ctrl)
	updateAnchor := mock_service.NewMockCertificateUpdateAnchorService(ctrl)
	server, err := rest.NewServer(
		"localhost",
		"8085",
		processService,
		availabilityService,
		updateAnchor,
		notifyService,
		"",
		zerolog.Logger{},
		true,
	)
	require.NoError(t, err)
	go server.Start()
	require.NoError(t, err)
	engine := server.Engine()

	tests := []struct {
		name string
		url  string
		file []byte
	}{
		{name: "given a valid request with file it should return 202 when no error occurs", file: fixtures.PDFContent},
		{name: "given a valid request with url it should return 202 when no error occurs", url: "https://valid.url"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			integrityEnabled := true
			authenticityEnabled := true
			authenticityKeySource := domain.MANAGED_KEY.String()
			authenticityKeyType := "EcP256k"
			authenticityKid := "768e7955-9690-4ba5-8ff9-23206d14ceb8"
			authenticityUseEnsResolution := true
			encryptionEnabled := true
			encryptionKeySource := domain.MANAGED_KEY.String()
			encryptionKeyType := "EcP256k"
			encryptionKid := "768e7955-9690-4ba5-8ff9-23206d14ceb8"
			availabilityType := domain.HOSTED.String()

			processRequest, err := request.NewProcessRequest(test.file, integrityEnabled, authenticityEnabled, authenticityKeySource, authenticityKeyType, authenticityKid, authenticityUseEnsResolution, encryptionEnabled, encryptionKeySource, encryptionKeyType, encryptionKid, availabilityType)
			require.NoError(t, err)
			processResponse := response.NewProcessResponseBuilder().Build()
			processService.EXPECT().Process(gomock.Any(), *processRequest).Return(processResponse, nil)
			buf := new(bytes.Buffer)
			writer := multipart.NewWriter(buf)

			if len(test.file) > 0 {
				part, _ := writer.CreateFormFile("file", uuid.New().String())
				_, err = part.Write(test.file)
			}

			if len(test.url) > 0 {
				availabilityService.EXPECT().Download(gomock.Any(), test.url).Return(test.file, nil)
				_ = writer.WriteField("url", test.url)
			}

			require.NoError(t, err)
			_ = writer.WriteField("integrity.enabled", strconv.FormatBool(integrityEnabled))
			_ = writer.WriteField("authenticity.enabled", strconv.FormatBool(authenticityEnabled))
			_ = writer.WriteField("authenticity.keySource", authenticityKeySource)
			_ = writer.WriteField("authenticity.keyType", authenticityKeyType)
			_ = writer.WriteField("authenticity.key", authenticityKid)
			_ = writer.WriteField("authenticity.useEnsResolution", strconv.FormatBool(authenticityUseEnsResolution))
			_ = writer.WriteField("encryption.enabled", strconv.FormatBool(encryptionEnabled))
			_ = writer.WriteField("encryption.keySource", encryptionKeySource)
			_ = writer.WriteField("encryption.keyType", encryptionKeyType)
			_ = writer.WriteField("encryption.key", encryptionKid)
			_ = writer.WriteField("availability.type", availabilityType)

			err = writer.Close()
			require.NoError(t, err)
			rec := httptest.NewRecorder()
			require.NoError(t, err)
			url := "/v1/process"
			req := httptest.NewRequest(http.MethodPost, url, buf)
			req.Header.Add("Content-Type", writer.FormDataContentType())
			engine.ServeHTTP(rec, req)

			res := rec.Result()
			assert.Equal(t, http.StatusAccepted, res.StatusCode)
		})
	}
}

func TestPostProcessBadRequests(t *testing.T) {
	ctrl := gomock.NewController(t)
	processService := mock_service.NewMockBaseProcessService(ctrl)
	availabilityService := mock_service.NewMockAvailabilityService(ctrl)
	notifyService := mock_service.NewMockNotifyService(ctrl)
	updateAnchor := mock_service.NewMockCertificateUpdateAnchorService(ctrl)
	server, err := rest.NewServer(
		"localhost",
		"8085",
		processService,
		availabilityService,
		updateAnchor,
		notifyService,
		"",
		zerolog.Logger{},
		true,
	)
	require.NoError(t, err)
	go server.Start()
	require.NoError(t, err)
	engine := server.Engine()

	kid := "768e7955-9690-4ba5-8ff9-23206d14ceb8"
	ensResolutionEnabled := "true"
	data := []byte("Hello World")
	ecp256k := "EcP256k"
	hosted := domain.HOSTED.String()
	managedKey := domain.HOSTED.String()
	integrityEnabled := "true"
	authenticityEnabled := "true"
	tests := []struct {
		file             []byte
		url              string
		name             string
		integrity        string
		authenticity     string
		keyType          string
		kty              string
		key              string
		availability     string
		useEnsResolution string
	}{
		{name: "given bad url it should return 400 status code", url: "invalid url", integrity: "a", authenticity: authenticityEnabled, keyType: managedKey, kty: ecp256k, key: kid, availability: hosted, useEnsResolution: ensResolutionEnabled},
		{name: "given bad integrity value it should return 400 status code", file: data, integrity: "a", authenticity: authenticityEnabled, keyType: managedKey, kty: ecp256k, key: kid, availability: hosted, useEnsResolution: ensResolutionEnabled},
		{name: "given bad authenticity value it should return 400 status code", file: data, integrity: integrityEnabled, authenticity: "a", keyType: managedKey, kty: ecp256k, key: kid, availability: hosted, useEnsResolution: ensResolutionEnabled},
		{name: "given bad  keyType value it should return 400 status code", file: data, integrity: integrityEnabled, authenticity: authenticityEnabled, keyType: "a", kty: ecp256k, key: kid, availability: hosted, useEnsResolution: ensResolutionEnabled},
		{name: "given bad  kty value it should return 400 status code", file: data, integrity: integrityEnabled, authenticity: authenticityEnabled, keyType: managedKey, kty: "a", key: kid, availability: hosted, useEnsResolution: ensResolutionEnabled},
		{name: "given bad  keyID value it should return 400 status code", file: data, integrity: integrityEnabled, authenticity: authenticityEnabled, keyType: managedKey, kty: ecp256k, key: "a", availability: hosted, useEnsResolution: ensResolutionEnabled},
		{name: "given bad  availability value it should return 400 status code", file: data, integrity: integrityEnabled, authenticity: authenticityEnabled, keyType: managedKey, kty: ecp256k, key: kid, availability: "a", useEnsResolution: ensResolutionEnabled},
		{name: "given bad  ensResolution value it should return 400 status code", file: data, integrity: integrityEnabled, authenticity: authenticityEnabled, keyType: managedKey, kty: ecp256k, key: kid, availability: "a", useEnsResolution: "a"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			writer := multipart.NewWriter(buf)

			part, _ := writer.CreateFormFile("file", uuid.New().String())
			_, err = part.Write(data)
			require.NoError(t, err)
			_ = writer.WriteField("url", test.url)
			_ = writer.WriteField("integrity.enabled", test.integrity)
			_ = writer.WriteField("authenticity.enabled", test.authenticity)
			_ = writer.WriteField("authenticity.keyType", test.keyType)
			_ = writer.WriteField("authenticity.kty", test.kty)
			_ = writer.WriteField("authenticity.key", test.key)
			_ = writer.WriteField("authenticity.useEnsResolution", test.useEnsResolution)
			_ = writer.WriteField("availability.type", test.availability)

			rec := httptest.NewRecorder()
			require.NoError(t, err)
			url := "/v1/process"
			req := httptest.NewRequest(http.MethodPost, url, buf)
			req.Header.Add("Content-Type", writer.FormDataContentType())
			engine.ServeHTTP(rec, req)

			res := rec.Result()
			assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		})
	}
}

func TestWithoutFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	processService := mock_service.NewMockBaseProcessService(ctrl)
	availabilityService := mock_service.NewMockAvailabilityService(ctrl)
	notifyService := mock_service.NewMockNotifyService(ctrl)
	updateAnchor := mock_service.NewMockCertificateUpdateAnchorService(ctrl)
	server, err := rest.NewServer(
		"localhost",
		"8085",
		processService,
		availabilityService,
		updateAnchor,
		notifyService,
		"",
		zerolog.Logger{},
		true,
	)
	require.NoError(t, err)
	go server.Start()
	require.NoError(t, err)
	engine := server.Engine()

	t.Run("given no file or url it should return bad request", func(t *testing.T) {
		kid := "768e7955-9690-4ba5-8ff9-23206d14ceb8"
		useEnsResolution := "true"
		ecp256k := "EcP256k"
		hosted := domain.HOSTED.String()
		managedKey := domain.MANAGED_KEY.String()
		integrityEnabled := "true"
		authenticityEnabled := "true"

		buf := new(bytes.Buffer)
		writer := multipart.NewWriter(buf)

		_ = writer.WriteField("integrity.enabled", integrityEnabled)
		_ = writer.WriteField("authenticity.enabled", authenticityEnabled)
		_ = writer.WriteField("authenticity.keyType", managedKey)
		_ = writer.WriteField("authenticity.kty", ecp256k)
		_ = writer.WriteField("authenticity.key", kid)
		_ = writer.WriteField("authenticity.useEnsResolution", useEnsResolution)
		_ = writer.WriteField("availability.type", hosted)

		assert.NoError(t, err)
		err = writer.Close()
		require.NoError(t, err)
		rec := httptest.NewRecorder()
		require.NoError(t, err)
		url := "/v1/process"
		req := httptest.NewRequest(http.MethodPost, url, buf)
		req.Header.Add("Content-Type", writer.FormDataContentType())
		engine.ServeHTTP(rec, req)

		res := rec.Result()
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	t.Run("given empty file it should return bad request", func(t *testing.T) {
		kid := "768e7955-9690-4ba5-8ff9-23206d14ceb8"
		useEnsResolution := "true"
		ecp256k := "EcP256k"
		hosted := domain.HOSTED.String()
		managedKey := domain.MANAGED_KEY.String()
		integrityEnabled := "true"
		authenticityEnabled := "true"

		buf := new(bytes.Buffer)
		writer := multipart.NewWriter(buf)
		part, err := writer.CreateFormFile("file", uuid.New().String())
		require.NoError(t, err)
		_, err = part.Write([]byte{})
		require.NoError(t, err)
		_ = writer.WriteField("integrity.enabled", integrityEnabled)
		_ = writer.WriteField("authenticity.enabled", authenticityEnabled)
		_ = writer.WriteField("authenticity.keyType", managedKey)
		_ = writer.WriteField("authenticity.kty", ecp256k)
		_ = writer.WriteField("authenticity.key", kid)
		_ = writer.WriteField("authenticity.useEnsResolution", useEnsResolution)
		_ = writer.WriteField("availability.type", hosted)

		assert.NoError(t, err)
		err = writer.Close()
		require.NoError(t, err)
		rec := httptest.NewRecorder()
		require.NoError(t, err)
		url := "/v1/process"
		req := httptest.NewRequest(http.MethodPost, url, buf)
		req.Header.Add("Content-Type", writer.FormDataContentType())
		engine.ServeHTTP(rec, req)

		res := rec.Result()
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})
}
