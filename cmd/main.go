package main

import (
	"bloock-managed-api/internal/config"
	"bloock-managed-api/internal/platform/repository"
	http_repository "bloock-managed-api/internal/platform/repository/http"
	"bloock-managed-api/internal/platform/repository/sql"
	"bloock-managed-api/internal/platform/repository/sql/connection"
	"bloock-managed-api/internal/platform/rest"
	"bloock-managed-api/internal/service/create"
	"bloock-managed-api/internal/service/update"
	"github.com/rs/zerolog"
	"net/http"
	"time"
)

func main() {
	cfg, err := config.InitConfig()
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
	integrityRepository := repository.NewBloockIntegrityRepository(cfg.APIKey, logger)
	notificationRepository := http_repository.NewHttpNotificationRepository(http.Client{}, cfg.WebhookURL, logger)

	createCertification := create.NewCertification(certificationRepository, integrityRepository)
	updateCertificationAnchor := update.NewCertificationAnchor(certificationRepository, notificationRepository, integrityRepository)

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
