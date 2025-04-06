package main

import (
	"crypto/rand"
	"crypto/rsa"

	"github.com/akakou/ra-webs/log"
	"github.com/akakou/ra-webs/log/api"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	g := e.Group("/api")

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	l, err := log.Default(key)
	if err != nil {
		panic(err)
	}

	api.GetApi.Set(g, l)
	api.PostApi.Set(g, l)

	defer l.DB.Close()

	err = e.Start(":8080")

	if err != nil {
		panic(err)
	}
}
