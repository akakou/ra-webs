package ca

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"

	simplecertify "github.com/akakou/simple-certify"
)

func issueCertificate(domain string, uniqueId []byte, ca *simplecertify.Certifier) (*x509.Certificate, error) {
	templ := simplecertify.ServerTemplate()
	templ.PublicKey = ta.PublicKey
	templ.Subject = pkix.Name{
		Country:      []string{"Japan"},
		Organization: []string{"ra-webs"},
		Locality:     []string{"Kanagawa"},
		Province:     []string{"Yokohama"},
		CommonName:   domain,
	}

	templ.Issuer = ca.Certificate.Subject
	templ.Extensions = []pkix.Extension{
		{
			Id:    core.X509_EXTENSION_LABEL,
			Value: []byte(uniqueId),
		},
	}

	cert, err := ca.Certify(&templ)
	if err != nil {
		return nil, fmt.Errorf("failed to issue certificate: %w", err)
	}

	return cert, nil

}
