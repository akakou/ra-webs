package api

import (
	"errors"
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

			// fmt.Printf("domain: %v\n", domain)
			servs, err := verifier.DB.Client.TAServer.
				Query().
				Where(taserver.Domain(domain)).
				Where(taserver.HasActivated(true)).
				WithCode().
				WithViolation().
				Order(ent.Desc(taserver.FieldID)).
				All(*verifier.DB.Ctx)

			if err != nil {
				return errors.New("server is not found")
			}

			var res struct {
				TA      []*ent.TAServer `json:"ta"`
				IsValid bool            `json:"is_valid"`
			}

			res.TA = servs
			res.IsValid = checkTAValiditiy(servs)

			return c.JSON(http.StatusOK, res)
		}
	},
}
