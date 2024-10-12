package main

import (
	"fmt"

	"github.com/akakou/crtsh"
)

func main() {
	data, err := crtsh.Fetch("test.ochano.co", crtsh.EXCLUDE_EXPIRED)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	// fmt.Printf("Data: %v\n", data)
	fmt.Printf("Doamin[0]: %v\n", data[0].Certificate.DNSNames)
	fmt.Printf("Index[0]: %v\n", data[0].ID)
}
