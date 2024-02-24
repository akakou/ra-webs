package ttp

import (
	"errors"
	"fmt"
	"os"

	goutils "github.com/akakou/go-utils"
	golangutils "github.com/akakou/golang-utils"
	metact "github.com/akakou/meta-ct"
	"github.com/akakou/ra_webs/ttp/ent/taserver"
	simplecertify "github.com/akakou/simple-certify"
)

var ISSUER_NAME = "Let's Encrypt"

type Auditor struct {
	db         *DB
	ca         *simplecertify.Certifier
	ct         *metact.MetaCT
	adminToken string
}

func NewAuditor(db *DB, ca *simplecertify.Certifier, ct *metact.MetaCT, adminToken string) (*Auditor, error) {
	return &Auditor{
		db:         db,
		ca:         ca,
		ct:         ct,
		adminToken: adminToken,
	}, nil
}

func DefaultAuditor() (*Auditor, error) {
	dbType := golangutils.GetEnv("DB_TYPE", "sqlite3")
	dbConfig := golangutils.GetEnv("DB_CONFIG", "file:ent?mode=memory&cache=shared&_fk=1")
	fmt.Printf("We use %s as database type and %s as database config\n", dbType, dbConfig)

	metaAppId := os.Getenv("META_APP_ID")
	metaAccessToken := os.Getenv("META_ACCESS_TOKEN")

	adminToken, err := goutils.RandomHex(32)
	if err != nil {
		return nil, fmt.Errorf("failed to generate admin token: %w", err)
	}

	fmt.Printf("Admin token generated: %s\n", adminToken)

	caTempl := simplecertify.CATemplate()

	dbc := DBConfig{
		Type:   dbType,
		Config: dbConfig,
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

	return NewAuditor(db, ca, ct, adminToken)

}

func (auditor *Auditor) AuditOne(cert *metact.Certificate) error {
	domain, violatingDomains, err := validateDomains(cert.Domains)

	if err != nil || cert.IssuerName != ISSUER_NAME {
		revokeByDomain(auditor.db, violatingDomains)
		return fmt.Errorf("domain violation: %w", err)
	}

	taServ, err := auditor.db.Client.TAServer.
		Query().
		Where(taserver.DomainEQ(domain)).
		WithTa().
		Only(*auditor.db.Ctx)

	if err != nil {
		return fmt.Errorf("failed to get ta info: %w", err)
	}

	ta := taServ.Edges.Ta
	if !ta.IsValid {
		return errors.New("ct log is not valid")
	}

	taCode := ta.Edges.Code

	if validateAttestation(cert, taCode.UniqueID) != nil {
		ta.IsValid = false
		ta.Update().Save(*auditor.db.Ctx)
		return fmt.Errorf("failed to check ct logs: %w", err)
	}

	ta.LastCt = cert.Id
	ta.Update().Save(*auditor.db.Ctx)

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
