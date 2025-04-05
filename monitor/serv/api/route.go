package api

import (
	"github.com/akakou/ra_webs/monitor/serv"
	"github.com/labstack/echo/v4"
)

func Route(e *echo.Group, monitor *serv.MonitorServer) {
	GetServerApi.Set(e, monitor)
	GetServerFromDomainApi.Set(e, monitor)
	PostNotifierApi.Set(e, monitor)
}
