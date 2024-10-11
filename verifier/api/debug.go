package api

import (
	"fmt"

	goutils "github.com/akakou/go-utils"
	golangutils "github.com/akakou/golang-utils"
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

			adminToken, err := goutils.RandomHex(core.RANDOM_SIZE)
			if err != nil {
				return err
			}
			adminToken = golangutils.GetEnv("ADMIN_TOKEN", adminToken)

			fmt.Printf("Admin token generated: %s\n", adminToken)

			dbc := core.DBConfig{
				Type:   dbType,
				Config: dbConfig,
			}

			db, err := core.NewDB(&dbc)
			if err != nil {
				return err
			}

			verifier.DB = db

			return nil
		}
	},
}
