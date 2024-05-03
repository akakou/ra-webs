package main

import (
	"flag"
	"fmt"

	"github.com/akakou/ra_webs/service"
	"github.com/labstack/echo/v4"
)

const SERVER_MODE = "server"
const CODE_MODE = "code"

func main() {
	res := ""
	var err error = nil

	flag.Parse()
	args := flag.Args()

	if len(args) != 2 {
		err := fmt.Errorf("invalid number of args: %v", len(args))
		panic(err)
	}

	service := service.DefaultService()

	if args[0] == SERVER_MODE {
		e := echo.New()
		res, err = service.PostServer(args[1], e)
	} else if args[0] == CODE_MODE {
		go service.ServAuthDomain()
		res, err = service.PostCode(args[1])
	} else {
		panic("invalid mode")
	}

	if err != nil {
		panic(err)
	}

	fmt.Println(res)
}
