package ttp

import (
	"context"
	"fmt"

	"github.com/akakou/ra_webs/core"
	"github.com/akakou/ra_webs/ttp/ent"
	_ "github.com/mattn/go-sqlite3"
)

type DBConfig struct {
	Type   string
	Config string
}

type taInfoDB struct {
	client *ent.Client
	ctx    *context.Context
}

func newtTAInfoDB(dbConfig *DBConfig) (*taInfoDB, error) {
	client, err := ent.Open(dbConfig.Type, dbConfig.Config)
	if err != nil {
		return nil, fmt.Errorf("failed opening connection to sqlite: %w", err)
	}

	ctx := context.Background()

	if err := client.Schema.Create(ctx); err != nil {
		return nil, fmt.Errorf("failed creating schema resources: %w", err)
	}

	return &taInfoDB{
		client: client,
		ctx:    &ctx,
	}, nil
}

func (db *taInfoDB) store(req *core.ProvisioningRequest) error {
	_, err := db.client.TAInfo.
		Create().
		SetDomain(req.Domain).
		SetPublicKey(req.PublicKey).
		SetAttestation(req.Attestation).
		Save(*db.ctx)

	if err != nil {
		return fmt.Errorf("failed creating user: %w", err)
	}

	return nil
}
