package main

import (
	"net/http"

	"github.com/akakou/ra_webs/ta"
	"github.com/labstack/echo/v4"
)

const ttpUrl = "https://ttp.example.com"

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello")
	})

	ta.SetRaWebs(e, ttpUrl)
	e.StartAutoTLS(":443")

	e.Logger.Fatal(e.StartAutoTLS(":443"))
}
