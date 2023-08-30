package config

import (
	"github.com/bloock/bloock-sdk-go/v2"
	"github.com/spf13/viper"
)

type Config struct {
	DBConnectionString      string `mapstructure:"BLOOCK_DB_CONNECTION_STRING"`
	APIKey                  string `mapstructure:"BLOOCK_API_KEY"`
	APIHost                 string `mapstructure:"BLOOCK_API_HOST"`
	APIPort                 string `mapstructure:"BLOOCK_API_PORT"`
	WebhookURL              string `mapstructure:"BLOOCK_WEBHOOK_URL"`
	WebhookSecretKey        string `mapstructure:"BLOOCK_WEBHOOK_SECRET_KEY"`
	WebhookEnforceTolerance bool   `mapstructure:"BLOOCK_ENFORCE_TOLERANCE"`
	DebugMode               bool   `mapstructure:"BLOOCK_API_DEBUG_MODE"`
	PrivateKey              string `mapstructure:"BLOOCK_AUTHENTICITY_PRIVATE_KEY"`
	PublicKey               string `mapstructure:"BLOOCK_AUTHENTICITY_PUBLIC_KEY"`
	MaxMemory               int64  `mapstructure:"BLOOCK_MAX_MEMORY"`
	FileDir                 string `mapstructure:"BLOOCK_FILE_DIR"`
}

var Configuration = &Config{}

func InitConfig() (*Config, error) {

	setDefaultConfigValues()
	viper.AddConfigPath("./")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()

	err = viper.Unmarshal(Configuration)
	if err != nil {
		return &Config{}, err
	}

	bloock.ApiKey = Configuration.APIKey
	return Configuration, nil
}

func setDefaultConfigValues() {
	viper.SetDefault("bloock_db_connection_string", "file:managed?mode=memory&cache=shared&_fk=1")
	viper.SetDefault("bloock_api_key", "")
	viper.SetDefault("bloock_webhook_url", "")
	viper.SetDefault("bloock_webhook_secret_key", "")
	viper.SetDefault("bloock_api_host", "0.0.0.0")
	viper.SetDefault("bloock_api_port", "8080")
	viper.SetDefault("bloock_webhook_enforce_tolerance", false)
	viper.SetDefault("bloock_api_debug_mode", false)
	viper.SetDefault("bloock_max_memory", 10<<20) //10MB
	viper.SetDefault("bloock_file_dir", "./")
	viper.SetDefault("bloock_authenticity_private_key", "./")
	viper.SetDefault("bloock_authenticity_public_key", "./")
}
