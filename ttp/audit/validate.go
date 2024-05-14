package audit

import (
	"crypto/x509"
	"errors"
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
