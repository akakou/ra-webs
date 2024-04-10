package main

import "github.com/akakou/ra_webs/ttp"

const PORT = ":8081"

func main() {
	e, err := ttp.DefaultTTPServer("../views/*.html")
	if err != nil {
		panic(err)
	}

	e.Debug = true
	e.Logger.Fatal(e.Start(PORT))
}
