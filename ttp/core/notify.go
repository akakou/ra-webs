package core

import (
	"github.com/akakou/ra_webs/ttp/ent"
	"github.com/labstack/echo/v4"
)

type Notify interface {
	Notify(msg []byte, serv *ent.TAServer, ttp *TTP) error
	Setup(e *echo.Echo, ttp *TTP) error // todo: fix
}
