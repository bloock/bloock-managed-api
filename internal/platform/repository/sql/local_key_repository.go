package sql

import (
	"bloock-managed-api/internal/domain"
	"bloock-managed-api/internal/platform/repository/sql/connection"
	"bloock-managed-api/internal/platform/repository/sql/ent"
	"bloock-managed-api/internal/service"
	"context"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"time"
)

type SQLLocalKeyRepository struct {
	connection connection.EntConnection
	dbTimeout  time.Duration
	logger     zerolog.Logger
}

func (s SQLLocalKeyRepository) SaveKey(ctx context.Context, localKey domain.LocalKey) error {
	key := localKey.LocalKey()
	_, err := s.connection.DB().LocalKey.Create().
		SetID(localKey.Id()).
		SetLocalKey(&key).
		SetKeyType(localKey.KeyTypeStr()).Save(ctx)

	if err != nil {
		s.logger.Error().Err(err).Msg("")
		return err
	}

	return nil
}

func (s SQLLocalKeyRepository) FindKeyByID(ctx context.Context, id uuid.UUID) (*domain.LocalKey, error) {
	localKey, err := s.connection.DB().LocalKey.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	keyType, err := service.ValidateKeyType(localKey.KeyType)
	if err != nil {
		return nil, err
	}

	return domain.NewLocalKey(*localKey.LocalKey, keyType, localKey.ID), nil
}

func (s SQLLocalKeyRepository) FindKeys(ctx context.Context) ([]domain.LocalKey, error) {
	localKeys, err := s.connection.DB().LocalKey.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	return toLocalKeys(localKeys)
}

func toLocalKeys(keys []*ent.LocalKey) ([]domain.LocalKey, error) {
	var localKeys []domain.LocalKey
	for _, localKey := range keys {
		keyType, err := service.ValidateKeyType(localKey.KeyType)
		if err != nil {
			return nil, err
		}
		localKeys = append(localKeys, *domain.NewLocalKey(*localKey.LocalKey, keyType, localKey.ID))
	}

	return localKeys, nil
}

func NewSQLLocalKeyRepository(connection connection.EntConnection, dbTimeout time.Duration, logger zerolog.Logger) *SQLLocalKeyRepository {
	return &SQLLocalKeyRepository{connection: connection, dbTimeout: dbTimeout, logger: logger}
}
