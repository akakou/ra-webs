package api

import (
	"fmt"
	"net/http"
	"reflect"

	goutils "github.com/akakou/go-utils"
	"github.com/akakou/ra_webs/core"
	"github.com/akakou/ra_webs/verifier/builder"
	verifiercore "github.com/akakou/ra_webs/verifier/core"
	"github.com/akakou/ra_webs/verifier/ent"
	"github.com/akakou/ra_webs/verifier/ent/taserver"
	"github.com/labstack/echo/v4"
)

var RegisterApi = goutils.EchoRoute[verifiercore.Verifier]{
	Method: goutils.POST,
	Path:   "/ta",
	F: func(verifier *verifiercore.Verifier) goutils.EchoRouteFunc {
		return func(c echo.Context) error {
			service, err := authenticateService(verifier, c)
			if err != nil {
				return c.String(http.StatusUnauthorized, "token is invalid")
			}

			var req core.RegisterRequest
			err = c.Bind(&req)
			if err != nil {
				return err
			}

			exists, err := DomainExist(req.Domain, verifier)
			if err != nil {
				return err
			}

			output, err := builder.Build(req.Repository)
			if err != nil {
				return err
			}

			err = CheckValidity(output.UniqueId, req, exists, service, verifier)
			if err != nil {
				return err
			}

			err = Register(&req, output, exists, service, verifier)
			if err != nil {
				return err
			}

			return c.String(http.StatusOK, "ok")

		}
	},
}

func DomainExist(domain string, verifier *verifiercore.Verifier) (bool, error) {
	exists, err := verifier.DB.Client.TAServer.Query().Where(taserver.Domain(domain)).Exist(*verifier.DB.Ctx)

	if err != nil {
		return exists, err
	}

	return exists, nil
}

func CheckValidity(uniqueId []byte, req core.RegisterRequest, exist bool, service *ent.Service, verifier *verifiercore.Verifier) error {
	report, err := core.VerifyServer(req.Quote, req.PublicKey, service.Token)

	if err != nil {
		return err
	}

	fmt.Printf("Unique ID: %x == %x\n", report.UniqueID, uniqueId)
	if !reflect.DeepEqual(report.UniqueID, uniqueId) {
		return fmt.Errorf(ERROR_QUOTE_INVALID)
	}

	err = verifier.Monitor.PreCheck(req.Domain, exist, verifier)
	if err != nil {
		return err
	}

	return nil
}

func Register(req *core.RegisterRequest, output *builder.BuildOutput, exist bool, service *ent.Service, verifier *verifiercore.Verifier) error {
	err := verifier.Monitor.Register(req.Domain, exist, verifier)
	if err != nil {
		return err
	}

	err = verifiercore.NotifierUpdate(req.Domain, verifier)
	if err != nil {
		return err
	}

	code, err := RegisterCode(output, &req.CodeRequest, service, verifier)

	if err != nil {
		return err
	}

	err = RegisterServer(&req.ServerRequest, code, service, verifier)
	if err != nil {
		return err
	}

	return nil
}

func RegisterServer(req *core.ServerRequest, code *ent.TACode, service *ent.Service, verifier *verifiercore.Verifier) error {
	taServerCreate := verifier.DB.Client.TAServer.
		Create().
		SetDomain(req.Domain).
		SetService(service).
		SetCode(code).
		SetPublicKey(req.PublicKey).
		SetQuote(req.Quote).
		SetIsActive(false)

	_, err := taServerCreate.Save(*verifier.DB.Ctx)
	if err != nil {
		return err
	}

	return nil
}

func RegisterCode(output *builder.BuildOutput, req *core.CodeRequest, service *ent.Service, verifier *verifiercore.Verifier) (*ent.TACode, error) {
	codeCreate := verifier.DB.Client.TACode.
		Create().
		SetRepository(req.Repository).
		SetCommitID(output.CommitId).
		SetUniqueID(output.UniqueId).
		SetIsActive(true).
		SetService(service)

	code, err := codeCreate.Save(*verifier.DB.Ctx)

	if err != nil {
		return nil, err
	}

	return code, nil
}
