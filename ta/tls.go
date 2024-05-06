package ta

import (
	"crypto/tls"
)

const CERT_DIER_CACHE = "./tmp/ra-webs.cache"

func (ap *TA) TLSConfig() (*tls.Config, error) {
	cert := IssueCertificate(ap.privateKey, ap.config.Domain, ap.config.Email)

	return &tls.Config{
		GetCertificate: func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
			return &tls.Certificate{
				Certificate: [][]byte{cert.Certificate},
				PrivateKey:  ap.privateKey,
			}, nil
		},
	}, nil
}
