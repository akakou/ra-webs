package ta

import (
	"crypto/rsa"
	"crypto/tls"
)

type RAConfig struct {
	TTPDomain string
	Domain    string
}

type RA struct {
	config       *RAConfig
	privKeyStore *privKeyStore
	certStore    *certStore
}

func NewRA(config *RAConfig) *RA {
	return &RA{
		config:       config,
		privKeyStore: &privKeyStore{},
		certStore:    &certStore{},
	}
}

func TLSConfig(ra *RA) (*tls.Config, error) {
	privKeyStore := privKeyStore{}
	certStore := certStore{}

	var privKey *rsa.PrivateKey
	var cert *tls.Certificate

	var err error

	if hasFileExists() {
		privKey, cert, err = ra.Load()
	} else {
		privKey, cert, err = ra.Provisioning(&privKeyStore, &certStore)
	}

	if err != nil {
		return nil, err
	}

	cert.PrivateKey = privKey

	tlsConfig := tls.Config{
		Certificates: []tls.Certificate{
			*cert,
		},
	}

	return &tlsConfig, nil

}
