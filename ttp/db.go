package ttp

import (
	"context"
	"fmt"

	"github.com/akakou/ra_webs/ttp/ent"
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
		return nil, fmt.Errorf("failed opening connection to sqlite: %w", err)
	}

	ctx := context.Background()

	if err := client.Schema.Create(ctx); err != nil {
		return nil, fmt.Errorf("failed creating schema resources: %w", err)
	}

	return &DB{
		Client: client,
		Ctx:    &ctx,
	}, nil
}

func (db *DB) close() {
	db.Client.Close()
}
