package config

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitConfig(t *testing.T) {

	t.Run("given env vars it should populate config", func(t *testing.T) {
		value := "1234abcd"
		intVal := "1234"
		_ = os.Setenv("BLOOCK_API_PORT", value)
		_ = os.Setenv("BLOOCK_API_HOST", value)
		_ = os.Setenv("BLOOCK_API_KEY", value)
		_ = os.Setenv("BLOOCK_CLIENT_ENDPOINT_URL", value)
		_ = os.Setenv("BLOOCK_WEBHOOK_SECRET_KEY", value)
		_ = os.Setenv("BLOOCK_ENFORCE_TOLERANCE", "0")
		_ = os.Setenv("BLOOCK_DB_CONNECTION_STRING", value)
		_ = os.Setenv("BLOOCK_API_DEBUG_MODE", "true")
		_ = os.Setenv("BLOOCK_AUTHENTICITY_PRIVATE_KEY", value)
		_ = os.Setenv("BLOOCK_AUTHENTICITY_PUBLIC_KEY", value)
		_ = os.Setenv("BLOOCK_MAX_MEMORY", intVal)
		_ = os.Setenv("BLOOCK_FILE_DIR", value)
		config, err := InitConfig()

		assert.NoError(t, err)
		assert.NotEmpty(t, config)
		assert.Equal(t, value, config.APIHost)
		assert.Equal(t, value, config.APIKey)
		assert.Equal(t, value, config.APIPort)
		assert.Equal(t, value, config.DBConnectionString)
		assert.Equal(t, value, config.ClientEndpointUrl)
		assert.Equal(t, value, config.WebhookSecretKey)
		assert.Equal(t, value, config.PublicKey)
		assert.Equal(t, value, config.PrivateKey)
		atoi, _ := strconv.Atoi(intVal)
		assert.Equal(t, int64(atoi), config.MaxMemory)
		assert.Equal(t, value, config.FileDir)
		assert.Equal(t, true, config.DebugMode)
		assert.Equal(t, false, config.WebhookEnforceTolerance)

	})

}
