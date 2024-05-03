package main

import (
	"embed"
	"fmt"
	"html/template"
	"io"

	extractembed "github.com/akakou/extract-embed"
	"github.com/akakou/ra_webs/core"
	"github.com/akakou/ra_webs/ttp"
	"github.com/labstack/echo/v4"
)

//go:embed views/*.html
var embedFiles embed.FS

const TMP_FOLDER_NAME = "views"

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	embeded, err := extractembed.ExtractFs(TMP_FOLDER_NAME, &embedFiles)
	if err != nil {
		panic(err)
	}

	defer embeded.Close()

	fmt.Printf("extracted: %s\n", embeded.Path)

	tmpPath := fmt.Sprintf("%s/%s", embeded.Path, "views/*.html")

	e, err := ttp.DebugTTPServer()
	if err != nil {
		panic(err)
	}

	e.Renderer = &Template{
		templates: template.Must(template.ParseGlob(tmpPath)),
	}

	e.GET("/redirect", func(c echo.Context) error {
		return c.Render(200, "redirect.html", nil)
	})

	e.GET("/example", func(c echo.Context) error {
		return c.Render(200, "example.html", nil)
	})

	e.Debug = true
	e.Logger.Fatal(e.Start(core.TTPPort))
}
