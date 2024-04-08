package ca

import (
	"net/http"

	goutils "github.com/akakou/go-utils"
	"github.com/labstack/echo"
)

var certApi = goutils.EchoRoute[ttpcore.TTP]{
	Method: goutils.GET,
	Path:   "/cert",
	F: func(ttp *ttpcore.TTP) goutils.EchoRouteFunc {
		return func(c echo.Context) error {
			cert := ttp.CA.Certificate.Raw
			return c.Blob(http.StatusOK, "application/x-x509-ca-cert", cert)
		}
	},
}
