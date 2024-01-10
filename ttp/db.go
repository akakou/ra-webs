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

func newTtpDB(dbConfig *DBConfig) (*ttpDB, error) {
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
		SetPublicKeyHash(taInfo.PublicKeyHash).
		SetAttestation(taInfo.Attestation)

	return entTaInfo
}

func (db *ttpDB) toCoreTaInfo(taInfo *ent.TAInfo) *core.TAInfo {
	return &core.TAInfo{
		Domain:        taInfo.Domain,
		PublicKeyHash: taInfo.PublicKeyHash,
		Attestation:   taInfo.Attestation,
	}
}

func (db *ttpDB) selectTaInfoByDomain(domain string) (*ent.TAInfo, error) {
	taInfo, err := db.client.TAInfo.
		Query().
		Where(tainfo.DomainEQ(domain)).
		WithCtLog().
		Only(*db.ctx)

	if err != nil {
		return nil, fmt.Errorf("failed querying ta info: %w", err)
	}

	return taInfo, nil
}

func (db *ttpDB) toEntCTLogAudit(ctLog *CTLogAudit, taInfo *ent.TAInfo) (*ent.CTLogAuditCreate, error) {
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
	taInfo, err := db.selectTaInfoByDomain(domain)
	if err != nil {
		return nil, err
	}

	ctLog := taInfo.Edges.CtLog
	ctLog.Edges.TaInfo = taInfo

	return ctLog, err
}
