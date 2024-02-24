package ttp

import (
	"net/http"
	"strconv"

	"github.com/akakou/ra_webs/ttp/ent/tacode"
	"github.com/labstack/echo/v4"
)

var postCodeApi = echoRoute{
	method: POST,
	path:   "/code",
	f: func(auditor *Auditor) echoRouteFunc {
		return func(c echo.Context) error {
			req := new(struct {
				Repository string `json:"repository"`
				CommitId   string `json:"commit_id"`
				UniqueID   []byte `json:"unique_id"`
			})

			if c.Bind(req) != nil {
				return c.String(http.StatusBadRequest, "bad attestation")
			}

			codeCreate := auditor.db.Client.TACode.
				Create().
				SetRepository(req.Repository).
				SetCommitID(req.CommitId).
				SetUniqueID(req.UniqueID)

			code, err := codeCreate.Save(*auditor.db.Ctx)

			if err != nil {
				c.Error(err)
			}

			return c.String(http.StatusOK, strconv.Itoa(code.ID))
		}
	},
}

var postActivateCodeApi = echoRoute{
	method: POST,
	path:   "/code/:id/activate",
	f: func(auditor *Auditor) echoRouteFunc {
		return func(c echo.Context) error {
			paramId := c.Param("id")

			codeId, err := strconv.Atoi(paramId)
			if err != nil {
				c.Error(err)
				return err
			}

			authorization := c.Request().Header["Authorization"][0]
			token := authorization[len("Bearer "):]

			if token != auditor.adminToken {
				return c.String(http.StatusUnauthorized, "token is invalid")
			}

			code, err := auditor.db.Client.TACode.Get(*auditor.db.Ctx, codeId)
			if err != nil {
				c.Error(err)
				return err
			}

			_, err = code.Update().SetActivate(true).Save(*auditor.db.Ctx)
			if err != nil {
				c.Error(err)
				return err
			}

			return c.String(http.StatusOK, strconv.Itoa(code.ID))
		}
	},
}

var getCodeApi = echoRoute{
	method: GET,
	path:   "/code",
	f: func(auditor *Auditor) echoRouteFunc {
		return func(c echo.Context) error {
			activate := c.QueryParam("activate") != "false"

			code, err := auditor.db.Client.TACode.Query().Where(tacode.Activate(activate)).All(*auditor.db.Ctx)
			if err != nil {
				c.Error(err)
				return err
			}

			return c.JSON(http.StatusOK, code)
		}
	},
}
