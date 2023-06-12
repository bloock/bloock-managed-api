package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestInitConfig(t *testing.T) {

	t.Run("given no config file path it should read from env", func(t *testing.T) {
		value := "1234abcd"
		err := os.Setenv("BLOOCK_API_PORT", value)
		require.NoError(t, err)
		err = os.Setenv("BLOOCK_API_HOST", value)
		err = os.Setenv("BLOOCK_API_KEY", value)
		err = os.Setenv("BLOOCK_WEBHOOK_URL", value)
		err = os.Setenv("BLOOCK_WEBHOOK_SECRET_KEY", value)
		err = os.Setenv("BLOOCK_ENFORCE_TOLERANCE_STRING", "0")
		err = os.Setenv("BLOOCK_DB_CONNECTION_STRING", value)
		config, err := InitConfig()

		assert.NotEmpty(t, config)
		assert.Equal(t, value, config.APIHost)
		assert.Equal(t, value, config.APIKey)
		assert.Equal(t, value, config.APIPort)
		assert.Equal(t, value, config.DBConnectionString)
		assert.Equal(t, value, config.WebhookURL)
		assert.Equal(t, value, config.WebhookSecretKey)
		assert.Equal(t, false, config.WebhookEnforceTolerance)
		assert.NoError(t, err)

	})

	t.Run("given config path it should read from config file", func(t *testing.T) {
		err := os.Setenv("BLOOCK_CONFIG_PATH", "../../")
		require.NoError(t, err)

		config, err := InitConfig()

		assert.NotEmpty(t, config)
		assert.NoError(t, err)
	})

}
