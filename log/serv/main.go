package main

import (
	"github.com/akakou/ra-webs/log"
	"github.com/akakou/ra-webs/log/api"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	g := e.Group("/api")

	l, err := log.Default()
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
