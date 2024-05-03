package core

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/edgelesssys/ego/attestation"
	"github.com/edgelesssys/ego/enclave"
)

const SECURITY_VERSION = 1

var AttestByAzure = attestByAzure
var VerifyByAzure = verifyByAzure

func AttestServer(publicKey []byte, token string) (string, error) {
	buf := append([]byte(token), publicKey...)
	quote, err := AttestByAzure(buf)
	return quote, err
}

func VerifyServer(quote string, publicKey []byte, token string) (*attestation.Report, error) {
	buf := append([]byte(token), publicKey...)
	report, err := verifyByAzure(quote, []byte(buf))
	return report, err
}

func attestByAzure(data []byte) (string, error) {
	if DEBUG {
		return "", nil
	}

	// publicKeyHash := hashPublicKey(publicKey)
	token, err := enclave.CreateAzureAttestationToken(data, ATTEST_PROVIDER_URL)
	if err != nil {
		return "", fmt.Errorf("%s: %w", ERROR_CREATE_ATTESTATION, err)
	}

	return token, nil
}

func verifyByAzure(quote string, data []byte) (*attestation.Report, error) {
	report, err := attestation.VerifyAzureAttestationToken(quote, ATTEST_PROVIDER_URL)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ERROR_VERIFY_ATTESTATION, err)
	}

	if report.SecurityVersion < SECURITY_VERSION {
		return nil, errors.New(ERROR_INVALID_SECURITY_VERSION_IN_ATTESTATION)
	}

	if !bytes.Equal(report.Data, data) {
		return nil, errors.New(ERROR_INVALID_REPORT_DATA_IN_ATTESTATION)
	}

	return &report, nil
}
