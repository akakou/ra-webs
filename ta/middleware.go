package ta

import (
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func MakeMiddleware(ttpUrl string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess, err := session.Get("session", c)
			if err != nil {
				return err
			}

			if c.Request().Header.Get("Referer") == ttpUrl {
				sess.Values["ttp_redirected"] = true
				sess.Save(c.Request(), c.Response())
			}

			if !sess.Values["ttp_redirected"].(bool) {
				return c.Redirect(http.StatusSeeOther, ttpUrl)
			} else {
				return next(c)
			}
		}
	}
}
