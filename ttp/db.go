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

type ttpDB struct {
	client *ent.Client
	ctx    *context.Context
}

func newTTPDB(dbConfig *DBConfig) (*ttpDB, error) {
	client, err := ent.Open(dbConfig.Type, dbConfig.Config)
	if err != nil {
		return nil, fmt.Errorf("failed opening connection to sqlite: %w", err)
	}

	ctx := context.Background()

	if err := client.Schema.Create(ctx); err != nil {
		return nil, fmt.Errorf("failed creating schema resources: %w", err)
	}

	return &ttpDB{
		client: client,
		ctx:    &ctx,
	}, nil
}

func (db *ttpDB) close() {
	db.client.Close()
}
