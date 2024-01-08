package main

import (
	"net/http"

	"github.com/akakou/ra_webs/core"
	"github.com/labstack/echo/v4"
)

const PORT = ":1323"

func main() {
	e := NewRouter()
	e.Logger.Fatal(e.Start(PORT))
}

func NewRouter() *echo.Echo {
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

		if storeTAInfo(provReq) != nil {
			return c.String(http.StatusInternalServerError, "internal error")
		}

		return c.String(http.StatusOK, "ok")
	})

	return e
}

func verifyAttestation(attestation string) error {
	return nil
}

func storeTAInfo(req *core.ProvisioningRequest) error {
	return nil
}
