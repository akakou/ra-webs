package ttp

import (
	"fmt"
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewTTPServer(auditor *Auditor, templatePath string) *echo.Echo {
	e := echo.New()

	e.Renderer = &Template{
		templates: template.Must(template.ParseGlob(templatePath)),
	}

	Route(e, auditor)

	return e
}

func DefaultTTPServer(templatePath string) (*echo.Echo, error) {
	auditor, err := DefaultAuditor()
	if err != nil {
		return nil, fmt.Errorf("failed to init auditor: %w", err)
	}

	return NewTTPServer(auditor, templatePath), nil
}
