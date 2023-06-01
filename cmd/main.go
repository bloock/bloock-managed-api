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
	"net/http"
	"time"
)

func main() {
	logger := zerolog.Logger{}
	entConnector := connection.NewEntConnector(logger)
	conn, err := connection.NewEntConnection("file:ent?mode=memory&cache=shared&_fk=1", entConnector, logger)
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
	integrityRepository := repository.NewBloockIntegrityRepository("Nm1sFmrojcrRgfZ4v0H0w0d1d22GookjcJl7y-2jr51qx0RioCR3nVm1z74hDEzZ", logger)
	notificationRepository := http_repository.NewHttpNotificationRepository(http.Client{}, "http://localhost:8081/v1/certification", logger)

	createCertification := create.NewCertification(certificationRepository, integrityRepository)
	updateCertificationAnchor := update.NewCertificationAnchor(certificationRepository, notificationRepository)

	server := rest.NewServer("", "8081", *createCertification, *updateCertificationAnchor, 5, logger)
	err = server.Start()
	if err != nil {
		panic(err)
		return
	}
}
