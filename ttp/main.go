package ttp

import (
	"fmt"
	"html/template"
	"io"

	"github.com/akakou/ra_webs/ttp/core"
	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewTTPServer(ttp *core.TTP, templatePath string) *echo.Echo {
	e := echo.New()

	e.Renderer = &Template{
		templates: template.Must(template.ParseGlob(templatePath)),
	}

	Route(e, ttp)

	return e
}

func DefaultTTPServer(templatePath string) (*echo.Echo, error) {
	ttp, err := core.DefaultTTP()
	if err != nil {
		return nil, fmt.Errorf("failed to init ttp: %w", err)
	}

	return NewTTPServer(ttp, templatePath), nil
}
