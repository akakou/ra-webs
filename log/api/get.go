package api

import (
	"encoding/base64"
	"net/http"

	goutils "github.com/akakou/go-utils"
	"github.com/akakou/ra-webs/log"
	"github.com/akakou/ra-webs/log/ent/ta"
	"github.com/labstack/echo/v4"
)

var GetApi = goutils.EchoRoute[log.Log]{
	Method: goutils.GET,
	Path:   "/ta",
	F: func(log *log.Log) goutils.EchoRouteFunc {
		return func(c echo.Context) error {
			encodedPublicKey := c.QueryParam("publicKey")
			publicKey, err := base64.URLEncoding.DecodeString(encodedPublicKey)

			if err != nil {
				return c.String(http.StatusBadRequest, "invalid base64")
			}

			ta, err := log.DB.Client.TA.Query().
				Where(ta.PublicKeyEQ(publicKey)).
				Only(*log.DB.Ctx)

			if err != nil {
				return c.String(http.StatusBadRequest, "invalid request")
			}

			return c.JSON(http.StatusOK, ta)
		}
	},
}
