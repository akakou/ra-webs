package api

import (
	"net/http"

	goutils "github.com/akakou/go-utils"
	"github.com/akakou/ra_webs/verifier/core"
	"github.com/labstack/echo/v4"
)

var PostNotifierApi = goutils.EchoRoute[core.Verifier]{
	Method: goutils.POST,
	Path:   "/notify",
	F: func(verifier *core.Verifier) goutils.EchoRouteFunc {
		return func(c echo.Context) error {
			err := authenticateAdmin(verifier, c)
			if err != nil {
				return c.String(http.StatusUnauthorized, "token is invalid")
			}

			var data struct {
				Domain  string `json:"domain"`
				Message string `json:"message"`
			}

			err = c.Bind(&data)
			if err != nil {
				return err
			}

			err = verifier.Notifier.Notify([]byte(data.Message), data.Domain, verifier)
			if err != nil {
				return err
			}

			return c.JSON(http.StatusOK, "ok")
		}
	},
}
