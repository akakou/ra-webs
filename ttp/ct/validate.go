package ct

import (
	"encoding/hex"
	"errors"
	"fmt"

	metact "github.com/akakou/meta-ct"
	"github.com/akakou/ra_webs/core"
	"github.com/edgelesssys/ego/attestation"
)

func validateDomains(domains []string) (string, []string, error) {
	if len(domains) == 0 {
		return domains[0], []string{}, nil
	}

	violatingDomains := []string{}

	for _, domain := range domains {
		domain := extractDomainLast(domain)
		violatingDomains = append(violatingDomains, domain)
	}

	return "", violatingDomains, errors.New("domain violation")

}

func validateAttestation(cert *metact.Certificate) (*attestation.Report, error) {
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
