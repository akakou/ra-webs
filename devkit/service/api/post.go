package api

import (
	"fmt"
	"net/http"

	goutils "github.com/akakou/go-utils"
	"github.com/akakou/ra-webs/devkit/service"
	"github.com/akakou/ra-webs/devkit/service/api/auth"
	"github.com/akakou/ra-webs/devkit/service/api/io"
	"github.com/labstack/echo/v4"
)

var PostApi = goutils.EchoRoute[service.Service]{
	Method: goutils.POST,
	Path:   "/ta",
	F: func(service *service.Service) goutils.EchoRouteFunc {
		return func(c echo.Context) error {
			var req io.PostRequest

			err := auth.Authenticate(service, c)
			if err != nil {
				return c.String(http.StatusUnauthorized, "unauthorized")
			}

			err = c.Bind(&req)
			if err != nil {
				return c.String(http.StatusBadRequest, "invalid json body")
			}

			ta, err := service.DB.Client.TA.Create().
				SetRepository(req.Repository).
				SetPublicKey(req.PublicKey).
				SetCommitID(req.CommitId).
				SetEvidence(req.Evidence).
				Save(*service.DB.Ctx)

			if err != nil {
				return c.String(http.StatusBadRequest, "failed to store the log")
			}

			return c.String(http.StatusOK, fmt.Sprintf("%d", ta.ID))
		}
	},
}
