package ct

import (
	"fmt"

	metact "github.com/akakou/meta-ct"
	"github.com/akakou/ra_webs/ttp/core"
	"github.com/akakou/ra_webs/ttp/ent/tacode"
	"github.com/akakou/ra_webs/ttp/ent/taserver"
)

var ISSUER_NAME = "Let's Encrypt"

func AuditOne(ttp *core.TTP, cert *metact.Certificate) error {
	ta := ttp.DB.Client.TA.Create()
	domain, violatingDomains, err := validateDomains(cert.Domains)

	if err != nil || cert.IssuerName != ISSUER_NAME {
		revokeByDomain(ttp.DB, violatingDomains)
		return fmt.Errorf("domain violation: %w", err)
	}

	taServ, err := ttp.DB.Client.TAServer.
		Query().
		Where(taserver.DomainEQ(domain)).
		WithTa().
		Only(*ttp.DB.Ctx)

	if err != nil {
		return fmt.Errorf("failed to get ta info: %w", err)
	}

	report, err := validateAttestation(cert)

	if err != nil {
		ta.SetIsValid(false)
		ta.Save(*ttp.DB.Ctx)
		return fmt.Errorf("failed to check ct logs: %w", err)
	}

	taCode, err := ttp.DB.Client.TACode.
		Query().
		Where(tacode.UniqueID(report.UniqueID)).
		Only(*ttp.DB.Ctx)

	if err != nil {
		ta.SetIsValid(false)
		return fmt.Errorf("failed to get ta code: %w", err)
	}

	ta.SetCode(taCode).
		SetServer(taServ).
		SetPublicKey(report.Data).
		SetIsValid(true).
		Save(*ttp.DB.Ctx)

	if err != nil {
		return fmt.Errorf("failed to create ta: %w", err)
	}

	return nil
}

func AuditAll(ttp *core.TTP, cert []metact.Certificate) error {
	for _, c := range cert {
		err := AuditOne(ttp, &c)
		if err != nil {
			return fmt.Errorf("failed to audit: %w", err)
		}
	}
	return nil
}
