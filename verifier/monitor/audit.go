package monitor

import (
	"crypto/rsa"
	"fmt"

	"github.com/akakou/ra_webs/verifier/core"
	"github.com/akakou/ra_webs/verifier/ent"
	"github.com/akakou/ra_webs/verifier/ent/taserver"
	"github.com/google/certificate-transparency-go/x509"
)

func Monitor(verifier *core.Verifier, cert *x509.Certificate) error {
	domain, err := validateDomains(cert)
	if err != nil {
		revokeByDomains(cert.DNSNames, verifier)
		return fmt.Errorf("%s: %w", ERROR_DOMAIN_INVALID, err)
	}

	unmarshaledPublicKey, isRSA := cert.PublicKey.(*rsa.PublicKey)

	if !isRSA {
		revokeByDomain(domain, lastValidID(domain, verifier.DB), verifier)
		return fmt.Errorf("%s", ERROR_PUBLIC_KEY_NOT_RSA)
	}

	publicKey := x509.MarshalPKCS1PublicKey(unmarshaledPublicKey)

	lastID := lastValidID(domain, verifier.DB)

	serv, err := verifier.DB.Client.TAServer.
		Query().
		Where(taserver.DomainEQ(domain)).
		Where(taserver.HasActivated(false)).
		Where(taserver.PublicKey(publicKey)).
		Where(taserver.IDGT(lastID - 1)).
		Order(ent.Desc(taserver.FieldID)).
		First(*verifier.DB.Ctx)

	if err != nil {
		revokeByDomain(domain, lastID, verifier)
		return fmt.Errorf("%v: %v", ERROR_CERTIFICATE_NOT_FOUND, err)
	}

	if !serv.HasActivated {
		serv.Update().SetHasActivated(true).Save(*verifier.DB.Ctx)
	}

	return nil
}
