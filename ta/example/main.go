package main

import (
	"net/http"

	goutils "github.com/akakou/go-utils"
	"github.com/akakou/ra_webs/core"
	"github.com/akakou/ra_webs/ta"
	"github.com/labstack/echo/v4"
)

var Token = goutils.GetEnv("RA_WEBS_SERVICE_TOKEN", core.DEBUG_TOKEN)
var Domain = goutils.GetEnv("RA_WEBS_TA_DOMAIN", "localhost")
var TTPBase = goutils.GetEnv("RA_WEBS_TTP_BASE", "http://localhost")
var Repository = goutils.GetEnv("RA_WEBS_TA_REPOSITORY", "github.com/akakou/ra_webs")

const REDIRECT_PATH = "/app/redirect"

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, Repository+REDIRECT_PATH, http.StatusTemporaryRedirect)
}

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello")
	})

	core.EnableDebug()

	ta, err := ta.InitTA(
		&ta.Config{
			Token:      Token,
			Domain:     Domain,
			TTP:        TTPBase + core.TTPPort,
			Repository: Repository,
		},
	)

	if err != nil {
		panic(err)
	}

	err = ta.Config(e)
	if err != nil {
		panic(err)
	}

	e.Debug = true
	e.Logger.Fatal(e.StartAutoTLS(core.TAPort))
}
