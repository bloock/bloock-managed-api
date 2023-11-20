package process_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/bloock/bloock-managed-api/internal/config"
	"github.com/bloock/bloock-managed-api/internal/domain"
	"github.com/bloock/bloock-managed-api/internal/platform/rest"
	"github.com/bloock/bloock-managed-api/internal/platform/rest/handler/process/request"
	"github.com/bloock/bloock-managed-api/internal/platform/rest/handler/process/response"
	"github.com/bloock/bloock-managed-api/internal/platform/test_utils/fixtures"
	"github.com/bloock/bloock-sdk-go/v2/client"
	"github.com/bloock/bloock-sdk-go/v2/entity/key"
	"github.com/google/uuid"
	"github.com/mcuadros/go-defaults"
	"github.com/spf13/viper"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	ValidJWT           = ""
	UrlFile            = "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf"
	FileContent        = fixtures.PDFContent
	ManagedKey         = "89fa206f-5213-4d6b-bfb7-c6a421f926b4"
	ManagedCertificate = "46f1ddc9-8997-41cc-8172-a076f3200f4f"
)

func init() {
	apiHost := os.Getenv("API_HOST")
	if apiHost == "" {
		panic("No API Host defined on API_HOST env variable")
	}

	apiKey := os.Getenv("API_KEY")
	if apiHost == "" {
		panic("No API Key defined on API_KEY env variable")
	}

	ValidJWT = os.Getenv("SESSION_TOKEN")
	if apiHost == "" {
		panic("No Session Token defined on SESSION_TOKEN env variable")
	}

	viper.Set("bloock.api_host", apiHost)
	viper.Set("bloock.api_key", apiKey)

	config.InitConfig(zerolog.Logger{})
}

