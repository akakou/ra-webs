package ct

import (
	"crypto/rsa"
	"crypto/x509"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/akakou/ra_webs/core"
	"github.com/edgelesssys/ego/attestation"
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

// for testability
var ValidateAttestation = _validateAttestation

func _validateAttestation(token []byte, publicKey any) (*attestation.Report, error) {
	publicKeyBuf := x509.MarshalPKCS1PublicKey(publicKey.(*rsa.PublicKey))

	report, err := core.VerifyByAzure(string(token), publicKeyBuf)
	if err != nil {
		return nil, fmt.Errorf("%v: %v", ERROR_QUOTE_INVALID, err)
	}
	return report, nil
}
