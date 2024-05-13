package config

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/bloock/bloock-sdk-go/v2"
	"github.com/mcuadros/go-defaults"
	"github.com/spf13/viper"
)

type APIConfig struct {
	Host      string `mapstructure:"host" default:"0.0.0.0"`
	Port      string `mapstructure:"port" default:"8080"`
	DebugMode bool   `mapstructure:"debug_mode" default:"false"`
}

type AuthConfig struct {
	Secret string `mapstructure:"secret"`
}

type DBConfig struct {
	ConnectionString string `mapstructure:"connection_string" default:"file:managed?mode=memory&cache=shared&_fk=1"`
}

type BloockConfig struct {
	ApiHost          string `mapstructure:"api_host" default:"https://api.bloock.com"`
	ApiKey           string `mapstructure:"api_key"`
	CdnHost          string `mapstructure:"cdn_host" default:"https://cdn.bloock.com"`
	WebhookSecretKey string `mapstructure:"webhook_secret_key"`
}

type WebhookConfig struct {
	ClientEndpointUrl string `mapstructure:"client_endpoint_url"`
	MaxRetries        uint64 `mapstructure:"max_retries"`
}

type KeyConfig struct {
	KeyType string `mapstructure:"key_type"`
	Key     string `mapstructure:"key"`
}

type CertificateConfig struct {
	Pkcs12Path     string `mapstructure:"pkcs12_path"`
	Pkcs12Password string `mapstructure:"pkcs12_password"`
}

type AuthenticityConfig struct {
	KeyConfig         KeyConfig         `mapstructure:"key"`
	CertificateConfig CertificateConfig `mapstructure:"certificate"`
}

type EncryptionConfig struct {
	KeyConfig         KeyConfig         `mapstructure:"key"`
	CertificateConfig CertificateConfig `mapstructure:"certificate"`
}

type StorageConfig struct {
	TmpDir        string `mapstructure:"tmp_dir" default:"./tmp"`
	LocalPath     string `mapstructure:"local_path" default:"./data"`
	LocalStrategy string `mapstructure:"local_strategy" default:"HASH"`
}

type IntegrityConfig struct {
	AggregateMode       bool `mapstructure:"aggregate_mode" default:"false"`
	AggregateWorker     bool `mapstructure:"aggregate_worker" default:"false"`
	AggregateInterval   int  `mapstructure:"aggregate_interval" default:"3600"`
	MaxProofMessageSize int  `mapstructure:"max_proof_message_size" default:"1000"`
}

type Config struct {
	Api          APIConfig
	Auth         AuthConfig
	Db           DBConfig
	Bloock       BloockConfig
	Webhook      WebhookConfig
	Authenticity AuthenticityConfig
	Encryption   EncryptionConfig
	Storage      StorageConfig
	Integrity    IntegrityConfig
}

var Configuration = Config{}

func InitConfig() (*Config, error) {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("bloock")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.AddConfigPath("./")
	err := viper.ReadInConfig()
	if err != nil {
		switch err.(type) {
		default:
			return nil, fmt.Errorf("fatal error loading config file: %s", err)
		case viper.ConfigFileNotFoundError:
			return nil, errors.New("No config file found. Using defaults and environment variables")
		}
	}

	bindEnvs(Configuration)

	err = viper.Unmarshal(&Configuration)
	if err != nil {
		return nil, fmt.Errorf("fatal error loading config file: %s", err)
	}
	defaults.SetDefaults(&Configuration)

	bloock.ApiHost = Configuration.Bloock.ApiHost

	if Configuration.Integrity.AggregateMode {
		if Configuration.Bloock.ApiKey == "" {
			return nil, errors.New("aggregate mode requires a BLOOCK Api Key set")
		}
	}

	return &Configuration, nil
}

func bindEnvs(iface interface{}, parts ...string) {
	ifv := reflect.ValueOf(iface)
	ift := reflect.TypeOf(iface)
	for i := 0; i < ift.NumField(); i++ {
		v := ifv.Field(i)
		t := ift.Field(i)

		var tv string
		tv, ok := t.Tag.Lookup("mapstructure")
		if !ok {
			tv = toSnakeCase(t.Name)
		}
		switch v.Kind() {
		case reflect.Struct:
			bindEnvs(v.Interface(), append(parts, tv)...)
		default:
			viper.BindEnv(strings.Join(append(parts, tv), "."))
		}
	}
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
