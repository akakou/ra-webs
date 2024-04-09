package ct

import (
	"crypto/x509"
	"encoding/asn1"
	"errors"
	"reflect"
)

func findCertExtensions(label asn1.ObjectIdentifier, cert *x509.Certificate) ([]byte, error) {
	for _, ext := range cert.Extensions {
		if reflect.DeepEqual(ext.Id, label) {
			return ext.Value, nil
		}
	}

	return []byte{}, errors.New("extension not found")
}
