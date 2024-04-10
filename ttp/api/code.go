package api

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"

	goutils "github.com/akakou/go-utils"
	"github.com/akakou/ra_webs/ttp/builder"
	ttpcore "github.com/akakou/ra_webs/ttp/core"
	"github.com/akakou/ra_webs/ttp/ent/tacode"
	"github.com/labstack/echo/v4"
)

var BUILD_DOCKER_PATH = "./builder"

var PostCodeApi = goutils.EchoRoute[ttpcore.TTP]{
	Method: goutils.POST,
	Path:   "/code",
	F: func(ttp *ttpcore.TTP) goutils.EchoRouteFunc {
		return func(c echo.Context) error {
			service, err := authenticateService(ttp, c)
			if err != nil {
				return c.String(http.StatusUnauthorized, "token is invalid")
			}

			req := new(struct {
				Repository string `json:"repository"`
			})

			if c.Bind(req) != nil {
				return c.String(http.StatusBadRequest, "bad attestation")
			}

			rawSha256 := sha256.Sum256([]byte(req.Repository))
			sha256 := fmt.Sprintf("%x", rawSha256)

			commitId, uniqueIdString, err := builder.Build(sha256, req.Repository)
			if err != nil {
				return err
			}

			uniqueId, _ := hex.DecodeString(uniqueIdString)

			codeCreate := ttp.DB.Client.TACode.
				Create().
				SetRepository(req.Repository).
				SetCommitID(commitId).
				SetUniqueID(uniqueId).
				SetIsActive(true).
				SetService(service)

			_, err = codeCreate.Save(*ttp.DB.Ctx)

			if err != nil {
				return err
			}

			return c.String(http.StatusOK, uniqueIdString)
		}
	},
}

var GetCodeApi = goutils.EchoRoute[ttpcore.TTP]{
	Method: goutils.GET,
	Path:   "/code",
	F: func(ttp *ttpcore.TTP) goutils.EchoRouteFunc {
		return func(c echo.Context) error {
			code, err := ttp.DB.Client.TACode.Query().
				Where(tacode.IsActive(true)).
				All(*ttp.DB.Ctx)
			if err != nil {
				return err
			}

			return c.JSON(http.StatusOK, code)
		}
	},
}
