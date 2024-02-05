package ttp

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/akakou/metact"
	"github.com/akakou/ra_webs/core"
)

func validateDomains(domains []string) (string, []string, error) {
	if len(domains) == 0 {
		return domains[0], []string{}, nil
	}

	violatingDomains := []string{}

	for _, domain := range domains {
		domain := splitDomainLast(domain)
		violatingDomains = append(violatingDomains, domain)
	}

	return "", violatingDomains, errors.New("domain violation")

}

func validateAttestation(cert *metact.Certificate, productId uint16) error {
	token, err := findCertExtensions(cert.Extensions, "core.X509_EXTENSION_LABEL")
	if err != nil {
		return fmt.Errorf("extension not found: %v", err)
	}

	hashedPublicKey, err := hex.DecodeString(cert.PublicKeyHashSha256)
	if err != nil {
		return fmt.Errorf("failed to decode public key hash: %v", err)
	}

	_, err = core.VerifyByAzure(token, hashedPublicKey, productId)
	if err != nil {
		return fmt.Errorf("failed to verify attestation: %v", err)
	}
	return nil
}
