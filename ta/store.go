package ta

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/edgelesssys/ego/ecrypto"
)

const PRIVATE_PATH = "private.key"
const CERTIFICATE_PATH = "certificate.pem"

func (ra *RA) Load() (*rsa.PrivateKey, *tls.Certificate, error) {
	err := ra.privKeyStore.Load()
	if err != nil {
		return nil, nil, err
	}

	err = ra.certStore.Load()
	if err != nil {
		return nil, nil, err
	}

	return ra.privKeyStore.privKey, ra.certStore.cert, nil
}

type privKeyStore struct {
	privKey *rsa.PrivateKey
}

func (store *privKeyStore) Load() error {
	bytes, err := os.ReadFile(PRIVATE_PATH)
	if err != nil {
		return fmt.Errorf("failed to read private storage: %w", err)
	}

	plaintext, err := ecrypto.Unseal(bytes, []byte{})
	if err != nil {
		return fmt.Errorf("failed to unseal private storage: %w", err)
	}

	privKey, err := x509.ParsePKCS1PrivateKey(plaintext)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}

	store.Set(privKey)

	return nil
}

func (store *privKeyStore) Set(privKey *rsa.PrivateKey) {
	store.privKey = privKey
}

func (store *privKeyStore) Store() error {
	raw := x509.MarshalPKCS1PrivateKey(store.privKey)

	cipher, err := ecrypto.SealWithProductKey(raw, []byte{})
	if err != nil {
		return err
	}

	return os.WriteFile(PRIVATE_PATH, cipher, 0644)
}

type certStore struct {
	cert *tls.Certificate
}

func (store *certStore) Load() error {
	raw, err := os.ReadFile(CERTIFICATE_PATH)
	if err != nil {
		return err
	}

	cert := tls.Certificate{
		Certificate: [][]byte{raw},
	}

	store.Set(&cert)

	return nil
}

func (store *certStore) Set(cert *tls.Certificate) {
	store.cert = cert
}

func (store *certStore) Store() error {
	return os.WriteFile(CERTIFICATE_PATH, store.cert.Certificate[0], 0644)
}

func hasFileExists() bool {
	_, err := os.Stat(PRIVATE_PATH)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}
