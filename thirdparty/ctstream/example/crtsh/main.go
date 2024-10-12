package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/akakou/ctstream/core"
	"github.com/akakou/ctstream/monitor/crtsh"
	ctx509 "github.com/google/certificate-transparency-go/x509"
)

func main() {
	core.DefaultEpochSleep = time.Second * 10

	m, err := crtsh.DefaultCTsStream([]string{
		"test2.ochano.co",
	}, context.Background())

	if err != nil {
		fmt.Printf("Failed to create new ctstream: ")
		os.Exit(1)
	}

	m.Start(func(cert *ctx509.Certificate, i int, opt any, err error) {
		if err != nil {
			fmt.Printf("Failed to fetch %v: \n", err)
		}

		params := opt.(*crtsh.CrtshCTParams)

		fmt.Printf("%v: %v (target: %v)\n", params.ID, cert.DNSNames, params.Client.Domain)
	})

	go func() {
		time.Sleep(40 * time.Second)
		m.Stop()
	}()

	m.Await()
}
