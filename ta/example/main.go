package main

import (
	"fmt"
	"net/http"

	goutils "github.com/akakou/go-utils"
	"github.com/akakou/ra_webs/core"
	"github.com/akakou/ra_webs/ta"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var Token = goutils.GetEnv("RA_WEBS_SERVICE_TOKEN", core.DEBUG_TOKEN)
var Domain = goutils.GetEnv("RA_WEBS_TA_DOMAIN", "localhost")
var TTPBase = goutils.GetEnv("RA_WEBS_TTP_BASE", "http://localhost"+core.TTPPort)
var Repository = goutils.GetEnv("RA_WEBS_TA_REPOSITORY", "github.com/akakou/ra_webs")

const REDIRECT_PATH = "/app/redirect/"

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	e.GET("/", func(c echo.Context) error {
		_, err := c.Cookie("isFirstAccess")
		if err != nil {
			c.SetCookie(&http.Cookie{
				Name:  "isFirstAccess",
				Value: "true",
			})

			html := fmt.Sprintf("<script>location.href = '%v'</script>", TTPBase+REDIRECT_PATH)
			c.HTML(http.StatusFound, html)
		}

		return c.String(http.StatusOK, "Hello, World!")
	})

	// core.EnableDebug()

	ta, err := ta.InitTA(
		&ta.Config{
			Token:      Token,
			Domain:     Domain,
			TTP:        TTPBase,
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
	// e.Start(":8002")
	e.Logger.Fatal(e.StartAutoTLS(core.TAPort))
}
