package main

import (
	"bloock-managed-api/internal/config"
	"bloock-managed-api/internal/platform/repository"
	http_repository "bloock-managed-api/internal/platform/repository/http"
	"bloock-managed-api/internal/platform/repository/sql"
	"bloock-managed-api/internal/platform/repository/sql/connection"
	"bloock-managed-api/internal/platform/rest"
	"bloock-managed-api/internal/service/authenticity"
	"bloock-managed-api/internal/service/availability"
	"bloock-managed-api/internal/service/integrity"
	"bloock-managed-api/internal/service/process"
	"github.com/bloock/bloock-sdk-go/v2/client"
	"github.com/rs/zerolog"
	"net/http"
	"os"
	"time"
)

func main() {
	cfg, err := config.InitConfig()
	if err != nil {
		panic(err)
	}
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
	entConnector := connection.NewEntConnector(logger)
	conn, err := connection.NewEntConnection(cfg.DBConnectionString, entConnector, logger)
	if err != nil {
		panic(err)
	}
	err = conn.Migrate()
	if err != nil {
		panic(err)
	}

	recordClient := client.NewRecordClient()
	availabilityClient := client.NewAvailabilityClient()
	keyClient := client.NewKeyClient()
	authenticityClient := client.NewAuthenticityClient()
	integrityClient := client.NewIntegrityClient()

	certificationRepository := sql.NewSQLCertificationRepository(*conn, 5*time.Second, logger)
	integrityRepository := repository.NewBloockIntegrityRepository(integrityClient, logger)
	notificationRepository := http_repository.NewHttpNotificationRepository(http.Client{}, cfg.WebhookURL, logger)
	authenticityRepository := repository.NewBloockAuthenticityRepository(keyClient, authenticityClient, recordClient, logger)
	availabilityRepository := repository.NewBloockAvailabilityRepository(recordClient, availabilityClient, logger)

	integrityService := integrity.NewIntegrityService(certificationRepository, integrityRepository)
	authenticityService := authenticity.NewAuthenticityService(authenticityRepository)
	updateAnchorService := integrity.NewUpdateAnchorService(certificationRepository, notificationRepository, integrityRepository)
	availabilityService := availability.NewAvailabilityService(availabilityRepository)
	processService := process.NewProcessService(integrityService, authenticityService, availabilityService)

	server, err := rest.NewServer(cfg.APIHost, cfg.APIPort, processService, updateAnchorService, cfg.WebhookSecretKey, cfg.WebhookEnforceTolerance, logger, cfg.DebugMode)
	if err != nil {
		panic(err)
	}
	err = server.Start()
	if err != nil {
		panic(err)
	}
}
