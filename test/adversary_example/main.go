package main

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

var MonitorBase = "https://rare-charter-vancouver-many.trycloudflare.com"

func main() {
	e := echo.New()

	body := "<button onclick=\"window.open('%v/');\">Monitor Page (https://rare-charter-vancouver-many.trycloudflare.com)</button><br/>"
	body = fmt.Sprintf(body, MonitorBase)

	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, body)
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
