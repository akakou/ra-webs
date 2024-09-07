package main

import (
	"crypto/tls"
	"embed"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"

	"github.com/akakou/ra_webs/verifier"
	"github.com/akakou/ra_webs/verifier/api"
	"github.com/akakou/ra_webs/verifier/notifier"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
)

//go:embed views/*/*.html
var viewEmbedFiles embed.FS

//go:embed static/js/*.js static/js/*/*.js
var staticEmbedFiles embed.FS

const TMP_FOLDER_NAME = "views"

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func InjectSWHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Service-Worker-Allowed", "/app/redirect/")
		return next(c)
	}
}

func main() {
	verifierHost := os.Getenv("RA_WEBS_VERIFIER_HOST")

	autoTLSManager := autocert.Manager{
		Prompt: autocert.AcceptTOS,
		Cache:  autocert.DirCache("/var/www/.cache"),
	}

	e := echo.New()

	verifier, err := verifier.DefaultVerifier()
	if err != nil {
		panic(err)
	}

	api.Route(e, verifier)

	viewSubFS := echo.MustSubFS(viewEmbedFiles, "views")
	e.StaticFS("/app", viewSubFS)

	staticSubFS := echo.MustSubFS(staticEmbedFiles, "static")
	e.StaticFS("/static", staticSubFS)

	e.Use(middleware.Logger())
	e.Use(InjectSWHeader)

	e.Debug = true
	err = verifier.Setup(e)
	if err != nil {
		panic(err)
	}

	fmt.Printf("public: %v\nprivate: %v", verifier.Notifier.(*notifier.BrowserNotifier).VapidPublicKey, verifier.Notifier.(*notifier.BrowserNotifier).VapidPrivateKey)

	s := http.Server{
		Addr:    verifierHost + ":443",
		Handler: e,
		TLSConfig: &tls.Config{
			GetCertificate: autoTLSManager.GetCertificate,
			NextProtos:     []string{acme.ALPNProto},
		},
	}

	go verifier.Monitor.Run(verifier)
	// fmt.Printf("public: %v\nprivate: %v", verifier.Notifier.(*notifier.BrowserNotifier).VapidPublicKey, verifier.Notifier.(*notifier.BrowserNotifier).VapidPrivateKey)

	if err := s.ListenAndServeTLS("", ""); err != http.ErrServerClosed {
		e.Logger.Fatal(err)
	}
}
