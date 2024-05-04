package main

import (
	"embed"
	"html/template"
	"io"

	"github.com/akakou/ra_webs/core"
	"github.com/akakou/ra_webs/ttp"
	"github.com/labstack/echo/v4"
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
	e, err := ttp.DebugTTPServer()
	if err != nil {
		panic(err)
	}

	subFS := echo.MustSubFS(embedFiles, "views")
	e.StaticFS("/app", subFS)

	e.Debug = true
	e.Logger.Fatal(e.Start(core.TTPPort))
}
