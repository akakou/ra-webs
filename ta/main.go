package ta

import (
	"crypto/x509/pkix"
	"fmt"
	"net/http"

	goutils "github.com/akakou/go-utils"
	"github.com/akakou/ra_webs/core"
	"github.com/labstack/echo/v4"
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
		ExtraExtensions: []pkix.Extension{
			{
				Id:       core.X509_EXTENSION_LABEL,
				Critical: false,
				Value:    []byte(quote),
			},
		},
	}, nil
}

func (ap *TA) DomainAuthServer(e *echo.Echo) {
	e.GET(DOMAIN_AUTH_PATH, func(c echo.Context) error {
		nonce, _ := goutils.RandomHex(64)
		h := core.DomainToken(ap.Config.Token, nonce)
		return c.String(http.StatusOK, h)
	})
}
