package main

import (
	"net/http"

	golangutils "github.com/akakou/go-utils"
	"github.com/akakou/ra_webs/core"
	"github.com/akakou/ra_webs/ta"
	"github.com/labstack/echo/v4"
)

var ttpRedirectUrl = golangutils.GetEnv("TTP_REDIRECT", "http://localhost:8000/redirect")

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, ttpRedirectUrl, http.StatusTemporaryRedirect)
}

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello")
	})

	core.EnableDebug()
	ta, err := ta.DefaultTA()

	if err != nil {
		panic(err)
	}

	e.AutoTLSManager, err = ta.TLSConfig()

	if err != nil {
		panic(err)
	}

	e.Logger.Fatal(e.StartAutoTLS(core.TAPort))
}
