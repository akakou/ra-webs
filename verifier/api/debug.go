package api

import (
	"fmt"

	golangutils "github.com/akakou/go-utils"
	goutils "github.com/akakou/go-utils"
	"github.com/akakou/ra_webs/verifier/core"
	"github.com/labstack/echo/v4"
)

var GetReloadDBApi = goutils.EchoRoute[core.Verifier]{
	Method: goutils.GET,
	Path:   "/reloaddb",
	F: func(verifier *core.Verifier) goutils.EchoRouteFunc {
		return func(c echo.Context) error {
			dbType := golangutils.GetEnv("DB_TYPE", "sqlite3")
			dbConfig := golangutils.GetEnv("DB_CONFIG", "file:ent?mode=memory&cache=shared&_fk=1")
			fmt.Printf("We use %s as database type and %s as database config\n", dbType, dbConfig)

			dbc := core.DBConfig{
				Type:   dbType,
				Config: dbConfig,
			}

			db, err := core.NewDB(&dbc)
			if err != nil {
				return err
			}

			_, err = verifier.DB.Client.TAServer.
				Query().
				All(*verifier.DB.Ctx)

			if err != nil {
				return err
			}

			verifier.DB = db

			return nil
		}
	},
}
