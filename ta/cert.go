package ta

import (
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"

	"github.com/akakou/ra_webs/core"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

func (ap *TA) IssueAcmeCert(e *echo.Echo) error {
	quote, err := attestPublicKey(ap)
	if err != nil {
		return fmt.Errorf("failed to attest public key: %w", err)
	}

	acmeClient := acme.Client{DirectoryURL: autocert.DefaultACMEDirectory}
	acmeClient.Key = ap.PrivateKey

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

	return nil
}

func (ap *TA) IssueTTPCert(taId int, e *echo.Echo) error {
	quote, err := attestPublicKey(ap)
	if err != nil {
		return fmt.Errorf("failed to attest public key: %w", err)
	}

	issueCertUrl := fmt.Sprintf(TTP_ISSUE_CERT, taId)

	body := map[string]any{
		"quote": quote,
	}

	resp, err := ap.requestToTTP(issueCertUrl, body)
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
				PrivateKey:  ap.PrivateKey,
			},
		},
	}

	return nil
}
