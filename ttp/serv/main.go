package main

import (
	"embed"
	"html/template"
	"io"

	"github.com/akakou/ra_webs/core"
	"github.com/akakou/ra_webs/ttp"
	"github.com/akakou/ra_webs/ttp/api"
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
		c.Response().Header().Set("Service-Worker-Allowed", "/app/redirect/")
		return next(c)
	}
}

func main() {
	e := echo.New()
	ttp, err := ttp.DefaultTTP()
	if err != nil {
		panic(err)
	}

	api.Route(e, ttp)

	viewSubFS := echo.MustSubFS(viewEmbedFiles, "views")
	e.StaticFS("/app", viewSubFS)

	staticSubFS := echo.MustSubFS(staticEmbedFiles, "static")
	e.StaticFS("/static", staticSubFS)

	e.Use(middleware.Logger())
	e.Use(InjectSWHeader)

	e.Debug = true
	err = ttp.Setup(e)
	if err != nil {
		panic(err)
	}

	// go ttp.Audit.Run(ttp)
	e.Logger.Fatal(e.Start(core.TTPPort))
}
