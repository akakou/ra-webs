package main

import (
	"net/http"

	"github.com/akakou/ra_webs/ta"
	"github.com/labstack/echo/v4"
)

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "http://localhost:8081/redirect", http.StatusTemporaryRedirect)
}

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello")
	})

	ta.SetRaWebs(e)
	e.StartAutoTLS(":443")

	e.Logger.Fatal(e.StartAutoTLS(":443"))
}
