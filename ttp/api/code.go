package api

import (
	"net/http"
	"strconv"

	goutils "github.com/akakou/go-utils"
	ttpcore "github.com/akakou/ra_webs/ttp/core"
	"github.com/akakou/ra_webs/ttp/ent/tacode"
	"github.com/labstack/echo/v4"
)

var postCodeApi = goutils.EchoRoute[ttpcore.TTP]{
	Method: goutils.POST,
	Path:   "/code",
	F: func(ttp *ttpcore.TTP) goutils.EchoRouteFunc {
		return func(c echo.Context) error {
			req := new(struct {
				Repository string `json:"repository"`
				CommitId   string `json:"commit_id"`
				UniqueID   []byte `json:"unique_id"`
			})

			if c.Bind(req) != nil {
				return c.String(http.StatusBadRequest, "bad attestation")
			}

			codeCreate := ttp.DB.Client.TACode.
				Create().
				SetRepository(req.Repository).
				SetCommitID(req.CommitId).
				SetUniqueID(req.UniqueID)

			code, err := codeCreate.Save(*ttp.DB.Ctx)

			if err != nil {
				return err
			}

			return c.String(http.StatusOK, strconv.Itoa(code.ID))
		}
	},
}

var postActivateCodeApi = goutils.EchoRoute[ttpcore.TTP]{
	Method: goutils.POST,
	Path:   "/code/:id/activate",
	F: func(ttp *ttpcore.TTP) goutils.EchoRouteFunc {
		return func(c echo.Context) error {
			paramId := c.Param("id")

			codeId, err := strconv.Atoi(paramId)
			if err != nil {
				return err
			}

			err = authenticateAdmin(ttp, c)
			if err != nil {
				return c.String(http.StatusUnauthorized, "token is invalid")
			}

			code, err := ttp.DB.Client.TACode.Get(*ttp.DB.Ctx, codeId)
			if err != nil {
				return err
			}

			_, err = code.Update().SetHasActivated(true).Save(*ttp.DB.Ctx)
			if err != nil {
				return err
			}

			return c.String(http.StatusOK, strconv.Itoa(code.ID))
		}
	},
}

var getCodeApi = goutils.EchoRoute[ttpcore.TTP]{
	Method: goutils.GET,
	Path:   "/code",
	F: func(ttp *ttpcore.TTP) goutils.EchoRouteFunc {
		return func(c echo.Context) error {
			activate := c.QueryParam("activate") != "false"

			code, err := ttp.DB.Client.TACode.Query().Where(tacode.HasActivated(activate)).All(*ttp.DB.Ctx)
			if err != nil {
				return err
			}

			return c.JSON(http.StatusOK, code)
		}
	},
}
