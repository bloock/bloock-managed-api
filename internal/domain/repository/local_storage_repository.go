package repository

import "context"

type LocalStorageRepository interface {
	Save(ctx context.Context, data []byte, filename string) error
	Retrieve(ctx context.Context, directory string, filename string) ([]byte, error)
}
