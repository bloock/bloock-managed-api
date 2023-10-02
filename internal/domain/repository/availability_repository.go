package repository

import (
	"context"

	"github.com/bloock/bloock-sdk-go/v2/entity/record"
)

type AvailabilityRepository interface {
	UploadHosted(ctx context.Context, record *record.Record) (string, error)
	UploadIpfs(ctx context.Context, record *record.Record) (string, error)
	FindFile(ctx context.Context, dataID string) ([]byte, error)
}
