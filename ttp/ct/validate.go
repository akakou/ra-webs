package ct

import (
	"crypto/rsa"
	"crypto/x509"
	"errors"
	"fmt"

	"github.com/akakou/ra_webs/core"
	"github.com/edgelesssys/ego/attestation"
)

func validateDomains(cert *x509.Certificate) (string, error) {
	domains := cert.DNSNames

	if len(domains) != 1 {
		return "", errors.New(ERROR_DOMAIN_INVALID_BY_NUM_DOMAIN)
	}

	domain := domains[0]

	if domain != cert.Subject.CommonName {
		return "", errors.New(ERROR_DOMAIN_INVALID_NOT_MATCH_COMMONNAME_AND_SAT)
	}

	return domain, nil
}

// for testability
func ValidateAttestation(token []byte, publicKey any) (*attestation.Report, error) {
	publicKeyBuf := x509.MarshalPKCS1PublicKey(publicKey.(*rsa.PublicKey))

	report, err := core.VerifyByAzure(string(token), publicKeyBuf)
	if err != nil {
		return nil, fmt.Errorf("%v: %v", ERROR_QUOTE_INVALID, err)
	}
	return report, nil
}
