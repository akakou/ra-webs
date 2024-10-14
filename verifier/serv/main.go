package main

import (
	"embed"
	"fmt"
	"html/template"
	"io"

	"github.com/akakou/ra_webs/core"
	"github.com/akakou/ra_webs/verifier"
	"github.com/akakou/ra_webs/verifier/api"
	"github.com/akakou/ra_webs/verifier/notifier"
	"github.com/labstack/echo/v4"
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
		c.Response().Header().Set("Service-Worker-Allowed", "/app/verification-status/")
		return next(c)
	}
}

func main() {
	e := echo.New()
	apiGroup := e.Group(core.API_ROOT)

	verifier, err := verifier.DefaultVerifier()
	if err != nil {
		panic(err)
	}

	api.Route(apiGroup, verifier)

	viewSubFS := echo.MustSubFS(viewEmbedFiles, "views")
	e.StaticFS("/app", viewSubFS)

	staticSubFS := echo.MustSubFS(staticEmbedFiles, "static")
	e.StaticFS("/static", staticSubFS)

	// e.Use(middleware.Logger())
	e.Use(InjectSWHeader)

	e.Debug = true
	err = verifier.Setup(apiGroup)
	if err != nil {
		panic(err)
	}

	fmt.Printf("public: %v\nprivate: %v", verifier.Notifier.(*notifier.BrowserNotifier).VapidPublicKey, verifier.Notifier.(*notifier.BrowserNotifier).VapidPrivateKey)

	go verifier.Monitor.Run(verifier)
	// fmt.Printf("public: %v\nprivate: %v", verifier.Notifier.(*notifier.BrowserNotifier).VapidPublicKey, verifier.Notifier.(*notifier.BrowserNotifier).VapidPrivateKey)

	e.Logger.Fatal(e.Start(":8080"))
}
