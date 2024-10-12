package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/akakou/ctstream/core"
	"github.com/akakou/ctstream/monitor/sslmate"

	ctx509 "github.com/google/certificate-transparency-go/x509"
)

func main() {
	core.DefaultEpochSleep = time.Second * 10

	m, err := sslmate.DefaultCTsStream([]string{
		"google.com",
	}, context.Background())

	if err != nil {
		fmt.Printf("Failed to create new ctstream: ")
		os.Exit(1)
	}

	err = m.Init()
	if err != nil {
		fmt.Printf("Failed to initialize ctstream: ")
		os.Exit(1)
	}

	m.Start(func(cert *ctx509.Certificate, i int, opt any, err error) {
		if err != nil {
			fmt.Printf("Failed to fetch %v: \n", err)
		}

		params := opt.(*sslmate.SSLMateCTParams)

		fmt.Printf("%v ~ %v: %v (target: %v)\n", params.Index.First, params.Index.Last, cert.DNSNames, params.Client.Domain)
	})

	go func() {
		time.Sleep(40 * time.Second)
		m.Stop()
	}()

	m.Await()

	fmt.Printf("last: %v: ", sslmate.GetFirst(m.Client))
}
