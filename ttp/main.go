package ttp

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"os"

	"github.com/akakou/metact"
	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewTTPServer(ct *metact.MetaCT, dbConfig *DBConfig, templatePath string) *echo.Echo {
	e := echo.New()

	e.Renderer = &Template{
		templates: template.Must(template.ParseGlob(templatePath)),
	}

	db, err := newAuditDB(dbConfig)
	if err != nil {
		panic(err)
	}

	auditor, err := NewAuditor(db, ct)
	if err != nil {
		panic(err)
	}

	auditServ := NewAuditServer(auditor)
	auditServ.Route(e)

	return e
}

func DefaultTTPServer(templatePath string) *echo.Echo {
	dbType := flag.String("db_type", "sqlite3", "database type")
	dbConfig := flag.String("db_config", "file:ent?mode=memory&cache=shared&_fk=1", "database config")

	metaAppId := os.Getenv("META_APP_ID")
	metaAccessToken := os.Getenv("META_ACCESS_TOKEN")

	fmt.Printf("We use %s as database type and %s as database config\n", *dbType, *dbConfig)

	ct := metact.NewCT(metaAppId, metaAccessToken)
	dbc := DBConfig{
		Type:   *dbType,
		Config: *dbConfig,
	}

	return NewTTPServer(ct, &dbc, templatePath)
}
