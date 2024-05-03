package api

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"reflect"

	goutils "github.com/akakou/go-utils"
	"github.com/akakou/ra_webs/core"
	"github.com/akakou/ra_webs/ttp/builder"
	ttpcore "github.com/akakou/ra_webs/ttp/core"
	"github.com/akakou/ra_webs/ttp/ct"
	"github.com/akakou/ra_webs/ttp/ent"
	"github.com/labstack/echo/v4"
)

var RegisterApi = goutils.EchoRoute[ttpcore.TTP]{
	Method: goutils.POST,
	Path:   "/register",
	F: func(ttp *ttpcore.TTP) goutils.EchoRouteFunc {
		return func(c echo.Context) error {
			service, err := authenticateService(ttp, c)
			if err != nil {
				return c.String(http.StatusUnauthorized, "token is invalid")
			}

			var req core.RegisterRequest
			err = c.Bind(&req)
			if err != nil {
				return err
			}

			code, err := RegisterCode(&req.CodeRequest, service, ttp)

			if err != nil {
				return err
			}

			err = RegisterServer(&req.ServerRequest, code, service, ttp)
			if err != nil {
				return err
			}

			return c.String(http.StatusOK, "ok")

		}
	},
}

func RegisterServer(req *core.ServerRequest, code *ent.TACode, service *ent.Service, ttp *ttpcore.TTP) error {
	report, err := core.VerifyServer(req.Quote, req.PublicKey, service.Token)

	if err != nil {
		return err
	}

	if reflect.DeepEqual(report.UniqueID, code.UniqueID) {
		return fmt.Errorf(ERROR_QUOTE_INVALID)
	}

	err = ct.SubscribeCT(req.Domain, ttp.CT)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to subscribe %s due to %v", req.Domain, err)
	}

	taServerCreate := ttp.DB.Client.TAServer.
		Create().
		SetDomain(req.Domain).
		SetService(service).
		SetCode(code).
		SetPublicKey(req.PublicKey).
		SetHasActivated(false)

	_, err = taServerCreate.Save(*ttp.DB.Ctx)
	if err != nil {
		return err
	}

	return nil
}

func RegisterCode(req *core.CodeRequest, service *ent.Service, ttp *ttpcore.TTP) (*ent.TACode, error) {
	rawSha256 := sha256.Sum256([]byte(req.Repository))
	sha256 := fmt.Sprintf("%x", rawSha256)

	commitId, uniqueIdString, err := builder.Build(sha256, req.Repository)
	if err != nil {
		return nil, err
	}

	uniqueId, _ := hex.DecodeString(uniqueIdString)

	codeCreate := ttp.DB.Client.TACode.
		Create().
		SetRepository(req.Repository).
		SetCommitID(commitId).
		SetUniqueID(uniqueId).
		SetIsActive(false).
		SetService(service)

	code, err := codeCreate.Save(*ttp.DB.Ctx)

	if err != nil {
		return nil, err
	}

	return code, nil
}
