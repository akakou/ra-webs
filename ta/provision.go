package ta

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"time"
)

func (raConfig *RAConfig) generateKeyPair() (*tls.Config, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		return nil, err
	}

	template := &x509.Certificate{
		SerialNumber: &big.Int{},
		Subject:      pkix.Name{CommonName: raConfig.Domain},
		NotAfter:     time.Now().Add(time.Hour),
		DNSNames:     []string{raConfig.Domain},
	}

	cert, err := x509.CreateCertificate(rand.Reader, template, template, &priv.PublicKey, priv)

	if err != nil {
		return nil, err
	}

	tlsCfg := tls.Config{
		Certificates: []tls.Certificate{
			{
				Certificate: [][]byte{cert},
				PrivateKey:  priv,
			},
		},
	}

	return &tlsCfg, nil
}
