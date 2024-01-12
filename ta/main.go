package ta

import (
	"crypto/rsa"
	"crypto/tls"
)

type RAConfig struct {
	TTPDomain string
	Domain    string
	Email     string
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

func (ra *RA) TLSConfig()(*tls.Config, error) {
	var privKey *rsa.PrivateKey
	var cert *tls.Certificate

	var err error

	if hasFileExists() {
		privKey, cert, err = ra.Load()
	} else {
		privKey, cert, err = ra.Provisioning()
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
