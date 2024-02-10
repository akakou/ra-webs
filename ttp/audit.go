package ttp

import (
	"errors"
	"fmt"

	"github.com/akakou/metact"
	"github.com/akakou/ra_webs/ttp/ent/ta"
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

	taInfo, err := auditor.db.client.TA.
		Query().
		Where(ta.DomainEQ(domain)).
		WithAuditLog().
		WithCode().
		Only(*auditor.db.ctx)

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
		auditLog.Update().Save(*auditor.db.ctx)
		return fmt.Errorf("failed to check ct logs: %w", err)
	}

	auditLog.LatestCtID = cert.Id
	auditLog.Update().Save(*auditor.db.ctx)

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
