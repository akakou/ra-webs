package main

import (
	"fmt"

	"github.com/akakou/ra-webs/domainowner"
)

func main() {
	// == DNS server setup == //
	dnsServer, err := domainowner.NewDNSServer()
	if err != nil {
		e := fmt.Errorf("failed to initialize DNS server: %w", err)
		panic(e)
	}

	err = dnsServer.LoadConfig()
	if err != nil {
		e := fmt.Errorf("failed to initialize DNS server: %w", err)
		panic(e)
	}

	err = dnsServer.Start("0.0.0.0:53")
	panic(err)
}
