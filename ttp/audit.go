package ttp

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/akakou/metact"
	"github.com/akakou/ra_webs/ttp/ent/ta"
	simplecertify "github.com/akakou/simple-certify"
)

var ISSUER_NAME = "Let's Encrypt"

type Auditor struct {
	db *DB
	ca *simplecertify.Certifier
	ct *metact.MetaCT
}

func NewAuditor(db *DB, ca *simplecertify.Certifier, ct *metact.MetaCT) (*Auditor, error) {
	return &Auditor{
		db: db,
		ca: ca,
		ct: ct,
	}, nil
}

func DefaultAuditor() (*Auditor, error) {
	dbType := flag.String("db_type", "sqlite3", "database type")
	dbConfig := flag.String("db_config", "file:ent?mode=memory&cache=shared&_fk=1", "database config")

	metaAppId := os.Getenv("META_APP_ID")
	metaAccessToken := os.Getenv("META_ACCESS_TOKEN")

	fmt.Printf("We use %s as database type and %s as database config\n", *dbType, *dbConfig)

	caTempl := simplecertify.CATemplate()

	dbc := DBConfig{
		Type:   *dbType,
		Config: *dbConfig,
	}

	db, err := NewDB(&dbc)
	if err != nil {
		return nil, fmt.Errorf("failed to init db: %w", err)
	}

	ct := metact.NewCT(metaAppId, metaAccessToken)

	ca, err := simplecertify.LoadOrInit(&caTempl, &caTempl)

	if err != nil {
		return nil, fmt.Errorf("failed to init ca: %w", err)

	}

	return NewAuditor(db, ca, ct)

}

func (auditor *Auditor) AuditOne(cert *metact.Certificate) error {
	domain, violatingDomains, err := validateDomains(cert.Domains)

	if err != nil || cert.IssuerName != ISSUER_NAME {
		revokeAllDomain(auditor.db, violatingDomains)
		return fmt.Errorf("domain violation: %w", err)
	}

	taInfo, err := auditor.db.Client.TA.
		Query().
		Where(ta.DomainEQ(domain)).
		WithAuditLog().
		WithCode().
		Only(*auditor.db.Ctx)

	if err != nil {
		return fmt.Errorf("failed to get ta info: %w", err)
	}

	auditLog := taInfo.Edges.AuditLog
	if !auditLog.IsValid {
		return errors.New("ct log is not valid")
	}

	taCode := taInfo.Edges.Code[len(taInfo.Edges.Code)-1]

	if validateAttestation(cert, taCode.UniqueID) != nil {
		auditLog.IsValid = false
		auditLog.Update().Save(*auditor.db.Ctx)
		return fmt.Errorf("failed to check ct logs: %w", err)
	}

	auditLog.LatestCtID = cert.Id
	auditLog.Update().Save(*auditor.db.Ctx)

	return nil
}

func (auditor *Auditor) AuditAll(cert []metact.Certificate) error {
	for _, c := range cert {
		err := auditor.AuditOne(&c)
		if err != nil {
			return fmt.Errorf("failed to audit: %w", err)
		}
	}
	return nil
}
