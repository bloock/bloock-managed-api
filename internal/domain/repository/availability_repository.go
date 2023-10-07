package repository

import (
	"context"

	"github.com/bloock/bloock-sdk-go/v2/entity/record"
)

type AvailabilityRepository interface {
	UploadHosted(ctx context.Context, record *record.Record) (string, error)
	UploadIpfs(ctx context.Context, record *record.Record) (string, error)
	UploadLocal(ctx context.Context, record *record.Record) (string, error)
	UploadTmp(ctx context.Context, record *record.Record) (string, error)
	RetrieveTmp(ctx context.Context, filename string) ([]byte, error)
	FindFile(ctx context.Context, dataID string) ([]byte, error)
}
