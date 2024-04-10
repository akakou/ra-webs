package core

import (
	"crypto/sha256"
	"encoding/hex"
)

func DomainToken(serviceToken, nonce string) string {
	hashSource := []byte{}
	hashSource = append(hashSource, serviceToken...)
	hashSource = append(hashSource, nonce...)

	hash := sha256.Sum256(hashSource)
	expected := hex.EncodeToString(hash[:])

	return expected
}
