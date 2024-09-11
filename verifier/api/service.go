package api

import (
	"net/http"

	goutils "github.com/akakou/go-utils"
	verifiercore "github.com/akakou/ra_webs/verifier/core"
	"github.com/labstack/echo/v4"
)

var PostServiceByAdmin = goutils.EchoRoute[verifiercore.Verifier]{
	Method: goutils.POST,
	Path:   "/service",
	F: func(verifier *verifiercore.Verifier) goutils.EchoRouteFunc {
		return func(c echo.Context) error {
			err := authenticateAdmin(verifier, c)

			if err != nil {
				return c.String(http.StatusUnauthorized, "token is invalid")
			}

			token, err := goutils.RandomHex(verifiercore.RANDOM_SIZE)

			if err != nil {
				return err
			}

			service, err := verifier.DB.Client.Service.
				Create().
				SetName("").
				SetToken(token).
				SetIsActive(true).
				Save(*verifier.DB.Ctx)

			if err != nil {
				return err
			}

			return c.String(http.StatusOK, service.Token)
		}
	},
}
