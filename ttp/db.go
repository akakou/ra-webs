package ttp

import (
	"context"
	"fmt"

	"github.com/akakou/ra_webs/core"
	"github.com/akakou/ra_webs/ttp/ent"
	"github.com/akakou/ra_webs/ttp/ent/tainfo"
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

func newttpDB(dbConfig *DBConfig) (*ttpDB, error) {
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

func (db *ttpDB) findByDomain(domain string) (*core.TAInfo, error) {
	taInfoColumn, err := db.client.TAInfo.
		Query().
		Where(tainfo.DomainEQ(domain)).
		Only(*db.ctx)

	if err != nil {
		return nil, fmt.Errorf("failed querying ta info: %w", err)
	}

	taInfo := &core.TAInfo{
		Domain:      taInfoColumn.Domain,
		PublicKey:   taInfoColumn.PublicKey,
		Attestation: taInfoColumn.Attestation,
	}

	return taInfo, nil
}

func (db *ttpDB) store(taInfo *core.TAInfo) error {
	_, err := db.client.TAInfo.
		Create().
		SetDomain(taInfo.Domain).
		SetPublicKey(taInfo.PublicKey).
		SetAttestation(taInfo.Attestation).
		Save(*db.ctx)

	if err != nil {
		return fmt.Errorf("failed creating user: %w", err)
	}

	return nil
}
