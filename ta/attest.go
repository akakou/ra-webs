package ta

import (
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"fmt"

	"github.com/akakou/ra_webs/core"
	"github.com/edgelesssys/ego/enclave"
)

func AttestateByAzure(publicKey *rsa.PublicKey) (string, error) {
	rawPubKey := x509.MarshalPKCS1PublicKey(publicKey)

	publicKeyHashBytes := sha256.Sum256(rawPubKey)
	publicKeyHash := hex.EncodeToString(publicKeyHashBytes[:])

	token, err := enclave.CreateAzureAttestationToken([]byte(publicKeyHash), core.ATTEST_PROVIDER_URL)
	if err != nil {
		return "", fmt.Errorf("failed to create attestation token: %w", err)
	}

	return token, nil
}
