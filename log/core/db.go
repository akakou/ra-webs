package core

import (
	"context"

	"github.com/akakou/ra-webs/log/ent"
	"github.com/cockroachdb/errors"
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
		return nil, errors.Wrap(err, ErrFailedToOpenDB.Error())
	}

	ctx := context.Background()

	if err := client.Schema.Create(ctx); err != nil {
		return nil, errors.Wrap(err, ErrCreateSchema.Error())
	}

	return &DB{
		Client: client,
		Ctx:    &ctx,
	}, nil
}

func (db *DB) Close() {
	db.Client.Close()
}
