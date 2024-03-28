package process

import (
	"context"
	"github.com/bloock/bloock-managed-api/internal/domain"
	"github.com/bloock/bloock-managed-api/internal/domain/repository"
	bloock_repository "github.com/bloock/bloock-managed-api/internal/platform/repository"
	"github.com/bloock/bloock-managed-api/internal/platform/repository/sql/connection"
	"github.com/bloock/bloock-managed-api/internal/platform/utils"
	"github.com/rs/zerolog"
)

type ListProcess struct {
	processRepository repository.ProcessRepository
	logger            zerolog.Logger
}

func NewListProcess(ctx context.Context, l zerolog.Logger, ent *connection.EntConnection) *ListProcess {
	logger := l.With().Caller().Str("component", "list-process").Logger()

	return &ListProcess{
		processRepository: bloock_repository.NewProcessRepository(ctx, l, ent),
		logger:            logger,
	}
}

func (l ListProcess) List(ctx context.Context, pq utils.PaginationQuery) ([]domain.Process, utils.Pagination, error) {
	listProcess, pagination, err := l.processRepository.List(ctx, pq)
	if err != nil {
		return []domain.Process{}, utils.Pagination{}, err
	}

	return listProcess, pagination, nil
}
