package ta

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509/pkix"
	"fmt"

	"github.com/akakou/ra_webs/core"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

const CERT_DIER_CACHE = "/var/www/.cache"
const ATTEST_ENDPOINT = "/rawebs/attest"

var QUOTE_OBJECT_ID = []int{1, 3, 6, 1, 4, 1, 48181, 1, 1}

func SetRaWebs(e *echo.Echo) error {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		return fmt.Errorf("failed to generate rsa key: %w", err)
	}

	pubKey := privKey.Public()
	quote, err := core.AttestByAzure(pubKey.(*rsa.PublicKey))
	if err != nil {
		return fmt.Errorf("failed to attestate by azure: %w", err)
	}

	acmeClient := acme.Client{DirectoryURL: autocert.DefaultACMEDirectory}
	acmeClient.Key = privKey

	e.AutoTLSManager = autocert.Manager{
		Client: &acmeClient,
		Cache:  autocert.DirCache(CERT_DIER_CACHE),
		ExtraExtensions: []pkix.Extension{
			{
				Id:       QUOTE_OBJECT_ID,
				Critical: false,
				Value:    []byte(quote),
			},
		},
	}

	return nil
}
