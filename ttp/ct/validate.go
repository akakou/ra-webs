package ct

import (
	"crypto/rsa"
	"crypto/x509"
	"errors"
	"fmt"
	"strings"

	"github.com/akakou/ra_webs/core"
	"github.com/edgelesssys/ego/attestation"
	"golang.org/x/exp/slices"
)

func validateDomains(domains []string) (string, error) {
	if len(domains) != 1 {
		return "", errors.New(ERROR_DOMAIN_INVALID_BY_NUM_DOMAIN)
	}

	domain := domains[0]

	chars := strings.Split(domain, "")
	if slices.Contains(chars, "*") {
		return "", errors.New(ERROR_DOMAIN_INVALID_BY_WILDCARD)
	}

	return domain, nil
}

// for debuggability
var validateAttestation = _validateAttestation

func _validateAttestation(cert *x509.Certificate) (*attestation.Report, error) {
	token, err := findCertExtensions(core.X509_EXTENSION_LABEL, cert)
	if err != nil {
		return nil, fmt.Errorf("%v: %v", ERROR_EXTENSION_NOT_FOUND, err)
	}

	publicKey := cert.PublicKey
	publicKeyBuf := x509.MarshalPKCS1PublicKey(publicKey.(*rsa.PublicKey))

	report, err := core.VerifyByAzure(string(token), publicKeyBuf)
	if err != nil {
		return nil, fmt.Errorf("%v: %v", ERROR_QUOTE_INVALID, err)
	}
	return report, nil
}
