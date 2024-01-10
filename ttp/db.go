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

func (db *ttpDB) toEntTaInfo(taInfo *core.TAInfo) *ent.TAInfoCreate {
	entTaInfo := db.client.TAInfo.
		Create().
		SetDomain(taInfo.Domain).
		SetPublicKey(taInfo.PublicKey).
		SetAttestation(taInfo.Attestation)

	return entTaInfo
}

func (db *ttpDB) toCoreTaInfo(taInfo *ent.TAInfo) *core.TAInfo {
	return &core.TAInfo{
		Domain:      taInfo.Domain,
		PublicKey:   taInfo.PublicKey,
		Attestation: taInfo.Attestation,
	}
}

func (db *ttpDB) selectTaInfoByDomain(domain string) (*ent.TAInfo, error) {
	return db._selectTaInfoByDomain(domain, false)
}

func (db *ttpDB) _selectTaInfoByDomain(domain string, withCTLogs bool) (*ent.TAInfo, error) {
	query := db.client.TAInfo.
		Query().
		Where(tainfo.DomainEQ(domain))

	if withCTLogs {
		query = query.WithCtLog()
	}

	taInfo, err := query.Only(*db.ctx)

	if err != nil {
		return nil, fmt.Errorf("failed querying ta info: %w", err)
	}

	return taInfo, nil
}

func (db *ttpDB) toEntCTLogAuditWithRelation(ctLog *CTLogAudit) (*ent.CTLogAuditCreate, error) {
	taInfo, err := db.selectTaInfoByDomain(ctLog.TADomain)
	if err != nil {
		return nil, err
	}

	entCTLog := db.client.CTLogAudit.
		Create().
		SetLatestCtID(ctLog.LatestCTId).
		SetIsValid(ctLog.IsValid).
		SetTaInfo(taInfo)

	return entCTLog, nil
}

func (db *ttpDB) toCoreCTLogAudit(ctLogAudit *ent.CTLogAudit) *CTLogAudit {
	return &CTLogAudit{
		TADomain:   ctLogAudit.Edges.TaInfo.Domain,
		IsValid:    ctLogAudit.IsValid,
		LatestCTId: ctLogAudit.LatestCtID,
	}
}

func (db *ttpDB) selectCTLogByDomain(domain string) (*ent.CTLogAudit, error) {
	taInfo, err := db._selectTaInfoByDomain(domain, true)
	if err != nil {
		return nil, err
	}

	ctLog := taInfo.Edges.CtLog
	ctLog.Edges.TaInfo = taInfo

	return ctLog, err
}
