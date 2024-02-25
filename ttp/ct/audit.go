package ct

import (
	"errors"
	"fmt"

	metact "github.com/akakou/meta-ct"
	"github.com/akakou/ra_webs/ttp/core"
	"github.com/akakou/ra_webs/ttp/ent/taserver"
)

var ISSUER_NAME = "Let's Encrypt"

func AuditOne(ttp *core.TTP, cert *metact.Certificate) error {
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

	ta := taServ.Edges.Ta
	if !ta.IsValid {
		return errors.New("ct log is not valid")
	}

	taCode := ta.Edges.Code

	if validateAttestation(cert, taCode.UniqueID) != nil {
		ta.IsValid = false
		ta.Update().Save(*ttp.DB.Ctx)
		return fmt.Errorf("failed to check ct logs: %w", err)
	}

	_, err = ta.Edges.CtAudit.Update().SetLastCt(cert.Id).Save(*ttp.DB.Ctx)
	if err != nil {
		return fmt.Errorf("failed to update ct audit: %w", err)
	}

	_, err = ta.Update().Save(*ttp.DB.Ctx)
	if err != nil {
		return fmt.Errorf("failed to update ta: %w", err)
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
