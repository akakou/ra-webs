package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/akakou/ctstream/direct"
	ctX509 "github.com/google/certificate-transparency-go/x509"
)

func main() {
	m, err := direct.DefaultCTsStream([]string{
		// "https://ct.googleapis.com/logs/us1/argon2024/",
		// "https://ct.googleapis.com/logs/eu1/xenon2024/",
		// "https://ct.cloudflare.com/logs/nimbus2024/",
		// "https://yeti2024.ct.digicert.com/log/",
		// "https://nessie2024.ct.digicert.com/log/",
		// "https://wyvern.ct.digicert.com/2024h2/",
		// "https://sphinx.ct.digicert.com/2024h2/",
		// "https://sabre2024h2.ct.sectigo.com/",
		// "https://mammoth2024h2.ct.sectigo.com/",
		"https://oak.ct.letsencrypt.org/2024h2/",
		// "https://ct2024.trustasia.com/log2024/",
	}, context.Background())

	if err != nil {
		fmt.Printf("Failed to create new ctstream")
		os.Exit(1)
	}

	err = m.Init()
	if err != nil {
		fmt.Printf("Failed to initialize ctstream")
		os.Exit(1)
	}

	m.Start(func(cert *ctX509.Certificate, i int, option any, err error) {
		params := option.(*direct.CTClientParams)

		if err != nil {
			fmt.Printf("Failed to fetch: %v\n", err)
		}

		fmt.Printf("%v, %v, %v\n", params.Start+int64(i), cert.DNSNames, params.LogClient.BaseURI())
	})

	go func() {
		time.Sleep(10 * time.Second)
		m.Stop()
	}()

	m.Await()
}
