package main

import (
	"bloock-managed-api/internal/config"
	"bloock-managed-api/internal/platform/repository"
	"bloock-managed-api/internal/platform/repository/hard_drive"
	http_repository "bloock-managed-api/internal/platform/repository/http"
	"bloock-managed-api/internal/platform/repository/sql"
	"bloock-managed-api/internal/platform/repository/sql/connection"
	"bloock-managed-api/internal/platform/rest"
	"bloock-managed-api/internal/service/authenticity"
	"bloock-managed-api/internal/service/availability"
	"bloock-managed-api/internal/service/file"
	"bloock-managed-api/internal/service/integrity"
	"bloock-managed-api/internal/service/notify"
	"bloock-managed-api/internal/service/process"
	"github.com/bloock/bloock-sdk-go/v2"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

func main() {
	os.Setenv("BLOOCK_API_FILE_DIR", "tmp")

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

	bloock.ApiKey = cfg.APIKey
	bloock.ApiHost = "https://api.bloock.dev"

	certificationRepository := sql.NewSQLCertificationRepository(*conn, 5*time.Second, logger)
	integrityRepository := repository.NewBloockIntegrityRepository(logger)
	notificationRepository := http_repository.NewHttpNotificationRepository(http.Client{}, cfg.ClientEndpointUrl, logger)
	authenticityRepository := repository.NewBloockAuthenticityRepository(logger)
	availabilityRepository := repository.NewBloockAvailabilityRepository(logger)
	storageRepository := hard_drive.NewHardDriveLocalStorageRepository(cfg.FileDir, logger)

	integrityService := integrity.NewIntegrityService(certificationRepository, integrityRepository)
	authenticityService := authenticity.NewAuthenticityService(authenticityRepository)
	updateAnchorService := integrity.NewUpdateAnchorService(certificationRepository)
	availabilityService := availability.NewAvailabilityService(availabilityRepository)
	fileService := file.NewFileService(storageRepository)
	notifyService := notify.NewNotifyService(notificationRepository, storageRepository, availabilityRepository)
	processService := process.NewProcessService(integrityService, authenticityService, availabilityService, fileService, notifyService)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		server, err := rest.NewServer(cfg.APIHost, cfg.APIPort, processService, updateAnchorService, notifyService, cfg.WebhookSecretKey, cfg.WebhookEnforceTolerance, logger, cfg.DebugMode)
		if err != nil {
			panic(err)
		}
		err = server.Start()
		if err != nil {
			panic(err)
		}
	}()

	wg.Wait()
}
