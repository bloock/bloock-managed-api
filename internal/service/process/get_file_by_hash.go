package process

import (
	"context"
	"errors"
	"github.com/bloock/bloock-managed-api/internal/domain/repository"
	bloock_repository "github.com/bloock/bloock-managed-api/internal/platform/repository"
	"github.com/bloock/bloock-managed-api/internal/platform/repository/sql/connection"
	"github.com/rs/zerolog"
)

var ErrHashNotFound = errors.New("hash not found")

type GetFileByHash struct {
	metadataRepository     repository.MetadataRepository
	availabilityRepository repository.AvailabilityRepository
	logger                 zerolog.Logger
}

func NewGetFileByHash(ctx context.Context, l zerolog.Logger, ent *connection.EntConnection) *GetFileByHash {
	logger := l.With().Caller().Str("component", "get-file-by-hash").Logger()

	return &GetFileByHash{
		metadataRepository:     bloock_repository.NewBloockMetadataRepository(ctx, l, ent),
		availabilityRepository: bloock_repository.NewBloockAvailabilityRepository(ctx, l),
		logger:                 logger,
	}
}

func (s GetFileByHash) Get(ctx context.Context, hash string) ([]byte, error) {
	certification, err := s.metadataRepository.FindCertificationByHash(ctx, hash)
	if err != nil {
		return nil, err
	}

	var fileBytes []byte
	if certification.DataID != "" {
		fileBytes, err = s.availabilityRepository.FindFile(ctx, certification.DataID)
		if err != nil {
			fileBytes, err = s.availabilityRepository.RetrieveLocal(ctx, certification.DataID)
			if err != nil {
				return nil, ErrHashNotFound
			}
		}
	} else {
		fileBytes, err = s.availabilityRepository.RetrieveTmp(ctx, hash)
		if err != nil {
			return nil, ErrHashNotFound
		}
	}

	return fileBytes, nil
}
