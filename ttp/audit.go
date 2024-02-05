package ttp

import (
	"errors"
	"fmt"

	"github.com/akakou/metact"
	"github.com/akakou/ra_webs/ttp/ent/tainfo"
)

var ISSUER_NAME = "Let's Encrypt"

type Auditor struct {
	db *auditDB
	ct *metact.MetaCT
}

func NewAuditor(db *auditDB, ct *metact.MetaCT) (*Auditor, error) {
	return &Auditor{
		db: db,
		ct: ct,
	}, nil
}

func (auditor *Auditor) AuditOne(cert *metact.Certificate) error {
	domain, violatingDomains, err := validateDomains(cert.Domains)

	if err != nil || cert.IssuerName != ISSUER_NAME {
		revokeAllDomain(auditor.db, violatingDomains)
		return fmt.Errorf("domain violation: %w", err)
	}

	taInfo, err := auditor.db.client.TAInfo.
		Query().
		Where(tainfo.DomainEQ(domain)).
		WithCtLog().
		WithTaCode().
		Only(*auditor.db.ctx)

	if err != nil {
		return fmt.Errorf("failed to get ta info: %w", err)
	}

	ctLog := taInfo.Edges.CtLog
	if !ctLog.IsValid {
		return errors.New("ct log is not valid")
	}

	taCode := taInfo.Edges.TaCode[len(taInfo.Edges.TaCode)-1]

	if validateAttestation(cert, taCode.ProductID) != nil {
		ctLog.IsValid = false
		ctLog.Update().Save(*auditor.db.ctx)
		return fmt.Errorf("failed to check ct logs: %w", err)
	}

	ctLog.LatestCtID = cert.Id
	ctLog.Update().Save(*auditor.db.ctx)

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
