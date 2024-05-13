package test

import (
	"bytes"
	"encoding/json"
	"github.com/bloock/bloock-managed-api/internal/config"
	"github.com/bloock/bloock-managed-api/internal/platform/rest"
	"github.com/bloock/bloock-managed-api/internal/platform/rest/handler/aggregate"
	"github.com/bloock/bloock-managed-api/internal/platform/test_utils/fixtures"
	"github.com/bloock/bloock-managed-api/pkg/request"
	pkg "github.com/bloock/bloock-managed-api/pkg/response"
	"github.com/bloock/bloock-sdk-go/v2"
	"github.com/bloock/bloock-sdk-go/v2/client"
	"github.com/bloock/bloock-sdk-go/v2/entity/key"
	"github.com/mcuadros/go-defaults"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"
)

type File struct {
	Filename string
	Content  []byte
}

var (
	ValidJWT = ""
	UrlFile  = "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf"
	FilePNG  = File{
		Filename: "sample.png",
		Content:  fixtures.PNGContent,
	}
	ManagedKey         = ""
	ManagedCertificate = ""
)

func TestMain(m *testing.M) {
	apiHost := os.Getenv("API_HOST")
	if apiHost == "" {
		panic("No API Host defined on API_HOST env variable")
	}

	apiKey := os.Getenv("API_KEY")
	if apiHost == "" {
		panic("No API Key defined on API_KEY env variable")
	}
	bloock.ApiKey = apiKey
	bloock.ApiHost = apiHost

	ValidJWT = os.Getenv("SESSION_TOKEN")
	if ValidJWT == "" {
		panic("No Session Token defined on SESSION_TOKEN env variable")
	}

	viper.Set("bloock.api_host", apiHost)
	viper.Set("bloock.api_key", apiKey)
	viper.Set("integrity.aggregate_mode", true)

	config.InitConfig()

	keyClient := client.NewKeyClient()
	managedKey, err := keyClient.NewManagedKey(key.ManagedKeyParams{
		Name:    "api_managed_test",
		KeyType: key.Rsa2048,
	})
	if err != nil {
		panic("Couldn't create managed api")
	}
	ManagedKey = managedKey.ID

	managedCertificate, err := keyClient.NewManagedCertificate(key.ManagedCertificateParams{
		KeyType: key.Rsa2048,
		Subject: key.SubjectCertificateParams{
			CommonName: "Test Api Managed BLOOCK",
		},
		ExpirationMonths: 12,
	})
	if err != nil {
		panic("Couldn't create managed certificate")
	}
	ManagedCertificate = managedCertificate.ID
	time.Sleep(30 * time.Second) // waiting azure creating certificate
}

func SendRequest(server *rest.Server, r *request.ProcessFormRequest, file File, headers map[string]string) (*pkg.ProcessResponse, int, error) {
	defaults.SetDefaults(r)
	engine := server.Engine()

	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)

	if file.Filename != "" {
		part, err := writer.CreateFormFile("file", file.Filename)
		if err != nil {
			return nil, 0, err
		}
		_, err = part.Write(file.Content)
		if err != nil {
			return nil, 0, err
		}
	} else {
		_ = writer.WriteField("url", r.Url)
	}
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

	var data pkg.ProcessResponse
	err := decoder.Decode(&data)
	if err != nil {
		return nil, res.StatusCode, err
	}

	return &data, res.StatusCode, nil
}

func PutAggregate(server *rest.Server) (bool, int, error) {
	engine := server.Engine()

	url := "/v1/aggregate"
	req := httptest.NewRequest(http.MethodPut, url, nil)

	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	res := rec.Result()
	decoder := json.NewDecoder(res.Body)

	var response aggregate.PutAggregateResponse
	err := decoder.Decode(&response)
	if err != nil {
		return false, res.StatusCode, err
	}

	return response.Success, res.StatusCode, nil
}
