package file

import (
	"bloock-managed-api/internal/domain/repository"
	"context"
	"github.com/bloock/bloock-sdk-go/v2/client"
)

type FileService struct {
	localStorageRepository repository.LocalStorageRepository
	recordClient           client.RecordClient
}

func NewFileService(localStorageRepository repository.LocalStorageRepository, recordClient client.RecordClient) *FileService {
	return &FileService{localStorageRepository: localStorageRepository, recordClient: recordClient}
}

func (f FileService) SaveFile(ctx context.Context, file []byte) error {
	record, err := f.recordClient.FromFile(file).Build()
	if err != nil {
		return err
	}
	hash, err := record.GetHash()
	if err != nil {
		return err
	}

	return f.localStorageRepository.Save(ctx, file, hash)
}
