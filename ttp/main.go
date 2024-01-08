package main

import (
	"flag"
	"net/http"

	"github.com/akakou/ra_webs/core"
	"github.com/labstack/echo/v4"
)

const PORT = ":1323"

func main() {
	dbType := flag.String("db_type", "sqlite3", "database type")
	dbConfig := flag.String("db_config", "file:ent?mode=memory&cache=shared&_fk=1", "database config")

	db, err := newtTAInfoDB(*dbType, *dbConfig)
	if err != nil {
		panic(err)
	}

	e := NewRouter(db)
	e.Logger.Fatal(e.Start(PORT))
}

func NewRouter(db *taInfoDB) *echo.Echo {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/register", func(c echo.Context) error {
		provReq := new(core.ProvisioningRequest)

		if err := c.Bind(provReq); err != nil {
			return err
		}

		if verifyAttestation(provReq.Attestation) != nil {
			return c.String(http.StatusBadRequest, "bad attestation")
		}

		if db.store(provReq) != nil {
			return c.String(http.StatusInternalServerError, "internal error")
		}

		return c.String(http.StatusOK, "ok")
	})

	return e
}

func verifyAttestation(attestation string) error {
	return nil
}