func TestProcessService(t *testing.T) {
	server, err := rest.NewServer(
		zerolog.Logger{},
	)
	require.NoError(t, err)
	go server.Start()
	require.NoError(t, err)

	t.Run("given a valid url, should return a valid hash", func(t *testing.T) {
		req := request.ProcessFormRequest{
			Url: UrlFile,
		}

		res, status, err := sendRequest(server, &req, nil)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, status)
		assert.Equal(t, res.Hash, "c5a2180e2f97506550f1bba5d863bc6ed05ad8b51daf6fca1ac7d396ef3183c5")
	})

	t.Run("given a invalid url, should return an error", func(t *testing.T) {
		req := request.ProcessFormRequest{
			Url: "invalid_url",
		}

		_, status, err := sendRequest(server, &req, nil)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, status)
	})

	t.Run("given no file nor url, should return an error", func(t *testing.T) {
		req := request.ProcessFormRequest{}

		_, status, err := sendRequest(server, &req, nil)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, status)
	})

	t.Run("with integrity enabled, should return a valid response", func(t *testing.T) {
		req := request.ProcessFormRequest{
			Url: UrlFile,
			Integrity: request.ProcessFormIntegrityRequest{
				Enabled: true,
			},
		}

		res, status, err := sendRequest(server, &req, nil)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, status)
		assert.Equal(t, res.Hash, "c5a2180e2f97506550f1bba5d863bc6ed05ad8b51daf6fca1ac7d396ef3183c5")
		assert.NotEmpty(t, res.Integrity.AnchorId)
	})

	t.Run("with integrity enabled and diferent auth methods and environments, should return valid reponse", func(t *testing.T) {
		liveAnchorID := 0

		// with valid JWT Token and no environment set
		req := request.ProcessFormRequest{
			Url: UrlFile,
			Integrity: request.ProcessFormIntegrityRequest{
				Enabled: true,
			},
		}

		res, status, err := sendRequest(server, &req, map[string]string{
			"X-Api-Key": ValidJWT,
		})
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, status)
		assert.Equal(t, res.Hash, "c5a2180e2f97506550f1bba5d863bc6ed05ad8b51daf6fca1ac7d396ef3183c5")
		assert.NotEmpty(t, res.Integrity.AnchorId)
		liveAnchorID = res.Integrity.AnchorId

		// with valid JWT Token and live environment
		req = request.ProcessFormRequest{
			Url: UrlFile,
			Integrity: request.ProcessFormIntegrityRequest{
				Enabled: true,
			},
		}

		res, status, err = sendRequest(server, &req, map[string]string{
			"X-Api-Key":   ValidJWT,
			"Environment": "live",
		})
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, status)
		assert.Equal(t, res.Hash, "c5a2180e2f97506550f1bba5d863bc6ed05ad8b51daf6fca1ac7d396ef3183c5")
		assert.NotEmpty(t, res.Integrity.AnchorId)
		assert.Equal(t, liveAnchorID, res.Integrity.AnchorId)

		// with valid JWT Token and test environment
		req = request.ProcessFormRequest{
			Url: UrlFile,
			Integrity: request.ProcessFormIntegrityRequest{
				Enabled: true,
			},
		}

		res, status, err = sendRequest(server, &req, map[string]string{
			"X-Api-Key":   ValidJWT,
			"Environment": "test",
		})
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, status)
		assert.Equal(t, res.Hash, "c5a2180e2f97506550f1bba5d863bc6ed05ad8b51daf6fca1ac7d396ef3183c5")
		assert.NotEmpty(t, res.Integrity.AnchorId)
		assert.NotEqual(t, liveAnchorID, res.Integrity.AnchorId)
	})

	t.Run("with local key authenticity, should return a valid response", func(t *testing.T) {
		keyClient := client.NewKeyClient()
		authenticityKey, err := keyClient.NewLocalKey(key.Rsa2048)
		if err != nil {
			panic(err)
		}

		viper.Set("authenticity.key.key_type", "Rsa2048")
		viper.Set("authenticity.key.private_key", authenticityKey.PrivateKey)
		viper.Set("authenticity.key.public_key", authenticityKey.Key)

		config.InitConfig(zerolog.Logger{})

		req := request.ProcessFormRequest{
			Url: UrlFile,
			Authenticity: request.ProcessFormAuthenticityRequest{
				Enabled:   true,
				KeySource: domain.LOCAL_KEY.String(),
			},
		}

		res, status, err := sendRequest(server, &req, nil)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, status)
		assert.NotEqual(t, res.Hash, "c5a2180e2f97506550f1bba5d863bc6ed05ad8b51daf6fca1ac7d396ef3183c5")
		assert.Equal(t, authenticityKey.Key, res.Authenticity.Key)
		assert.NotEmpty(t, res.Authenticity.Signature)
	})

	t.Run("with managed key authenticity, should return a valid response", func(t *testing.T) {
		req := request.ProcessFormRequest{
			Url: UrlFile,
			Authenticity: request.ProcessFormAuthenticityRequest{
				Enabled:   true,
				KeySource: domain.MANAGED_KEY.String(),
				Key:       ManagedKey,
			},
		}

		res, status, err := sendRequest(server, &req, nil)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, status)
		assert.NotEqual(t, res.Hash, "c5a2180e2f97506550f1bba5d863bc6ed05ad8b51daf6fca1ac7d396ef3183c5")
		assert.Equal(t, ManagedKey, res.Authenticity.Key)
		assert.NotEmpty(t, res.Authenticity.Signature)
	})

	t.Run("with local certificate authenticity, should return a valid response", func(t *testing.T) {
		keyClient := client.NewKeyClient()
		authenticityKey, err := keyClient.NewLocalCertificate(key.LocalCertificateParams{
			KeyType:  key.Rsa2048,
			Password: "password",
			Subject: key.SubjectCertificateParams{
				CommonName: "a name",
			},
			ExpirationMonths: 0,
		})
		if err != nil {
			panic(err)
		}

		id, err := uuid.NewUUID()
		require.NoError(t, err)
		path := fmt.Sprintf("./tmp/%s.p12", id)
		f, err := os.Create(path)
		require.NoError(t, err)
		defer os.Remove(path)
		_, err = f.Write(authenticityKey.Pkcs12)
		require.NoError(t, err)

		viper.Set("authenticity.certificate.pkcs12_path", path)
		viper.Set("authenticity.certificate.pkcs12_password", "password")

		config.InitConfig(zerolog.Logger{})

		req := request.ProcessFormRequest{
			Url: UrlFile,
			Authenticity: request.ProcessFormAuthenticityRequest{
				Enabled:   true,
				KeySource: domain.LOCAL_CERTIFICATE.String(),
				Key:       ManagedCertificate,
			},
		}

		res, status, err := sendRequest(server, &req, nil)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, status)
		assert.NotEqual(t, res.Hash, "c5a2180e2f97506550f1bba5d863bc6ed05ad8b51daf6fca1ac7d396ef3183c5")
		assert.Equal(t, "certificate_id", res.Authenticity.Key)
		assert.NotEmpty(t, res.Authenticity.Signature)
	})

	t.Run("with managed certificate authenticity, should return a valid response", func(t *testing.T) {
		req := request.ProcessFormRequest{
			Url: UrlFile,
			Authenticity: request.ProcessFormAuthenticityRequest{
				Enabled:   true,
				KeySource: domain.MANAGED_CERTIFICATE.String(),
				Key:       ManagedCertificate,
			},
		}

		res, status, err := sendRequest(server, &req, nil)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, status)
		assert.NotEqual(t, res.Hash, "c5a2180e2f97506550f1bba5d863bc6ed05ad8b51daf6fca1ac7d396ef3183c5")
		assert.Equal(t, ManagedCertificate, res.Authenticity.Key)
		assert.NotEmpty(t, res.Authenticity.Signature)
	})

	t.Run("with local key encryption, should return a valid response", func(t *testing.T) {
		keyClient := client.NewKeyClient()
		encryptionKey, err := keyClient.NewLocalKey(key.Aes128)
		if err != nil {
			panic(err)
		}

		viper.Set("encryption.key.key_type", "Aes128")
		viper.Set("encryption.key.public_key", encryptionKey.Key)

		config.InitConfig(zerolog.Logger{})

		req := request.ProcessFormRequest{
			Url: UrlFile,
			Encryption: request.ProcessFormEncryptionRequest{
				Enabled:   true,
				KeySource: domain.LOCAL_KEY.String(),
			},
		}

		res, status, err := sendRequest(server, &req, nil)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, status)
		assert.Equal(t, res.Hash, "c5a2180e2f97506550f1bba5d863bc6ed05ad8b51daf6fca1ac7d396ef3183c5")
		assert.Equal(t, encryptionKey.Key, res.Encryption.Key)
	})

	// t.Run("with managed key encryption, should return a valid response", func(t *testing.T) {
	// 	req := request.ProcessFormRequest{
	// 		Url: UrlFile,
	// 		Encryption: request.ProcessFormEncryptionRequest{
	// 			Enabled:   true,
	// 			KeySource: domain.MANAGED_KEY.String(),
	// 			Key:       ManagedKey,
	// 		},
	// 	}

	// 	res, status, err := sendRequest(server, &req, nil)
	// 	require.NoError(t, err)

	// 	assert.Equal(t, http.StatusOK, status)
	// 	assert.Equal(t, res.Hash, "c5a2180e2f97506550f1bba5d863bc6ed05ad8b51daf6fca1ac7d396ef3183c5")
	// 	assert.Equal(t, ManagedKey, res.Encryption.Key)
	// })

	// t.Run("with local certificate encryption, should return a valid response", func(t *testing.T) {
	// 	keyClient := client.NewKeyClient()
	// 	authenticityKey, err := keyClient.NewLocalCertificate(key.LocalCertificateParams{
	// 		KeyType:  key.Rsa2048,
	// 		Password: "password",
	// 		Subject: key.SubjectCertificateParams{
	// 			CommonName: "a name",
	// 		},
	// 		ExpirationMonths: 0,
	// 	})
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	id, err := uuid.NewUUID()
	// 	require.NoError(t, err)
	// 	path := fmt.Sprintf("./tmp/%s.p12", id)
	// 	f, err := os.Create(path)
	// 	require.NoError(t, err)
	// 	_, err = f.Write(authenticityKey.Pkcs12)
	// 	require.NoError(t, err)

	// 	viper.Set("encryption.certificate.pkcs12_path", path)
	// 	viper.Set("encryption.certificate.pkcs12_password", "password")

	// 	config.InitConfig(zerolog.Logger{})

	// 	req := request.ProcessFormRequest{
	// 		Url: UrlFile,
	// 		Encryption: request.ProcessFormEncryptionRequest{
	// 			Enabled:   true,
	// 			KeySource: domain.LOCAL_CERTIFICATE.String(),
	// 		},
	// 	}

	// 	res, status, err := sendRequest(server, &req, nil)
	// 	require.NoError(t, err)

	// 	assert.Equal(t, http.StatusOK, status)
	// 	assert.NotEqual(t, res.Hash, "c5a2180e2f97506550f1bba5d863bc6ed05ad8b51daf6fca1ac7d396ef3183c5")
	// 	assert.Equal(t, ManagedKey, res.Encryption.Key)
	// })

	// t.Run("with managed certificate encryption, should return a valid response", func(t *testing.T) {
	// 	req := request.ProcessFormRequest{
	// 		Url: UrlFile,
	// 		Encryption: request.ProcessFormEncryptionRequest{
	// 			Enabled:   true,
	// 			KeySource: domain.MANAGED_CERTIFICATE.String(),
	// 			Key:       ManagedCertificate,
	// 		},
	// 	}

	// 	res, status, err := sendRequest(server, &req, nil)
	// 	require.NoError(t, err)

	// 	assert.Equal(t, http.StatusOK, status)
	// 	assert.Equal(t, res.Hash, "c5a2180e2f97506550f1bba5d863bc6ed05ad8b51daf6fca1ac7d396ef3183c5")
	// 	assert.Equal(t, ManagedCertificate, res.Encryption.Key)
	// })

	t.Run("with hosting availability, should return a valid response", func(t *testing.T) {
		req := request.ProcessFormRequest{
			Url: UrlFile,
			Availability: request.ProcessFormAvailabilityRequest{
				Enabled: true,
				Type:    domain.HOSTED.String(),
			},
		}

		res, status, err := sendRequest(server, &req, nil)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, status)
		assert.Equal(t, res.Hash, "c5a2180e2f97506550f1bba5d863bc6ed05ad8b51daf6fca1ac7d396ef3183c5")
		assert.NotEmpty(t, res.Availability.ID)
		assert.NotEmpty(t, res.Availability.Url)
	})

	t.Run("with ipfs availability, should return a valid response", func(t *testing.T) {
		req := request.ProcessFormRequest{
			Url: UrlFile,
			Availability: request.ProcessFormAvailabilityRequest{
				Enabled: true,
				Type:    domain.IPFS.String(),
			},
		}

		res, status, err := sendRequest(server, &req, nil)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, status)
		assert.Equal(t, res.Hash, "c5a2180e2f97506550f1bba5d863bc6ed05ad8b51daf6fca1ac7d396ef3183c5")
		assert.NotEmpty(t, res.Availability.ID)
		assert.NotEmpty(t, res.Availability.Url)
	})

	t.Run("with local availability and hash strategy, should return a valid response", func(t *testing.T) {
		viper.Set("storage.local_strategy", "HASH")
		viper.Set("storage.local_path", "./tmp")

		config.InitConfig(zerolog.Logger{})

		req := request.ProcessFormRequest{
			Url: UrlFile,
			Availability: request.ProcessFormAvailabilityRequest{
				Enabled: true,
				Type:    domain.LOCAL.String(),
			},
		}

		res, status, err := sendRequest(server, &req, nil)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, status)
		assert.Equal(t, res.Hash, "c5a2180e2f97506550f1bba5d863bc6ed05ad8b51daf6fca1ac7d396ef3183c5")
		assert.Equal(t, res.Availability.ID, "./tmp/c5a2180e2f97506550f1bba5d863bc6ed05ad8b51daf6fca1ac7d396ef3183c5")
		assert.Empty(t, res.Availability.Url)
	})

	t.Run("with local availability and filename strategy, should return a valid response", func(t *testing.T) {
		viper.Set("storage.local_strategy", "FILENAME")
		viper.Set("storage.local_path", "./tmp")

		config.InitConfig(zerolog.Logger{})

		req := request.ProcessFormRequest{
			Url: UrlFile,
			Availability: request.ProcessFormAvailabilityRequest{
				Enabled: true,
				Type:    domain.LOCAL.String(),
			},
		}

		res, status, err := sendRequest(server, &req, nil)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, status)
		assert.Equal(t, res.Hash, "c5a2180e2f97506550f1bba5d863bc6ed05ad8b51daf6fca1ac7d396ef3183c5")
		assert.Equal(t, res.Availability.ID, "./tmp/dummy.pdf")
		assert.Empty(t, res.Availability.Url)
	})
}

