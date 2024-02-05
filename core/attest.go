package core

import (
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"

	"github.com/edgelesssys/ego/attestation"
	"github.com/edgelesssys/ego/enclave"
)

const SECURITY_VERSION = 1

func hashPublicKey(publicKey *rsa.PublicKey) []byte {
	rawPubKey := x509.MarshalPKCS1PublicKey(publicKey)

	publicKeyHashBytes := sha256.Sum256(rawPubKey)
	return publicKeyHashBytes[:]
}

func AttestByAzure(data []byte) (string, error) {
	// publicKeyHash := hashPublicKey(publicKey)
	token, err := enclave.CreateAzureAttestationToken(data, ATTEST_PROVIDER_URL)
	if err != nil {
		return "", fmt.Errorf("failed to create attestation token: %w", err)
	}

	return token, nil
}

func VerifyByAzure(token string, data []byte, productId uint16) (*attestation.Report, error) {
	report, err := attestation.VerifyAzureAttestationToken(token, ATTEST_PROVIDER_URL)
	if err != nil {
		return nil, fmt.Errorf("failed to verify attestation token: %w", err)
	}

	if binary.LittleEndian.Uint16(report.ProductID) != uint16(productId) {
		return nil, errors.New("token contains invalid product id")
	}

	if report.SecurityVersion < SECURITY_VERSION {
		return nil, errors.New("token contains invalid security version number")
	}

	if !reflect.DeepEqual(report.Data, data) {
		return nil, errors.New("token contains invalid public key hash")
	}

	return &report, nil
}
