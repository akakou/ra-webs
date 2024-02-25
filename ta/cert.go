package ta

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"

	"github.com/akakou/ra_webs/core"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

func (ap *AttestProxy) IssueAcmeCert(taId int, privKey *rsa.PrivateKey, quote string, e *echo.Echo) {
	acmeClient := acme.Client{DirectoryURL: autocert.DefaultACMEDirectory}
	acmeClient.Key = privKey

	e.AutoTLSManager = autocert.Manager{
		Client: &acmeClient,
		Cache:  autocert.DirCache(CERT_DIER_CACHE),
		ExtraExtensions: []pkix.Extension{
			{
				Id:       core.X509_EXTENSION_LABEL,
				Critical: false,
				Value:    []byte(quote),
			},
		},
	}
}

func (ap *AttestProxy) IssueTTPCert(taId int, privKey *rsa.PrivateKey, quote string, e *echo.Echo) error {
	issueCertUrl := fmt.Sprintf(TTP_ISSUE_CERT, taId)

	resp, err := ap.requestToTTP(issueCertUrl, string(""))
	if err != nil {
		return fmt.Errorf("failed to register: %w", err)
	}

	cert, err := x509.ParseCertificate(resp)
	if err != nil {
		return fmt.Errorf("failed to parse certificate: %w", err)
	}

	e.Server.TLSConfig = &tls.Config{
		Certificates: []tls.Certificate{
			{
				Certificate: [][]byte{cert.Raw},
				PrivateKey:  privKey,
			},
		},
	}

	return nil
}
