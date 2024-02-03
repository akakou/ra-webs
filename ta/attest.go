package ta

import (
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/akakou/ra_webs/core"
	"github.com/edgelesssys/ego/enclave"
)

type privateClaims struct {
	Data            string `json:"x-ms-sgx-ehd"`
	SecurityVersion uint   `json:"x-ms-sgx-svn"`
	Debug           bool   `json:"x-ms-sgx-is-debuggable"`
	UniqueID        string `json:"x-ms-sgx-mrenclave"`
	SignerID        string `json:"x-ms-sgx-mrsigner"`
	ProductID       uint   `json:"x-ms-sgx-product-id"`
}

func attestateByAzure(publicKey *rsa.PublicKey) (string, string, error) {
	rawPubKey := x509.MarshalPKCS1PublicKey(publicKey)

	publicKeyHashBytes := sha256.Sum256(rawPubKey)
	publicKeyHash := hex.EncodeToString(publicKeyHashBytes[:])

	token, err := enclave.CreateAzureAttestationToken([]byte(publicKeyHash), core.ATTEST_PROVIDER_URL)
	if err != nil {
		return "", "", fmt.Errorf("failed to create attestation token: %w", err)
	}

	report, err := enclave.GetSelfReport()
	if err != nil {
		return "", "", fmt.Errorf("failed to create attestation token: %w", err)
	}

	rawReport, err := json.Marshal(report)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal report: %w", err)
	}

	return token, string(rawReport), nil
}
