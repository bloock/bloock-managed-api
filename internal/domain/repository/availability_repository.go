package repository

import (
	"context"

	"github.com/bloock/bloock-managed-api/internal/domain"
)

type AvailabilityRepository interface {
	UploadHosted(ctx context.Context, file *domain.File) (string, error)
	UploadIpfs(ctx context.Context, file *domain.File) (string, error)
	UploadLocal(ctx context.Context, file *domain.File) (string, error)
	UploadTmp(ctx context.Context, file *domain.File) (string, error)
	RetrieveTmp(ctx context.Context, filename string) ([]byte, error)
	FindFile(ctx context.Context, dataID string) ([]byte, error)
}
