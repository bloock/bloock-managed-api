package main

import (
	"bloock-managed-api/internal/config"
	mock_repository "bloock-managed-api/internal/domain/repository/mocks"
	"bloock-managed-api/internal/platform/repository"
	http_repository "bloock-managed-api/internal/platform/repository/http"
	"bloock-managed-api/internal/platform/repository/sql"
	"bloock-managed-api/internal/platform/repository/sql/connection"
	"bloock-managed-api/internal/platform/rest"
	"bloock-managed-api/internal/service/create"
	"bloock-managed-api/internal/service/get"
	"bloock-managed-api/internal/service/update"
	"github.com/bloock/bloock-sdk-go/v2/client"
	"github.com/rs/zerolog"
	"net/http"
	"time"
)

func main() {
	cfg, err := config.InitConfig()
	if err != nil {
		panic(err)
	}

	logger := zerolog.Logger{}
	entConnector := connection.NewEntConnector(logger)
	conn, err := connection.NewEntConnection(cfg.DBConnectionString, entConnector, logger)
	if err != nil {
		panic(err)
	}
	err = conn.Migrate()
	if err != nil {
		panic(err)
	}
	certificationRepository := sql.NewSQLCertificationRepository(*conn, 5*time.Second, logger)
	integrityRepository := repository.NewBloockIntegrityRepository(cfg.APIKey, logger)
	notificationRepository := http_repository.NewHttpNotificationRepository(http.Client{}, cfg.WebhookURL, logger)
	authenticityRepository := repository.NewBloockAuthenticityRepository(
		cfg.APIKey,
		client.NewKeyClient(),
		client.NewAuthenticityClient(),
		client.NewRecordClient(),
		logger,
	)
	localKeyRepository := mock_repository.NewMockLocalKeysRepository(nil)

	createCertification := create.NewCertification(certificationRepository, integrityRepository)
	updateCertificationAnchor := update.NewCertificationAnchor(certificationRepository, notificationRepository, integrityRepository)
	createManagedKey := create.NewManagedKey(authenticityRepository)
	createLocalKey := create.NewLocalKey(authenticityRepository, localKeyRepository)
	createSignature := create.NewSignature(authenticityRepository, localKeyRepository)
	getKeys := get.NewLocalKeys(localKeyRepository)
	server, err := rest.NewServer(cfg.APIHost, cfg.APIPort, getKeys, createManagedKey, createLocalKey, createSignature, *createCertification, *updateCertificationAnchor, cfg.WebhookSecretKey, cfg.WebhookEnforceTolerance, logger, cfg.DebugMode)
	if err != nil {
		panic(err)
	}
	err = server.Start()
	if err != nil {
		panic(err)
	}
}
