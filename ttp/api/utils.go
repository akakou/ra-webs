package api

import (
	"fmt"
	"net/http"

	ttpcore "github.com/akakou/ra_webs/ttp/core"
	"github.com/akakou/ra_webs/ttp/ent"
	"github.com/akakou/ra_webs/ttp/ent/service"
	"github.com/labstack/echo/v4"
)

func authenticateService(auditor *ttpcore.TTP, c echo.Context) (*ent.Service, error) {
	authorization := c.Request().Header["Authorization"][0]
	token := authorization[len("Bearer "):]

	service, err := auditor.DB.Client.Service.Query().Where(service.TokenEQ(token)).First(*auditor.DB.Ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate service: %w", err)
	}

	return service, nil
}

func authenticateAdmin(auditor *ttpcore.TTP, c echo.Context) error {
	authorization := c.Request().Header["Authorization"][0]
	token := authorization[len("Bearer "):]

	if token != auditor.AdminToken {
		return c.String(http.StatusUnauthorized, "token is invalid")
	}

	return nil
}
