package main

import (
	"fmt"
	"net/http"

	"github.com/akakou/ra_webs/core"
	"github.com/akakou/ra_webs/ta"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const REDIRECT_PATH = "/app/redirect/"

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	config, err := ta.DefaultConfig()
	if err != nil {
		panic(err)
	}

	ta, err := ta.InitTA(config)
	if err != nil {
		panic(err)
	}

	e.GET("/", func(c echo.Context) error {
		_, err := c.Cookie("isFirstAccess")
		if err != nil {
			c.SetCookie(&http.Cookie{
				Name:  "isFirstAccess",
				Value: "true",
			})

			html := fmt.Sprintf("<script>location.href = '%v'</script>", config.TTP+REDIRECT_PATH)
			c.HTML(http.StatusFound, html)
		}

		return c.String(http.StatusOK, "Hello, World!")
	})

	e.Server.TLSConfig, err = ta.TLSConfig()
	if err != nil {
		panic(err)
	}

	e.Debug = true
	e.Logger.Fatal(e.Start(core.TAPort))
}
