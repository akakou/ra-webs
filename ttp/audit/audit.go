package audit

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"

	"github.com/akakou/ra_webs/ttp/core"
	"github.com/akakou/ra_webs/ttp/ent"
	"github.com/akakou/ra_webs/ttp/ent/taserver"
)

func AuditOne(ttp *core.TTP, cert *x509.Certificate) error {
	domain, err := validateDomains(cert)
	if err != nil {
		revokeByDomains(cert.DNSNames, ttp.DB)
		return fmt.Errorf("%s: %w", ERROR_DOMAIN_INVALID, err)
	}

	publicKey := x509.MarshalPKCS1PublicKey(cert.PublicKey.(*rsa.PublicKey))
	lastID := lastValidID(domain, ttp.DB)

	serv, err := ttp.DB.Client.TAServer.
		Query().
		Where(taserver.DomainEQ(domain)).
		Where(taserver.HasActivated(false)).
		Where(taserver.PublicKey(publicKey)).
		Where(taserver.IDGT(lastID - 1)).
		Order(ent.Desc(taserver.FieldID)).
		First(*ttp.DB.Ctx)

	if err != nil {
		revokeByDomain(domain, lastID, ttp.DB)
		return fmt.Errorf("%v: %v", ERROR_CERTIFICATE_NOT_FOUND, err)
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
			fmt.Printf("failed to audit: %v", err)
		}
	}
	return nil
}
