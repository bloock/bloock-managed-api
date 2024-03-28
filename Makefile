mocks:
	mockgen -source=internal/service/application_service.go -destination=internal/service/mock/mock_application_service.go
	mockgen -source=internal/domain/repository/integrity_repository.go -destination=internal/domain/repository/mocks/mock_integrity_repository.go
	mockgen -source=internal/domain/repository/authenticity_repository.go -destination=internal/domain/repository/mocks/mock_authenticity_repository.go
	mockgen -source=internal/domain/repository/availability_repository.go -destination=internal/domain/repository/mocks/mock_availability_repository.go
	mockgen -source=internal/domain/repository/encryption_repository.go -destination=internal/domain/repository/mocks/mock_encryption_repository.go
	mockgen -source=internal/domain/repository/notification_repository.go -destination=internal/domain/repository/mocks/mock_notification_repository.go
	mockgen -source=internal/domain/repository/metadata_repository.go -destination=internal/domain/repository/mocks/mock_metadata_repository.go
	mockgen -source=internal/platform/repository/sql/connection/sql_connector.go -destination=internal/platform/repository/sql/connection/mocks/mock_sql_connector.go
	mockgen -source=internal/domain/repository/merlkle_tree_repository.go -destination=internal/domain/repository/mocks/mock_merlkle_tree_repository.go
	mockgen -source=internal/domain/repository/message_aggregator_repository.go -destination=internal/domain/repository/mocks/mock_message_aggregator_repository.go
	mockgen -source=internal/domain/repository/process_repository.go -destination=internal/domain/repository/mocks/mock_process_repository.go

schemas:
	go run -mod=mod entgo.io/ent/cmd/ent generate --target ./internal/platform/repository/sql/ent ./internal/platform/repository/sql/schema

down:
	docker compose down

up:
	docker-compose up -d --build

test:
	go test -v ./...

cache:
	go mod tidy
	go mod vendor
qa:
	staticcheck ./...
	go vet ./...