package config

import (
	"encoding/json"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"os"
	"strings"
)

type Config struct {
	DBConnectionString      string `json:"bloock_db_connection_string" default:"file:bloock?mode=memory&cache=shared&_fk=1"`
	APIKey                  string `json:"bloock_api_key"`
	APIHost                 string `json:"bloock_api_host"`
	APIPort                 string `json:"bloock_api_port"`
	WebhookURL              string `json:"bloock_webhook_url"`
	WebhookSecretKey        string `json:"bloock_webhook_secret_key"`
	WebhookEnforceTolerance bool   `json:"bloock_webhook_enforce_tolerance"`
}

func InitConfig() (*Config, error) {
	logger := zerolog.Logger{}
	v := viper.New()
	var cfg = &Config{}

	cfgPath := os.Getenv("BLOOCK_CONFIG_PATH")
	if cfgPath == "" {
		logger.Info().Msg("reading configuration from env")
		err := readFromEnv(cfg)
		if err != nil {
			return &Config{}, err
		}

		return cfg, err
	}
	logger.Info().Msgf("reading configuration from config file: %s", cfgPath)

	v.AddConfigPath(cfgPath)
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	err := v.ReadInConfig()
	if _, configFileNotFound := err.(viper.ConfigFileNotFoundError); err != nil || configFileNotFound {

		return &Config{}, err
	}

	err = v.Unmarshal(cfg)
	if err != nil {
		return &Config{}, err
	}

	return cfg, nil
}

func readFromEnv(cfg *Config) error {
	envMap := make(map[string]string)

	envVariables := os.Environ()

	for _, env := range envVariables {
		pair := strings.SplitN(env, "=", 2)
		key := pair[0]
		value := pair[1]

		if strings.HasPrefix(key, "BLOOCK") {
			envMap[key] = value
		}
	}

	bytes, err := json.Marshal(envMap)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, cfg)
	if err != nil {
		return err
	}

	return nil
}
