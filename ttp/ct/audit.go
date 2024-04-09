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

func AuditOne(ttp *core.TTP, c *metact.MetaCert) error {
	cert, err := c.Certificate()
	if err != nil {
		err = fmt.Errorf("failed to get certificate: %w", err)
		panic(err)
	}

	domain, err := validateDomains(cert.DNSNames)
	if err != nil {
		revokeTAByDomains(cert.DNSNames, ttp.DB)
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
	lastTA, err := taServ.QueryTa().
		Order(ent.Desc(ta.FieldID)).
		First(*ttp.DB.Ctx)

	if err != nil && !ent.IsNotFound(err) {
		return fmt.Errorf("failed to check ct logs: %w", err)
	}

	if lastTA != nil && !lastTA.IsValid {
		return fmt.Errorf("last TA is invalid: %w", err)
	}

	_ta := ttp.DB.Client.TA.Create().
		SetServer(taServ).
		SetPublicKey([]byte{}).
		SetIsValid(false).
		SaveX(*ttp.DB.Ctx)

	report, err := validateAttestation(cert)
	if err != nil {
		return fmt.Errorf("failed to get validate quote: %w", err)
	}

	// check if the ta code has been registered
	taCode, err := ttp.DB.Client.TACode.
		Query().
		Where(tacode.UniqueID(report.UniqueID)).
		Only(*ttp.DB.Ctx)

	if err != nil {
		return fmt.Errorf("failed to get ta code: %w", err)
	}

	_ta.Update().
		SetCode(taCode).
		SetPublicKey(_ta.PublicKey).
		SetIsValid(true).
		SaveX(*ttp.DB.Ctx)

	return nil
}

func AuditAll(ttp *core.TTP, cert []metact.MetaCert) error {
	for _, c := range cert {
		err := AuditOne(ttp, &c)
		if err != nil {
			return fmt.Errorf("failed to audit: %w", err)
		}
	}
	return nil
}
