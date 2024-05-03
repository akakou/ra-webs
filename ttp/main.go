package ttp

import (
	"fmt"

	rawebscore "github.com/akakou/ra_webs/core"
	"github.com/akakou/ra_webs/ttp/api"
	"github.com/akakou/ra_webs/ttp/builder"
	"github.com/akakou/ra_webs/ttp/core"
	"github.com/akakou/ra_webs/ttp/ct"
	"github.com/labstack/echo/v4"
)

func NewTTPServer(ttp *core.TTP) *echo.Echo {
	e := echo.New()
	api.Route(e, ttp)

	return e
}

func DefaultTTPServer() (*echo.Echo, error) {
	ttp, err := core.DefaultTTP()
	if err != nil {
		return nil, fmt.Errorf("failed to init ttp: %w", err)
	}

	return NewTTPServer(ttp), nil
}

func DebugTTPServer() (*echo.Echo, error) {
	rawebscore.EnableDebug()
	ct.EnableDebug()
	builder.EnableDebug()

	ttp, err := core.DefaultTTP()
	if err != nil {
		return nil, fmt.Errorf("failed to init ttp: %w", err)
	}

	ttp.DB.Client.Service.Create().SetName("test").SetIsActive(true).SetToken(rawebscore.DEBUG_TOKEN).SaveX(*ttp.DB.Ctx)

	return NewTTPServer(ttp), nil
}
