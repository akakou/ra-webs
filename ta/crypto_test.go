package ta

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecureChannel(t *testing.T) {
	expected := []byte("hello")

	sc, err := newSecureChannel([]byte("01234567890123456789012345678901"))
	assert.NoError(t, err)

	cipher, err := sc.encrypt(expected)
	assert.NoError(t, err)

	actual, err := sc.decrypt(cipher)
	assert.NoError(t, err)

	assert.Equal(t, expected, actual)
}

func TestSecureLoader(t *testing.T) {
	expected := []byte("hello")

	ra := NewRA(nil)
	privKey, pubKey, err := ra.generateKeyPair()
	assert.NoError(t, err)

	cipher, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pubKey, expected, []byte{})
	assert.NoError(t, err)

	receiver := newSCProvisioner(privKey)
	actual, err := receiver.decryptKey(cipher)
	assert.NoError(t, err)

	assert.Equal(t, expected, actual)
}
