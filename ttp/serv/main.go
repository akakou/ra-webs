package main

import "github.com/akakou/ra_webs/ttp"

const PORT = ":8081"

func main() {
	e := ttp.DefaultTTPServer("../views/*.html")
	e.Debug = true
	e.Logger.Fatal(e.Start(PORT))
}
