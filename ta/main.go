package ta

import (
	"crypto/rsa"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/acme/autocert"
)

const CERT_DIER_CACHE = "/var/www/.cache"
const ATTEST_ENDPOINT = "/rawebs/attest"

func SetRaWebs(e *echo.Echo) {
	e.AutoTLSManager.Cache = autocert.DirCache(CERT_DIER_CACHE)

	e.GET(ATTEST_ENDPOINT, func(c echo.Context) error {
		publicKey := e.AutoTLSManager.Client.Key.Public()
		rsaPublicKey := publicKey.(*rsa.PublicKey)

		quote, err := attestateByAzure(rsaPublicKey)
		if err != nil {
			c.Error(err)
		}

		return c.String(http.StatusOK, quote)
	})

}
