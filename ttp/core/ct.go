package core

import "github.com/labstack/echo/v4"

type CT interface {
	Setup(*echo.Echo, *TTP) error
	Subscribe(string, *TTP) error
	Run(ttp *TTP)
}
