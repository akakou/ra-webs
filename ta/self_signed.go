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

func (ra *RA) issueSelfSignedCert(privKey *rsa.PrivateKey, pubKey *rsa.PublicKey) (*tls.Certificate, error) {
	template := &x509.Certificate{
		SerialNumber: &big.Int{},
		Subject:      pkix.Name{CommonName: ra.config.Domain},
		NotAfter:     time.Now().Add(time.Hour),
		DNSNames:     []string{ra.config.Domain},
	}

	rawCert, err := x509.CreateCertificate(rand.Reader, template, template, pubKey, privKey)

	if err != nil {
		return nil, err
	}

	cert := tls.Certificate{
		Certificate: [][]byte{rawCert},
		PrivateKey:  privKey,
	}

	return &cert, nil
}
