package api

import (
	"fmt"

	verifiercore "github.com/akakou/ra_webs/verifier/core"
	"github.com/akakou/ra_webs/verifier/ent"
	"github.com/akakou/ra_webs/verifier/ent/service"
	"github.com/labstack/echo/v4"
)

const (
	ERROR_AUTHENTICATE_SERVICE      = "failed to authenticate service"
	ERROR_AUTHENTICATE_ADMIN        = "failed to authenticate admin"
	ERROR_ACCESS_DOMAIN_AUTH_TARGET = "failed to access domain auth target"
	ERROR_DOMAIN_AUTH_INVALID       = "domain auth token is invalid"
	ERROR_QUOTE_INVALID             = "quote is invalid"
)

var SCHEME = "https"

func authenticateService(verifier *verifiercore.Verifier, c echo.Context) (*ent.Service, error) {
	authorization := c.Request().Header["Authorization"][0]
	token := authorization[len("Bearer "):]

	service, err := verifier.DB.Client.Service.Query().Where(service.TokenEQ(token)).Only(*verifier.DB.Ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ERROR_AUTHENTICATE_SERVICE, err)
	}

	return service, nil
}

func authenticateAdmin(verifier *verifiercore.Verifier, c echo.Context) error {
	authorization := c.Request().Header["Authorization"][0]
	token := authorization[len("Bearer "):]

	if token != verifier.AdminToken {
		return fmt.Errorf(ERROR_AUTHENTICATE_ADMIN)
	}

	return nil
}
