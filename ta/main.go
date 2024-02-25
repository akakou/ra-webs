package ta

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/akakou/ra_webs/core"
	"github.com/labstack/echo/v4"
)

const CERT_DIER_CACHE = "/var/www/.cache"
const ATTEST_ENDPOINT = "/rawebs/attest"

func (ap *AttestProxy) Run() error {
	publicKey := ap.PrivateKey.Public()
	encodedPublicKey := x509.MarshalPKCS1PublicKey(publicKey.(*rsa.PublicKey))
	quote, err := core.AttestByAzure(encodedPublicKey)

	if err != nil {
		return fmt.Errorf("failed to attestate by azure: %w", err)
	}

	taId, err := ap.Register(publicKey.(*rsa.PublicKey))
	if err != nil {
		return fmt.Errorf("failed to register: %w", err)
	}

	ap.IssueAcmeCert(taId, ap.PrivateKey, quote, ap.Echo)
	ap.IssueTTPCert(taId, ap.PrivateKey, quote, ap.Echo)

	ap.Echo.Any("/*", func(c echo.Context) error {
		director := func(req *http.Request) {
			Director(req, c)
		}

		modifyResp := func(resp *http.Response) error {
			return ModifyResponse(resp, c)
		}

		proxy := httputil.ReverseProxy{
			Director:       director,
			ModifyResponse: modifyResp,
		}

		proxy.ServeHTTP(c.Response().Writer, c.Request())
		return nil
	})

	return nil
}
