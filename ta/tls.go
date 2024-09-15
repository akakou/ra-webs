package ta

import (
	"crypto/rsa"
	"crypto/tls"
	"encoding/pem"
	"fmt"
)

const CERT_DIER_CACHE = "./tmp/ra-webs.cache"

func parsePemCertificate(raw []byte, privateKey *rsa.PrivateKey) (*tls.Certificate, error) {
	certs := make([][]byte, 0)

	for block, rest := pem.Decode(raw); block != nil; block, rest = pem.Decode(rest) {
		if block.Type != "CERTIFICATE" {
			return nil, fmt.Errorf("unexpected block type %s", block.Type)
		}

		certs = append(certs, block.Bytes)
	}

	return &tls.Certificate{
		Certificate: certs,
		PrivateKey:  privateKey,
	}, nil
}

func (ap *TA) TLSConfig() (*tls.Config, error) {
	res, err := ap.Register()
	if err != nil {
		return nil, err
	}
	fmt.Print(res)

	resouce := IssueCertificate(ap.privateKey, ap.config.Domain, ap.config.Email)

	cert, err := parsePemCertificate(resouce.Certificate, ap.privateKey)

	if err != nil {
		return nil, err
	}

	return &tls.Config{
		Certificates: []tls.Certificate{
			*cert,
		},
	}, nil
}
