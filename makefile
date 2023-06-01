mocks:
	mockgen -source=integrity/domain/repository/certification_repository.go -destination=integrity/domain/repository/mocks/mock_certification_repository.go
	mockgen -source=integrity/domain/repository/integrity_repository.go -destination=integrity/domain/repository/mocks/mock_integrity_repository.go
	mockgen -source=integrity/domain/repository/notification_repository.go -destination=integrity/domain/repository/mocks/mock_notification_repository.go
	mockgen -source=integrity/platform/repository/sql/connection/sql_connector.go -destination=integrity/platform/repository/sql/connection/mocks/mock_sql_connector.go

schema:
	go generate ./ent
