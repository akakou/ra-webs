package main

import (
	"log"
	"net/http"

	"github.com/akakou/ra_webs/ta"
	"github.com/labstack/echo/v4"
)

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "http://localhost:8081/redirect", http.StatusTemporaryRedirect)
}

func main() {
	config := ta.RAConfig{
		TTPDomain: "",
		Domain:    "",
		Email:     "",
		JSFolder:  "../static",
	}

	ra := ta.NewRA(&config)
	_, err := ra.TLSConfig()
	if err != nil {
		panic(err)
	}

	ta.DEBUG = true

	e := echo.New()
	e.Use(ra.Middleware())

	ra.EndPoints(e)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/aa", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, aa!")
	})

	e.Debug = true

	err = e.Start(":8081")

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
