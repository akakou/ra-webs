package sign

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignature(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)

	log := LogPlain{
		Repository: "test",
		CommitId:   "test",
		Evidence:   "test",
	}

	signature, err := Sign(&log, privateKey)
	assert.NoError(t, err, "Sigining failed")

	err = Verify(signature, &log, &privateKey.PublicKey)
	assert.NoError(t, err, "Signature verification failed")
}
