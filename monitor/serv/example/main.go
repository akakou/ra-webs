package main

import (
	"embed"
	"fmt"
	"html/template"
	"io"

	"github.com/akakou/ra-webs/core"
	"github.com/akakou/ra-webs/monitor/serv"
	"github.com/akakou/ra-webs/monitor/serv/api"

	browsernotify "github.com/akakou/ra-webs/monitor/notifier/browser"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//go:embed views/*/*.html
var viewEmbedFiles embed.FS

//go:embed static/js/*.js static/js/*/*.js
var staticEmbedFiles embed.FS

const TMP_FOLDER_NAME = "views"

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func InjectSWHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Service-Worker-Allowed", "/app/verification-status/")
		return next(c)
	}
}

func main() {
	e := echo.New()
	apiGroup := e.Group(core.API_ROOT)

	s, err := serv.Default()
	if err != nil {
		panic(err)
	}

	defer s.Close()

	api.Route(apiGroup, s)

	viewSubFS := echo.MustSubFS(viewEmbedFiles, "views")
	e.StaticFS("/app", viewSubFS)

	staticSubFS := echo.MustSubFS(staticEmbedFiles, "static")
	e.StaticFS("/static", staticSubFS)

	e.Use(middleware.Logger())
	e.Use(InjectSWHeader)

	e.Debug = true
	fmt.Printf("public: %v\nprivate: %v",
		s.Monitor.Notifier.(*browsernotify.BrowserNotifier).VapidPublicKey,
		s.Monitor.Notifier.(*browsernotify.BrowserNotifier).VapidPrivateKey)

	err = s.Run(":8080", e)
	e.Logger.Fatal(err)
}
