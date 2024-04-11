package ttp

import (
	"fmt"

	"github.com/akakou/ra_webs/ttp/api"
	"github.com/akakou/ra_webs/ttp/core"
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
