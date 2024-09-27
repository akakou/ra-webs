package api

import (
	"net/http"

	goutils "github.com/akakou/go-utils"
	verifiercore "github.com/akakou/ra_webs/verifier/core"
	"github.com/akakou/ra_webs/verifier/ent"
	"github.com/akakou/ra_webs/verifier/ent/taserver"
	"github.com/labstack/echo/v4"
)

var GetServerApi = goutils.EchoRoute[verifiercore.Verifier]{
	Method: goutils.GET,
	Path:   "/ta",
	F: func(verifier *verifiercore.Verifier) goutils.EchoRouteFunc {
		return func(c echo.Context) error {
			code, err := verifier.DB.Client.TAServer.Query().All(*verifier.DB.Ctx)
			if err != nil {
				return err
			}

			return c.JSON(http.StatusOK, code)
		}
	},
}

var GetServerFromDomainApi = goutils.EchoRoute[verifiercore.Verifier]{
	Method: goutils.GET,
	Path:   "/ta/:domain",
	F: func(verifier *verifiercore.Verifier) goutils.EchoRouteFunc {
		return func(c echo.Context) error {
			domain := c.Param("domain")

			type Res struct {
				TA      []*ent.TAServer `json:"ta"`
				IsValid bool            `json:"is_valid"`
				Message string          `json:"message"`
			}

			res := Res{}

			handleError := func(err error, r *Res) error {
				r.IsValid = false
				r.Message = err.Error()
				return c.JSON(http.StatusInternalServerError, res)
			}

			// fmt.Printf("domain: %v\n", domain)
			servs, err := verifier.DB.Client.TAServer.
				Query().
				Where(taserver.Domain(domain)).
				WithViolation().
				WithCode().
				WithService().
				Order(ent.Desc(taserver.FieldID)).
				All(*verifier.DB.Ctx)

			res.TA = servs

			if err != nil {
				return handleError(err, &res)
			}

			isValid1, err := checkViolationLogs(servs)
			if err != nil {
				return handleError(err, &res)
			}

			isValid2, err := checkTAValidity(servs[0], verifier)
			if err != nil {
				return handleError(err, &res)
			}

			res.IsValid = isValid1 && isValid2

			return c.JSON(http.StatusOK, res)
		}
	},
}
