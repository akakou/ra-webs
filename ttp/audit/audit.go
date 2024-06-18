package audit

import (
	"crypto/rsa"
	"fmt"

	"github.com/akakou/ra_webs/ttp/core"
	"github.com/akakou/ra_webs/ttp/ent"
	"github.com/akakou/ra_webs/ttp/ent/taserver"
	"github.com/google/certificate-transparency-go/x509"
)

func Audit(ttp *core.TTP, cert *x509.Certificate) error {
	domain, err := validateDomains(cert)
	if err != nil {
		revokeByDomains(cert.DNSNames, ttp.DB)
		return fmt.Errorf("%s: %w", ERROR_DOMAIN_INVALID, err)
	}

	unmarshaledPublicKey, isRSA := cert.PublicKey.(*rsa.PublicKey)

	if !isRSA {
		revokeByDomain(domain, lastValidID(domain, ttp.DB), ttp.DB)
		return fmt.Errorf("%s", ERROR_PUBLIC_KEY_NOT_RSA)
	}

	publicKey := x509.MarshalPKCS1PublicKey(unmarshaledPublicKey)

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
