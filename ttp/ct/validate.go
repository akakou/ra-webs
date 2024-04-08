package ct

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	metact "github.com/akakou/meta-ct"
	"github.com/akakou/ra_webs/core"
	"github.com/edgelesssys/ego/attestation"
	"golang.org/x/exp/slices"
)

func validateDomains(domains []string) (string, error) {
	if len(domains) != 1 {
		return "", fmt.Errorf("number of domain must be 1")
	}

	domain := domains[0]

	chars := strings.Split(domain, "")
	if slices.Contains(chars, "*") {
		return "", fmt.Errorf("wildcard domain is not allowed")
	}

	return domain, nil
}

func validatePublicKey(cert *metact.Certificate) error {
	if len(cert.PublicKeyValues) == 1 {
		return nil
	} else {
		return errors.New("multiple public key not supported")
	}
}

// for debuggability
var validateAttestation = _validateAttestation

func _validateAttestation(cert *metact.Certificate) (*attestation.Report, error) {
	token, err := findCertExtensions(cert.Extensions, "core.X509_EXTENSION_LABEL")
	if err != nil {
		return nil, fmt.Errorf("extension not found: %v", err)
	}

	hashedPublicKey, err := hex.DecodeString(cert.PublicKeyHashSha256)
	if err != nil {
		return nil, fmt.Errorf("failed to decode public key hash: %v", err)
	}

	report, err := core.VerifyByAzure(token, hashedPublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to verify attestation: %v", err)
	}
	return report, nil
}
