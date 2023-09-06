package repository

import "context"

type AvailabilityRepository interface {
	UploadHosted(ctx context.Context, data []byte) (string, error)
	UploadIpfs(ctx context.Context, data []byte) (string, error)
	FindFile(ctx context.Context, dataID string) ([]byte, error)
}
