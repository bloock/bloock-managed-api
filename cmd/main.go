package main

import (
	"bloock-managed-api/internal/config"
	"bloock-managed-api/internal/platform/repository"
	http_repository "bloock-managed-api/internal/platform/repository/http"
	"bloock-managed-api/internal/platform/repository/sql/connection"
	"bloock-managed-api/internal/platform/rest"
	"bloock-managed-api/internal/service/notify"
	"bloock-managed-api/internal/service/process"
	"time"

	"github.com/bloock/bloock-sdk-go/v2"

	"net/http"
	"os"
	"sync"

	"github.com/rs/zerolog"
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

	bloock.ApiKey = cfg.APIKey
	bloock.DisableAnalytics = true

	integrityRepository := repository.NewBloockIntegrityRepository(logger)
	authenticityRepository := repository.NewBloockAuthenticityRepository(logger)
	encryptionRepository := repository.NewBloockEncryptionRepository(logger)
	availabilityRepository := repository.NewBloockAvailabilityRepository(cfg.StorageLocalPath, cfg.TmpDir, logger)
	metadataRepository := repository.NewBloockMetadataRepository(*conn, 5*time.Second, logger)
	notificationRepository := http_repository.NewHttpNotificationRepository(http.Client{}, cfg.ClientEndpointUrl, logger)

	processService := process.NewProcessService(integrityRepository, authenticityRepository, encryptionRepository, availabilityRepository, metadataRepository, notificationRepository)
	notifyServie := notify.NewNotifyService(availabilityRepository, metadataRepository, notificationRepository)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		server, err := rest.NewServer(cfg.APIHost, cfg.APIPort, processService, notifyServie, cfg.WebhookSecretKey, logger, cfg.DebugMode)
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
