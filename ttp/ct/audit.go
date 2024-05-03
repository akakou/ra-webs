package ct

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"

	"github.com/akakou/ra_webs/ttp/core"
	"github.com/akakou/ra_webs/ttp/ent/taserver"
)

func AuditOne(ttp *core.TTP, cert *x509.Certificate) error {
	domain, err := validateDomains(cert)
	if err != nil {
		logViolationsByDomains(cert.DNSNames, ttp.DB)
		return fmt.Errorf("%s: %w", ERROR_DOMAIN_INVALID, err)
	}

	// get the last ta from ta server
	serv, err := ttp.DB.Client.TAServer.
		Query().
		Where(taserver.DomainEQ(domain)).
		Only(*ttp.DB.Ctx)

	if err != nil {
		return fmt.Errorf("%s: %w", ERROR_SELECT_SERVER, err)
	}

	isValid := cert.PublicKey.(*rsa.PublicKey).Equal(serv.PublicKey)

	if !isValid {
		logViolationByDomain(domain, ttp.DB)
	}

	if !serv.HasActivated {
		serv.Update().SetHasActivated(true).Save(*ttp.DB.Ctx)
	}

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
