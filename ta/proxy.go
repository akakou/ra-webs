package ta

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Director(req *http.Request, c echo.Context) {
}

func ModifyResponse(res *http.Response, c echo.Context) error {
	return nil
}
