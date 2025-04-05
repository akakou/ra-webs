package api

import (
	"fmt"
	"net/http"

	goutils "github.com/akakou/go-utils"
	"github.com/akakou/ra-webs/log/api/auth"
	"github.com/akakou/ra-webs/log/api/io"
	"github.com/akakou/ra-webs/log/core"
	"github.com/labstack/echo/v4"
)

var PostApi = goutils.EchoRoute[core.Log]{
	Method: goutils.POST,
	Path:   "/ta",
	F: func(log *core.Log) goutils.EchoRouteFunc {
		return func(c echo.Context) error {
			var req io.PostRequest

			err := auth.Authenticate(log, c)
			if err != nil {
				return c.String(http.StatusUnauthorized, "unauthorized")
			}

			err = c.Bind(&req)
			if err != nil {
				return c.String(http.StatusBadRequest, "invalid json body")
			}

			signature, err := core.Sign(log, &req)
			if err != nil {
				return c.String(http.StatusBadRequest, "failed to sign")
			}

			ta, err := log.DB.Client.TA.Create().
				SetRepository(req.Repository).
				SetCommitID(req.CommitId).
				SetEvidence(req.Evidence).
				SetSignature(signature).
				SetSignature([]byte("hello")).
				Save(*log.DB.Ctx)

			if err != nil {
				return c.String(http.StatusBadRequest, "failed to store the log")
			}

			return c.String(http.StatusOK, fmt.Sprintf("%d", ta.ID))
		}
	},
}
