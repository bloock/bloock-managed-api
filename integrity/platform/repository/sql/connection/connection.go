package connection

import (
	"bloock-managed-api/ent"
	"context"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"
	"strings"
)

const (
	Mysql    = "mysql"
	Postgres = "postgres"
	Sqlite   = "sqlite3"
)

type Connection struct {
	db     *ent.Client
	logger zerolog.Logger
}

func NewConnection(connectionURL string, connector SQLConnector, logger zerolog.Logger) (*Connection, error) {
	if connectionURL == "" {
		return &Connection{}, errors.New("connectionURL cannot be empty")
	}

	if strings.Contains(connectionURL, "file") {
		client, err := open(connector, Sqlite, connectionURL)
		if err != nil {
			return nil, err
		}
		return &Connection{
			db: client,
		}, nil
	}
	if strings.Contains(connectionURL, "mysql") {
		client, err := open(connector, Mysql, connectionURL)
		if err != nil {
			return nil, err
		}
		return &Connection{
			db: client,
		}, nil
	}
	if strings.Contains(connectionURL, "postgres") {
		client, err := open(connector, Postgres, connectionURL)
		if err != nil {
			return nil, err
		}
		return &Connection{
			db: client,
		}, nil
	}

	err := errors.New("unsupported database")
	logger.Error().Err(err).Msgf(" with url: %s", connectionURL)
	return nil, err

}

func (c *Connection) DB() *ent.Client {
	return c.db
}

func open(connector SQLConnector, driver string, connectionURL string) (*ent.Client, error) {
	client, err := connector.Connect(driver, connectionURL)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (c *Connection) Migrate() error {
	if err := c.db.Schema.Create(context.Background()); err != nil {
		return err
	}
	return nil
}
