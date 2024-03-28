package test

import (
	"fmt"
	"github.com/bloock/bloock-managed-api/internal/platform/repository/sql/connection"
	"github.com/bloock/bloock-sdk-go/v2/client"
	"github.com/bloock/bloock-sdk-go/v2/entity/key"
	"github.com/google/uuid"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/bloock/bloock-managed-api/internal/config"
	"github.com/bloock/bloock-managed-api/internal/domain"
	"github.com/bloock/bloock-managed-api/internal/platform/rest"
	"github.com/bloock/bloock-managed-api/pkg/request"
	"github.com/spf13/viper"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProcessService(t *testing.T) {
	logger := zerolog.Logger{}
	entConnector := connection.NewEntConnector(logger)
	conn, err := connection.NewEntConnection(config.Configuration.Db.ConnectionString, entConnector, logger)
	if err != nil {
		panic(err)
	}
	err = conn.Migrate()
	if err != nil {
		panic(err)
	}

	server, err := rest.NewServer(
		zerolog.Logger{},
		conn,
	)
	require.NoError(t, err)
	go func() {
		err = server.Start()
		require.NoError(t, err)
	}()

	t.Run("given a valid url, should return a valid hash", func(t *testing.T) {
		req := request.ProcessFormRequest{
			Url: UrlFile,
		}

		res, status, err := SendRequest(server, &req, File{}, nil)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, status)
		assert.Equal(t, res.Hash, "c5a2180e2f97506550f1bba5d863bc6ed05ad8b51daf6fca1ac7d396ef3183c5")
	})

	t.Run("given a invalid url, should return an error", func(t *testing.T) {
		req := request.ProcessFormRequest{
			Url: "invalid_url",
		}

		_, status, err := SendRequest(server, &req, File{}, nil)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, status)
	})

	t.Run("given no file nor url, should return an error", func(t *testing.T) {
		req := request.ProcessFormRequest{}

		_, status, err := SendRequest(server, &req, File{}, nil)
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

		res, status, err := SendRequest(server, &req, File{}, nil)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, status)
		assert.Equal(t, res.Hash, "c5a2180e2f97506550f1bba5d863bc6ed05ad8b51daf6fca1ac7d396ef3183c5")
		assert.True(t, res.Integrity.Enabled)
		assert.NotEmpty(t, res.Integrity.AnchorId)
	})

	t.Run("with integrity enabled and api key set, should return valid response", func(t *testing.T) {
		// with valid JWT Token and no environment set
		req := request.ProcessFormRequest{
			Url: UrlFile,
			Integrity: request.ProcessFormIntegrityRequest{
				Enabled: true,
			},
		}

		res, status, err := SendRequest(server, &req, File{}, map[string]string{
			"X-Api-Key": ValidJWT,
		})
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, status)
		assert.Equal(t, res.Hash, "c5a2180e2f97506550f1bba5d863bc6ed05ad8b51daf6fca1ac7d396ef3183c5")
		assert.True(t, res.Integrity.Enabled)
		assert.NotEmpty(t, res.Integrity.AnchorId)
	})

	t.Run("with local key authenticity, should return a valid response", func(t *testing.T) {
		keyClient := client.NewKeyClient()
		authenticityKey, err := keyClient.NewLocalKey(key.Rsa2048)
		require.NoError(t, err)

		viper.Set("authenticity.key.key_type", "Rsa2048")
		viper.Set("authenticity.key.key", authenticityKey.PrivateKey)

		config.InitConfig(zerolog.Logger{})

		req := request.ProcessFormRequest{
			Url: UrlFile,
			Authenticity: request.ProcessFormAuthenticityRequest{
				Enabled:   true,
				KeySource: domain.LOCAL_KEY.String(),
			},
		}

		log.Println(config.Configuration.Authenticity.KeyConfig.Key)
		res, status, err := SendRequest(server, &req, FilePNG, nil)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, status)
		assert.True(t, res.Authenticity.Enabled)
		assert.Equal(t, authenticityKey.Key, res.Authenticity.Signatures[0].Kid)
		assert.NotEmpty(t, res.Authenticity.Signatures)
	})

	t.Run("with managed key authenticity but signing a pdf, should return a error", func(t *testing.T) {
		req := request.ProcessFormRequest{
			Url: UrlFile,
			Authenticity: request.ProcessFormAuthenticityRequest{
				Enabled:   true,
				KeySource: domain.MANAGED_KEY.String(),
				Key:       ManagedKey,
			},
		}

		_, status, err := SendRequest(server, &req, File{}, nil)
		require.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, status)
	})

	t.Run("with managed key authenticity but signing any other file like PNG, should return a valid response", func(t *testing.T) {
		req := request.ProcessFormRequest{
			Authenticity: request.ProcessFormAuthenticityRequest{
				Enabled:   true,
				KeySource: domain.MANAGED_KEY.String(),
				Key:       ManagedKey,
			},
		}

		res, status, err := SendRequest(server, &req, FilePNG, nil)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, status)
		assert.True(t, res.Authenticity.Enabled)
		assert.NotEmpty(t, res.Authenticity.Signatures)
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
		require.NoError(t, err)

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

		res, status, err := SendRequest(server, &req, File{}, nil)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, status)
		assert.True(t, res.Authenticity.Enabled)
		assert.NotEmpty(t, res.Authenticity.Signatures)
	})

	t.Run("with managed certificate authenticity when signing a PDF or any other file, should return a valid response", func(t *testing.T) {
		req := request.ProcessFormRequest{
			Url: UrlFile,
			Authenticity: request.ProcessFormAuthenticityRequest{
				Enabled:   true,
				KeySource: domain.MANAGED_CERTIFICATE.String(),
				Key:       ManagedCertificate,
			},
		}

		res, status, err := SendRequest(server, &req, File{}, nil)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, status)
		assert.True(t, res.Authenticity.Enabled)
		assert.NotEmpty(t, res.Authenticity.Signatures)

		res, status, err = SendRequest(server, &req, FilePNG, nil)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, status)
		assert.True(t, res.Authenticity.Enabled)
		assert.NotEmpty(t, res.Authenticity.Signatures)
	})

	t.Run("with local key encryption, should return a valid response", func(t *testing.T) {
		keyClient := client.NewKeyClient()
		encryptionKey, err := keyClient.NewLocalKey(key.Aes128)
		require.NoError(t, err)

		viper.Set("encryption.key.key_type", "Aes128")
		viper.Set("encryption.key.key", encryptionKey.Key)

		config.InitConfig(zerolog.Logger{})

		req := request.ProcessFormRequest{
			Url: UrlFile,
			Encryption: request.ProcessFormEncryptionRequest{
				Enabled:   true,
				KeySource: domain.LOCAL_KEY.String(),
			},
		}

		res, status, err := SendRequest(server, &req, File{}, nil)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, status)
		assert.True(t, res.Encryption.Enabled)
		assert.Equal(t, "A256GCM", res.Encryption.Alg)
	})

	t.Run("with managed key encryption, should return a valid response", func(t *testing.T) {
		req := request.ProcessFormRequest{
			Url: UrlFile,
			Encryption: request.ProcessFormEncryptionRequest{
				Enabled:   true,
				KeySource: domain.MANAGED_KEY.String(),
				Key:       ManagedKey,
			},
		}

		res, status, err := SendRequest(server, &req, File{}, nil)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, status)
		assert.True(t, res.Encryption.Enabled)
		assert.Equal(t, "RSA_M", res.Encryption.Alg)
	})

	/*t.Run("with local certificate encryption, should return a valid response", func(t *testing.T) {
	 	keyClient := client.NewKeyClient()
	 	authenticityKey, err := keyClient.NewLocalCertificate(key.LocalCertificateParams{
	 		KeyType:  key.Rsa2048,
	 		Password: "password",
			Subject: key.SubjectCertificateParams{
				CommonName: "a name",
			},
	 		ExpirationMonths: 0,
	 	})
	 	require.NoError(t, err)

	 	id, err := uuid.NewUUID()
	 	require.NoError(t, err)
	 	path := fmt.Sprintf("./tmp/%s.p12", id)
	 	f, err := os.Create(path)
	 	require.NoError(t, err)
	 	_, err = f.Write(authenticityKey.Pkcs12)
	 	require.NoError(t, err)

		viper.Set("encryption.certificate.pkcs12_path", path)
	 	viper.Set("encryption.certificate.pkcs12_password", "password")

	 	config.InitConfig(zerolog.Logger{})

	 	req := request.ProcessFormRequest{
	 		Url: UrlFile,
	 		Encryption: request.ProcessFormEncryptionRequest{
	 			Enabled:   true,
	 			KeySource: domain.LOCAL_CERTIFICATE.String(),
	 		},
	 	}

	 	res, status, err := SendRequest(server, &req, false, nil)
	 	require.NoError(t, err)

	 	assert.Equal(t, http.StatusOK, status)
	 	assert.Equal(t, "Rsa2048", res.Encryption.Alg)
	})*/

	t.Run("with managed certificate encryption, should return a valid response", func(t *testing.T) {
		req := request.ProcessFormRequest{
			Url: UrlFile,
			Encryption: request.ProcessFormEncryptionRequest{
				Enabled:   true,
				KeySource: domain.MANAGED_CERTIFICATE.String(),
				Key:       ManagedCertificate,
			},
		}

		res, status, err := SendRequest(server, &req, File{}, nil)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, status)
		assert.Equal(t, "RSA_M", res.Encryption.Alg)
	})

	t.Run("with hosting availability, should return a valid response", func(t *testing.T) {
		req := request.ProcessFormRequest{
			Url: UrlFile,
			Availability: request.ProcessFormAvailabilityRequest{
				Enabled: true,
				Type:    domain.HOSTED.String(),
			},
		}

		res, status, err := SendRequest(server, &req, File{}, nil)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, status)
		assert.Equal(t, res.Hash, "c5a2180e2f97506550f1bba5d863bc6ed05ad8b51daf6fca1ac7d396ef3183c5")
		assert.True(t, res.Availability.Enabled)
		assert.Equal(t, domain.HOSTED.String(), res.Availability.Type)
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

		res, status, err := SendRequest(server, &req, File{}, nil)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, status)
		assert.Equal(t, res.Hash, "c5a2180e2f97506550f1bba5d863bc6ed05ad8b51daf6fca1ac7d396ef3183c5")
		assert.True(t, res.Availability.Enabled)
		assert.Equal(t, domain.IPFS.String(), res.Availability.Type)
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

		res, status, err := SendRequest(server, &req, File{}, nil)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, status)
		assert.Equal(t, res.Hash, "c5a2180e2f97506550f1bba5d863bc6ed05ad8b51daf6fca1ac7d396ef3183c5")
		assert.True(t, res.Availability.Enabled)
		assert.Equal(t, domain.LOCAL.String(), res.Availability.Type)
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

		res, status, err := SendRequest(server, &req, File{}, nil)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, status)
		assert.Equal(t, res.Hash, "c5a2180e2f97506550f1bba5d863bc6ed05ad8b51daf6fca1ac7d396ef3183c5")
		assert.True(t, res.Availability.Enabled)
		assert.Equal(t, domain.LOCAL.String(), res.Availability.Type)
		assert.Equal(t, res.Availability.ID, "./tmp/dummy.pdf")
		assert.Empty(t, res.Availability.Url)
	})
}
