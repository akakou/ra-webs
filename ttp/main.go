package ttp

import (
	"net/http"

	"github.com/akakou/ra_webs/core"
	"github.com/labstack/echo/v4"
)


func NewTTPServer(dbConfig *DBConfig) *echo.Echo {
	e := echo.New()

	db, err := newtTAInfoDB(dbConfig)
	if err != nil {
		panic(err)
	}

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
