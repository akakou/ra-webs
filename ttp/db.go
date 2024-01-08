package main

import (
	"context"
	"fmt"

	"github.com/akakou/ra_webs/core"
	"github.com/akakou/ra_webs/ttp/ent"
	_ "github.com/mattn/go-sqlite3"
)

type taInfoDB struct {
	client *ent.Client
	ctx    *context.Context
}

func newtTAInfoDB(dbType, dbConfig string) (*taInfoDB, error) {
	client, err := ent.Open(dbType, dbConfig)
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
