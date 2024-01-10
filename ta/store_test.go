package ta

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestPrivKeyStore(t *testing.T) {
	privKeyStoreW := privKeyStore{}
	privKeyStoreR := privKeyStore{}

	expected, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.Equal(t, nil, err)

	privKeyStoreW.privKey = expected

	err = privKeyStoreW.Store()
	assert.Equal(t, nil, err)

	err = privKeyStoreR.Load()
	assert.Equal(t, nil, err)
	assert.Equal(t, expected, privKeyStoreR.privKey)
}

// func TestCertificate(t *testing.T) {
// 	certStoreW := certStore{}
// 	certStoreR := certStore{}

// 	ra := NewRA(&RAConfig{
// 		TTPDomain: "ttp.example.com",
// 		Domain:    "ta.example.com",
// 	})

// 	_, expected, err := ra.generateKeyPair()
// 	assert.Equal(t, nil, err)

// 	certStoreW.cert = expected

// 	err = certStoreW.Store()
// 	assert.Equal(t, nil, err)

// 	err = certStoreR.Load()
// 	assert.Equal(t, nil, err)
// 	assert.Equal(t, expected.Certificate, certStoreR.cert.Certificate)
// }
