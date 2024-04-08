package ct

import (
	"fmt"

	metact "github.com/akakou/meta-ct"
	"github.com/akakou/ra_webs/ttp/core"
	"github.com/akakou/ra_webs/ttp/ent"
	"github.com/akakou/ra_webs/ttp/ent/ta"
	"github.com/akakou/ra_webs/ttp/ent/tacode"
	"github.com/akakou/ra_webs/ttp/ent/taserver"
)

func AuditOne(ttp *core.TTP, cert *metact.Certificate) error {
	_ta := ttp.DB.Client.TA.Create()

	// validate the domain
	domain, err := validateDomains(cert.Domains)

	if err != nil {
		revokeTAByDomains(ttp.DB, cert.Domains)
		return fmt.Errorf("domain violation: %w", err)
	}

	// get the last ta from ta server
	taServ, err := ttp.DB.Client.TAServer.
		Query().
		Where(taserver.DomainEQ(domain)).
		Only(*ttp.DB.Ctx)

	if err != nil {
		return fmt.Errorf("failed to get ta info: %w", err)
	}

	// fetch and check the status of last ta
	lastTA, _ := taServ.QueryTa().Order(ent.Desc(ta.FieldID)).First(*ttp.DB.Ctx)

	if lastTA != nil && !lastTA.IsValid {
		return fmt.Errorf("server is not valid")
	}

	// validate quotes
	report, err := validateAttestation(cert)

	if err != nil {
		revokeTA(lastTA, ttp.DB)
		return fmt.Errorf("failed to check ct logs: %w", err)
	}

	// check if the ta code has been registered
	taCode, err := ttp.DB.Client.TACode.
		Query().
		Where(tacode.UniqueID(report.UniqueID)).
		Only(*ttp.DB.Ctx)

	if err != nil {
		revokeTA(lastTA, ttp.DB)
		return fmt.Errorf("failed to get ta code: %w", err)
	}

	// save ta
	_, err = _ta.SetCode(taCode).
		SetServer(taServ).
		SetPublicKey(report.Data).
		SetIsValid(true).
		Save(*ttp.DB.Ctx)

	if err != nil {
		revokeTA(lastTA, ttp.DB)
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
