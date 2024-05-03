package ct

import (
	"crypto/x509"

	metact "github.com/akakou/meta-ct"
)

func MetaCertsToCerts(cs []metact.MetaCert) ([]x509.Certificate, error) {
	certs := []x509.Certificate{}

	for _, c := range cs {
		cert, err := c.Certificate()

		if err != nil {
			return []x509.Certificate{}, err
		}

		certs = append(certs, *cert)
	}

	return certs, nil
}

func subscribeCT(domain string, ct *metact.MetaCT) error {
	return ct.Subscribe(domain)
}

var SubscribeCT = subscribeCT
