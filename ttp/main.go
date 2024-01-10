package ttp

import (
	"flag"
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

func NewTTPServer(dbConfig *DBConfig, templatePath string) *echo.Echo {
	e := echo.New()

	e.Renderer = &Template{
		templates: template.Must(template.ParseGlob(templatePath)),
	}

	db, err := newTtpDB(dbConfig)
	if err != nil {
		panic(err)
	}

	Route(e, db)

	return e
}

func DefaultTTPServer(templatePath string) *echo.Echo {
	dbType := flag.String("db_type", "sqlite3", "database type")
	dbConfig := flag.String("db_config", "file:ent?mode=memory&cache=shared&_fk=1", "database config")

	fmt.Printf("We use %s as database type and %s as database config\n", *dbType, *dbConfig)

	return NewTTPServer(&DBConfig{
		Type:   *dbType,
		Config: *dbConfig,
	}, templatePath)
}

func verifyAttestation(attestation string) error {
	return nil
}
