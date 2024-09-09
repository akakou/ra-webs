package main

import (
	"crypto/tls"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	autoTLSManager := autocert.Manager{
		Prompt: autocert.AcceptTOS,
		Cache:  autocert.DirCache("/var/www/.cache"),
	}

	s := http.Server{
		Addr:    ":443",
		Handler: e,
		TLSConfig: &tls.Config{
			GetCertificate: autoTLSManager.GetCertificate,
			NextProtos:     []string{acme.ALPNProto},
		},
	}

	if err := s.ListenAndServeTLS("", ""); err != http.ErrServerClosed {
		e.Logger.Fatal(err)
	}
}
