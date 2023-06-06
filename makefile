mocks:
	mockgen -source=internal/domain/repository/certification_repository.go -destination=internal/domain/repository/mocks/mock_certification_repository.go
	mockgen -source=internal/domain/repository/integrity_repository.go -destination=internal/domain/repository/mocks/mock_integrity_repository.go
	mockgen -source=internal/domain/repository/notification_repository.go -destination=internal/domain/repository/mocks/mock_notification_repository.go
	mockgen -source=internal/platform/repository/sql/connection/sql_connector.go -destination=internal/platform/repository/sql/connection/mocks/mock_sql_connector.go

schema:
	go run -mod=mod entgo.io/ent/cmd/ent generate --target ./internal/platform/repository/sql/ent ./internal/platform/repository/sql/schema

up:
	docker-compose up -d --build

test:
	go test -v ./...