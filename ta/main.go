package ta

import (
	"crypto/x509/pkix"
	"fmt"

	"github.com/akakou/ra_webs/core"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

const CERT_DIER_CACHE = "/var/www/.cache"
const DOMAIN_AUTH_PATH = "/ra-webs"

func (ap *TA) TLSConfig() (autocert.Manager, error) {
	quote, err := attestPublicKey(ap)
	if err != nil {
		return autocert.Manager{}, fmt.Errorf("%s: %w", ERROR_ATTEST_PUBLIC_KEY, err)
	}

	acmeClient := acme.Client{DirectoryURL: ap.Config.ACMEUrl}
	acmeClient.Key = ap.PrivateKey

	return autocert.Manager{
		Client: &acmeClient,
		Cache:  autocert.DirCache(CERT_DIER_CACHE),
		Prompt: autocert.AcceptTOS,
		ExtraExtensions: []pkix.Extension{
			{
				Id:       core.X509_EXTENSION_LABEL,
				Critical: false,
				Value:    []byte(quote),
			},
		},
	}, nil
}
