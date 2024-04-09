package api

import (
	"fmt"
	"net/http"

	ttpcore "github.com/akakou/ra_webs/ttp/core"
	"github.com/akakou/ra_webs/ttp/ent"
	"github.com/akakou/ra_webs/ttp/ent/service"
	"github.com/labstack/echo/v4"
)

const ERROR_AUTHENTICATE_SERVICE = "failed to authenticate service"

func authenticateService(ttp *ttpcore.TTP, c echo.Context) (*ent.Service, error) {
	authorization := c.Request().Header["Authorization"][0]
	token := authorization[len("Bearer "):]

	service, err := ttp.DB.Client.Service.Query().Where(service.TokenEQ(token)).Only(*ttp.DB.Ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ERROR_AUTHENTICATE_SERVICE, err)
	}

	return service, nil
}

func authenticateAdmin(ttp *ttpcore.TTP, c echo.Context) error {
	authorization := c.Request().Header["Authorization"][0]
	token := authorization[len("Bearer "):]

	if token != ttp.AdminToken {
		return c.String(http.StatusUnauthorized, "token is invalid")
	}

	return nil
}
