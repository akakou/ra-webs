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
var embedFiles embed.FS

const TMP_FOLDER_NAME = "views"

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()
	ttp, err := ttp.DefaultTTP()
	if err != nil {
		panic(err)
	}

	api.Route(e, ttp)

	subFS := echo.MustSubFS(embedFiles, "views")
	e.StaticFS("/app", subFS)

	e.Use(middleware.Logger())

	e.Debug = true
	err = ttp.Setup(e)
	if err != nil {
		panic(err)
	}

	e.Logger.Fatal(e.Start(core.TTPPort))
}
