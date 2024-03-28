package process

import (
	"context"
	"errors"
	"github.com/bloock/bloock-managed-api/internal/domain"
	"github.com/bloock/bloock-managed-api/internal/domain/repository"
	bloock_repository "github.com/bloock/bloock-managed-api/internal/platform/repository"
	"github.com/bloock/bloock-managed-api/internal/platform/repository/sql/connection"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

var (
	ErrInvalidUUID     = errors.New("invalid uuid provided")
	ErrProcessNotFound = errors.New("process not found")
)

type GetProcessByID struct {
	processRepository repository.ProcessRepository
	logger            zerolog.Logger
}

func NewGetProcessByID(ctx context.Context, l zerolog.Logger, ent *connection.EntConnection) *GetProcessByID {
	logger := l.With().Caller().Str("component", "get-process-by-id").Logger()

	return &GetProcessByID{
		processRepository: bloock_repository.NewProcessRepository(ctx, l, ent),
		logger:            logger,
	}
}

func (g GetProcessByID) Get(ctx context.Context, processID string) (domain.Process, error) {
	processUUID, err := uuid.Parse(processID)
	if err != nil {
		g.logger.Error().Err(err).Msg("")
		return domain.Process{}, ErrInvalidUUID
	}

	process, err := g.processRepository.FindProcessByID(ctx, processUUID)
	if err != nil {
		return domain.Process{}, err
	}
	if process.Hash == "" {
		err = ErrProcessNotFound
		g.logger.Error().Err(err).Msg("")
		return domain.Process{}, err
	}

	return process, nil
}
