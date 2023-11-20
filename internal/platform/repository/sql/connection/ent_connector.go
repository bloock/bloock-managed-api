package connection

import (
	"github.com/bloock/bloock-managed-api/internal/platform/repository/sql/ent"
	"github.com/rs/zerolog"
)

type EntConnector struct {
	logger zerolog.Logger
}

func NewEntConnector(logger zerolog.Logger) *EntConnector {
	return &EntConnector{logger: logger}
}

func (c EntConnector) Connect(driver string, connectionURL string) (*ent.Client, error) {
	client, err := ent.Open(driver, connectionURL)
	if err != nil {
		c.logger.Error().Err(err).Msgf("failed opening connection for driver: %s %v", driver, err)
		return nil, err
	}

	return client, err
}
