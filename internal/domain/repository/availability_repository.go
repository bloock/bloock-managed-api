package repository

import (
	"context"

	"github.com/bloock/bloock-managed-api/internal/domain"
	"github.com/bloock/bloock-sdk-go/v2/entity/record"
)

type AvailabilityRepository interface {
	UploadHosted(ctx context.Context, file *domain.File, record record.Record) (string, error)
	UploadIpfs(ctx context.Context, file *domain.File, record record.Record) (string, error)
	UploadLocal(ctx context.Context, file *domain.File) (string, error)
	UploadTmp(ctx context.Context, file *domain.File, record record.Record) (string, error)
	RetrieveTmp(ctx context.Context, filename string) ([]byte, error)
	RetrieveLocal(ctx context.Context, filePath string) ([]byte, error)
	FindFile(ctx context.Context, dataID string) ([]byte, error)
	FindAll(ctx context.Context, id string) ([]domain.File, error)
}
