package main

import (
	"crypto/rand"
	"crypto/rsa"

	"github.com/akakou/ra-webs/log/api"
	"github.com/akakou/ra-webs/log/core"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	g := e.Group("/")

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	log, err := core.Default(key)
	if err != nil {
		panic(err)
	}

	api.GetApi.Set(g, log)
	api.PostApi.Set(g, log)

	defer log.DB.Close()

	err = e.Start(":8080")

	if err != nil {
		panic(err)
	}
}
