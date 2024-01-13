package main

import (
	"io"
	"net/http"
	"os"

	"github.com/akakou/ra_webs/ta"
	"github.com/labstack/echo/v4"
)

func readFile(name string) string {
	file, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	bytes, _ := io.ReadAll(file)
	html := string(bytes)

	return html
}

func jsEndpoint(name string, e *echo.Echo) {
	e.GET("/ra-webs/static/"+name, func(c echo.Context) error {
		js := readFile(ta.STATIC_FOLDER + "/" + name)
		c.Response().Header().Set("Content-Type", "application/javascript")
		c.Response().Header().Set("Service-Worker-Allowed", "/")

		return c.String(http.StatusOK, js)
	})
}

func main() {
	e := echo.New()

	ra := ta.NewRA(&ta.RAConfig{
		Domain:    "localhost:8081",
		TTPDomain: "localhost:8082",
	})

	ta.STATIC_FOLDER = "../../ta/static"

	_, err := ra.TLSConfig()
	if err != nil {
		panic(err)
	}

	e.GET("/", func(c echo.Context) error {
		html := readFile("index.html")
		c.Response().Header().Set("Service-Worker-Allowed", "/aaa")

		return c.HTML(http.StatusOK, html)
	})

	jsEndpoint("crypto.js", e)
	jsEndpoint("app.js", e)
	jsEndpoint("sw.js", e)

	e.GET("/ra-webs/public_key.js", func(c echo.Context) error {
		c.Response().Header().Set("Service-Worker-Allowed", "/")

		js, err := ra.MakeServiceWorker()
		if err != nil {
			panic(err)
		}

		c.Response().Header().Set("Content-Type", "application/javascript")
		return c.String(http.StatusOK, js)
	})

	e.Start(":8081")

}
