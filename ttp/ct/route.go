package ct

import (
	"github.com/akakou/ra_webs/ttp/core"
	"github.com/labstack/echo/v4"
)

func Route(e *echo.Echo, auditor *core.TTP) {
	webhook().Set(e, auditor)
}
