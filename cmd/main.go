package main

import (
	"bloock-managed-api/internal/platform/repository"
	http_repository "bloock-managed-api/internal/platform/repository/http"
	"bloock-managed-api/internal/platform/repository/sql"
	"bloock-managed-api/internal/platform/repository/sql/connection"
	"bloock-managed-api/internal/platform/rest"
	"bloock-managed-api/internal/service/create"
	"bloock-managed-api/internal/service/update"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"net/http"
	"time"
)

func main() {
	v := viper.New()
	v.AddConfigPath(".")
	err := v.ReadInConfig()
	if _, configFileNotFound := err.(viper.ConfigFileNotFoundError); err != nil && !configFileNotFound {
		panic(err)
	}

	var cfg Config

	err = v.Unmarshal(&cfg)
	if err != nil {
		panic(err)
		return
	}

	logger := zerolog.Logger{}
	entConnector := connection.NewEntConnector(logger)
	conn, err := connection.NewEntConnection(cfg.DBConnectionString, entConnector, logger)
	if err != nil {
		panic(err)
		return
	}
	err = conn.Migrate()
	if err != nil {
		panic(err)
		return
	}
	certificationRepository := sql.NewSQLCertificationRepository(*conn, 5*time.Second, logger)
	integrityRepository := repository.NewBloockIntegrityRepository(cfg.BloockAPIKey, logger)
	notificationRepository := http_repository.NewHttpNotificationRepository(http.Client{}, cfg.WebhookURL, logger)

	createCertification := create.NewCertification(certificationRepository, integrityRepository)
	updateCertificationAnchor := update.NewCertificationAnchor(certificationRepository, notificationRepository)

	server := rest.NewServer(
		cfg.APIHost,
		cfg.APIPort,
		*createCertification,
		*updateCertificationAnchor,
		cfg.WebhookSecretKey,
		cfg.WebhookEnforceTolerance,
		logger,
	)
	err = server.Start()
	if err != nil {
		panic(err)
		return
	}
}

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
