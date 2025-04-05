package monitor

import (
	"context"
	"fmt"

	"github.com/akakou/ra_webs/monitor/ent"
	_ "github.com/mattn/go-sqlite3"
)

type DBConfig struct {
	Type   string
	Config string
}

type DB struct {
	Client *ent.Client
	Ctx    *context.Context
}

func NewDB(dbConfig *DBConfig) (*DB, error) {
	client, err := ent.Open(dbConfig.Type, dbConfig.Config)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ERROR_OPEN_DB, err)
	}

	ctx := context.Background()

	if err := client.Schema.Create(ctx); err != nil {
		return nil, fmt.Errorf("%v: %w", ERROR_CREATE_SCHEMA, err)
	}

	return &DB{
		Client: client,
		Ctx:    &ctx,
	}, nil
}

func (db *DB) Close() {
	db.Client.Close()
}
