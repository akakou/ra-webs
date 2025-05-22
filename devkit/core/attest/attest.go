package attest

import (
	"errors"
	"fmt"

	"github.com/edgelesssys/ego/attestation"
	"github.com/edgelesssys/ego/enclave"
)

const SECURITY_VERSION = 1
const ATTEST_PROVIDER_URL = "https://shareduks.uks.attest.azure.net"

var Attest = attestByAzure
var Verify = verifyByAzure

func attestByAzure(data []byte) (string, error) {
	// publicKeyHash := hashPublicKey(publicKey)
	evidence, err := enclave.CreateAzureAttestationToken(data, ATTEST_PROVIDER_URL)
	if err != nil {
		return "", fmt.Errorf("%s: %w", ERROR_CREATE_ATTESTATION, err)
	}

	fmt.Printf("evidence is created! : %v\n", evidence)
	return evidence, nil
}

func verifyByAzure(quote string) (*attestation.Report, error) {
	report, err := attestation.VerifyAzureAttestationToken(quote, ATTEST_PROVIDER_URL)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ERROR_VERIFY_ATTESTATION, err)
	}

	if report.SecurityVersion < SECURITY_VERSION {
		return nil, errors.New(ERROR_INVALID_SECURITY_VERSION_IN_ATTESTATION)
	}

	return &report, nil
}
