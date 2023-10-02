package file

import (
	"bloock-managed-api/internal/domain/repository"
	"context"

	"github.com/bloock/bloock-sdk-go/v2/client"
	"github.com/bloock/bloock-sdk-go/v2/entity/record"
)

type FileService struct {
	localStorageRepository repository.LocalStorageRepository
	recordClient           client.RecordClient
}

func NewFileService(localStorageRepository repository.LocalStorageRepository) *FileService {
	return &FileService{
		localStorageRepository: localStorageRepository,
		recordClient:           client.NewRecordClient(),
	}
}

func (f FileService) GetRecord(ctx context.Context, file []byte) (*record.Record, error) {
	record, err := f.recordClient.FromFile(file).Build()
	if err != nil {
		return nil, err
	}

	return &record, nil
}

func (f FileService) GetFileHash(ctx context.Context, file []byte) (string, error) {
	record, err := f.recordClient.FromFile(file).Build()
	if err != nil {
		return "", err
	}
	hash, err := record.GetHash()
	if err != nil {
		return "", err
	}

	return hash, nil
}

func (f FileService) SaveFile(ctx context.Context, file []byte, hash string) error {
	return f.localStorageRepository.Save(ctx, file, hash)
}
