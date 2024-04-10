package ct

import (
	"crypto/x509"
	"fmt"

	rawebcore "github.com/akakou/ra_webs/core"
	"github.com/akakou/ra_webs/ttp/core"
	"github.com/akakou/ra_webs/ttp/ent"
	"github.com/akakou/ra_webs/ttp/ent/ta"
	"github.com/akakou/ra_webs/ttp/ent/tacode"
	"github.com/akakou/ra_webs/ttp/ent/taserver"
)

func AuditOne(ttp *core.TTP, cert *x509.Certificate) error {
	domain, err := validateDomains(cert.DNSNames)
	if err != nil {
		revokeTAByDomains(cert.DNSNames, ttp.DB)
		return fmt.Errorf("%s: %w", ERROR_DOMAIN_INVALID, err)
	}

	// get the last ta from ta server
	taServ, err := ttp.DB.Client.TAServer.
		Query().
		Where(taserver.DomainEQ(domain)).
		Only(*ttp.DB.Ctx)

	if err != nil {
		return fmt.Errorf("%s: %w", ERROR_SELECT_SERVER, err)
	}

	// fetch and check the status of last ta
	lastTA, err := taServ.QueryTa().
		Order(ent.Desc(ta.FieldID)).
		First(*ttp.DB.Ctx)

	if err != nil && !ent.IsNotFound(err) {
		return fmt.Errorf("%s: %w", ERROR_SELECT_LAST_LOG, err)
	}

	if lastTA != nil && !lastTA.IsValid {
		return fmt.Errorf("%s: %w", ERROR_LAST_TA_INVALID, err)
	}

	_ta := ttp.DB.Client.TA.Create().
		SetServer(taServ).
		SetPublicKey([]byte{}).
		SetQuote([]byte{}).
		SetIsValid(false).
		SaveX(*ttp.DB.Ctx)

	token, err := findCertExtensions(rawebcore.X509_EXTENSION_LABEL, cert)
	if err != nil {
		return fmt.Errorf("%v: %v", ERROR_EXTENSION_NOT_FOUND, err)
	}

	report, err := ValidateAttestation(token, cert.PublicKey)
	if err != nil {
		return fmt.Errorf("%s: %w", ERROR_QUOTE_INVALID, err)
	}

	// check if the ta code has been registered
	taCode, err := ttp.DB.Client.TACode.
		Query().
		Where(tacode.UniqueID(report.UniqueID)).
		Only(*ttp.DB.Ctx)

	if err != nil {
		return fmt.Errorf("%s %w", ERROR_SELECT_TA_CODE, err)
	}

	_ta.Update().
		SetCode(taCode).
		SetPublicKey(_ta.PublicKey).
		SetQuote(token).
		SetIsValid(true).
		SaveX(*ttp.DB.Ctx)

	return nil
}

func AuditAll(ttp *core.TTP, cert []x509.Certificate) error {
	for _, c := range cert {
		err := AuditOne(ttp, &c)
		if err != nil {
			return fmt.Errorf("failed to audit: %w", err)
		}
	}
	return nil
}