func sendRequest(server *rest.Server, r *request.ProcessFormRequest, headers map[string]string) (*response.ProcessResponse, int, error) {
	defaults.SetDefaults(r)
	engine := server.Engine()

	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)

	if r.File != nil {
		fileReader, err := r.File.Open()
		if err != nil {
			return nil, 0, err
		}

		filename := r.File.Filename
		file, err := io.ReadAll(fileReader)
		if err != nil {
			return nil, 0, err
		}

		part, err := writer.CreateFormFile("file", filename)
		if err != nil {
			return nil, 0, err
		}
		_, err = part.Write(file)
		if err != nil {
			return nil, 0, err
		}
	}

	_ = writer.WriteField("url", r.Url)
	_ = writer.WriteField("integrity.enabled", strconv.FormatBool(r.Integrity.Enabled))
	_ = writer.WriteField("authenticity.enabled", strconv.FormatBool(r.Authenticity.Enabled))
	_ = writer.WriteField("authenticity.keySource", r.Authenticity.KeySource)
	_ = writer.WriteField("authenticity.key", r.Authenticity.Key)
	_ = writer.WriteField("encryption.enabled", strconv.FormatBool(r.Encryption.Enabled))
	_ = writer.WriteField("encryption.keySource", r.Encryption.KeySource)
	_ = writer.WriteField("encryption.key", r.Encryption.Key)
	_ = writer.WriteField("availability.enabled", strconv.FormatBool(r.Availability.Enabled))
	_ = writer.WriteField("availability.type", r.Availability.Type)
	_ = writer.Close()

	url := "/v1/process"
	req := httptest.NewRequest(http.MethodPost, url, buf)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	res := rec.Result()
	decoder := json.NewDecoder(res.Body)

	var data response.ProcessResponse
	err := decoder.Decode(&data)
	if err != nil {
		return nil, res.StatusCode, err
	}

	return &data, res.StatusCode, nil
}
