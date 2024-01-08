package main

import "github.com/akakou/ra_webs/ttp"

const PORT = ":1323"

func main() {
	e := ttp.DefaultTTPServer()
	e.Logger.Fatal(e.Start(PORT))
}
