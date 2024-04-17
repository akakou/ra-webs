package ct

import (
	"crypto/x509"
	"encoding/asn1"
	"errors"
	"reflect"

	metact "github.com/akakou/meta-ct"
)

func findCertExtensions(label asn1.ObjectIdentifier, cert *x509.Certificate) ([]byte, error) {
	for _, ext := range cert.Extensions {
		if reflect.DeepEqual(ext.Id, label) {
			return ext.Value, nil
		}
	}

	return []byte{}, errors.New(ERROR_EXTENSION_NOT_FOUND)
}

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
