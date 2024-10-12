package utils

import (
	"crypto/x509"
	"errors"

	core "github.com/akakou/ctstream/core"
	ctx509 "github.com/google/certificate-transparency-go/x509"
)

func Callbacks(certs []*ctx509.Certificate, params any, callback core.Callback, err error) {
	for i, c := range certs {
		callback(c, i, params, err)
	}
}

func ReformatCertificate(cert *x509.Certificate) (*ctx509.Certificate, error) {
	return ctx509.ParseCertificate(cert.Raw)
}

func ReformatCertificates(certs []x509.Certificate) ([]*ctx509.Certificate, error) {
	var err error
	var result []*ctx509.Certificate

	for _, cert := range certs {
		c, tmpError := ReformatCertificate(&cert)
		err = errors.Join(err, tmpError)
		result = append(result, c)
	}

	return result, err
}
