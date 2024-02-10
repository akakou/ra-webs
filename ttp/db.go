package ttp

import (
	"context"
	"fmt"

	"github.com/akakou/ra_webs/ttp/ent"
	"github.com/akakou/ra_webs/ttp/ent/ta"
	_ "github.com/mattn/go-sqlite3"
)

type DBConfig struct {
	Type   string
	Config string
}

type auditDB struct {
	client *ent.Client
	ctx    *context.Context
}

func newAuditDB(dbConfig *DBConfig) (*auditDB, error) {
	client, err := ent.Open(dbConfig.Type, dbConfig.Config)
	if err != nil {
		return nil, fmt.Errorf("failed opening connection to sqlite: %w", err)
	}

	ctx := context.Background()

	if err := client.Schema.Create(ctx); err != nil {
		return nil, fmt.Errorf("failed creating schema resources: %w", err)
	}

	return &auditDB{
		client: client,
		ctx:    &ctx,
	}, nil
}

func (db *auditDB) close() {
	db.client.Close()
}

func revokeAllDomain(db *auditDB, domains []string) {
	for _, violatingDomain := range domains {
		taInfos, err := db.client.TA.
			Query().
			Where(ta.DomainEQ(violatingDomain)).
			WithAuditLog().
			WithCode().
			All(*db.ctx)

		if err != nil {
			continue
		}

		for _, taInfo := range taInfos {
			taInfo.Edges.AuditLog.IsValid = false
			taInfo.Edges.AuditLog.Update().Save(*db.ctx)

		}
	}
}
