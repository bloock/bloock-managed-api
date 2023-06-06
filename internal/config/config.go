package config

import (
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	DBConnectionString      string
	BloockAPIKey            string
	WebhookURL              string
	APIHost                 string
	APIPort                 string
	MaxMemoryPerRequest     string
	WebhookEnforceTolerance bool
	WebhookSecretKey        string
}

func InitConfig() (*Config, error) {
	v := viper.New()
	var in string
	if in = os.Getenv("CONFIG_PATH"); in == "" {
		in = "."
	}
	v.AddConfigPath(in)
	err := v.ReadInConfig()
	if _, configFileNotFound := err.(viper.ConfigFileNotFoundError); err != nil && !configFileNotFound {

		return &Config{}, err
	}

	var cfg Config
	err = v.Unmarshal(&cfg)
	if err != nil {
		return &Config{}, err
	}

	return &cfg, nil
}
