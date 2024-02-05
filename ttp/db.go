package ttp

import (
	"context"
	"fmt"

	"github.com/akakou/ra_webs/ttp/ent"
	"github.com/akakou/ra_webs/ttp/ent/tainfo"
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
		taInfo, err := db.client.TAInfo.
			Query().
			Where(tainfo.DomainEQ(violatingDomain)).
			WithCtLog().
			WithTaCode().
			All(*db.ctx)

		if err != nil {
			continue
		}

		taInfo[0].Edges.CtLog.IsValid = false
		taInfo[0].Edges.CtLog.Update().Save(*db.ctx)
	}
}
