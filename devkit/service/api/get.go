package api

import (
	"encoding/base64"
	"net/http"

	goutils "github.com/akakou/go-utils"
	"github.com/akakou/ra-webs/devkit/service"
	"github.com/akakou/ra-webs/devkit/service/ent/ta"
	"github.com/labstack/echo/v4"
)

var GetApi = goutils.EchoRoute[service.Service]{
	Method: goutils.GET,
	Path:   "/ta",
	F: func(service *service.Service) goutils.EchoRouteFunc {
		return func(c echo.Context) error {
			encodedPublicKey := c.QueryParam("public_key")
			publicKey, err := base64.URLEncoding.DecodeString(encodedPublicKey)

			if err != nil {
				return c.String(http.StatusBadRequest, "invalid base64")
			}

			ta, err := service.DB.Client.TA.Query().
				Where(ta.PublicKeyEQ(publicKey)).
				Only(*service.DB.Ctx)

			if err != nil {
				return c.String(http.StatusBadRequest, "invalid request")
			}

			return c.JSON(http.StatusOK, ta)
		}
	},
}
