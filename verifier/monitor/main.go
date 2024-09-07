package monitor

import (
	"context"
	"fmt"

	ctcore "github.com/akakou/ctstream/core"
	"github.com/akakou/ctstream/direct"
	"github.com/akakou/ctstream/thirdparty/sslmate"
	"github.com/akakou/ra_webs/verifier/core"
	ctx509 "github.com/google/certificate-transparency-go/x509"
)

var DefaultCTLogs = []string{
	"https://ct.googleapis.com/logs/us1/argon2024/",
	"https://ct.googleapis.com/logs/eu1/xenon2024/",
	"https://ct.cloudflare.com/logs/nimbus2024/",
	"https://yeti2024.ct.digicert.com/log/",
	"https://nessie2024.ct.digicert.com/log/",
	"https://wyvern.ct.digicert.com/2024h2/",
	"https://sphinx.ct.digicert.com/2024h2/",
	"https://sabre2024h2.ct.sectigo.com/",
	"https://mammoth2024h2.ct.sectigo.com/",
	"https://oak.ct.letsencrypt.org/2024h2/",
	"https://ct2024.trustasia.com/log2024/",
}

type Monitor struct {
	ctstream ctcore.CtStream
}

func NewMonitor[T ctcore.CtStream](stream T) (*Monitor, error) {
	return &Monitor{
		ctstream: stream,
	}, nil
}

func DefaultDirectMonitor() (*Monitor, error) {
	ctx := context.Background()
	stream, err := direct.DefaultCTsStream(DefaultCTLogs, ctx)
	if err != nil {
		return nil, err
	}

	return NewMonitor(stream)
}

func DefaultSSLMateMonitor() (*Monitor, error) {
	ctx := context.Background()
	stream, err := sslmate.DefaultCTsStream(DefaultCTLogs, ctx)
	if err != nil {
		return nil, err
	}

	return NewMonitor(stream)
}

func (a *Monitor) Setup(verifier *core.Verifier) error {
	return a.ctstream.Init()
}

func (a *Monitor) Run(verifier *core.Verifier) {
	a.ctstream.Run(func(cert *ctx509.Certificate, i int, params any, err error) {
		if err == nil {
		} else if err.Error() == direct.ERROR_FAILED_TO_NEW {
			return
		} else {
			fmt.Printf("Error: %v\n", err)
			return
		}

		err = Check(verifier, cert)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}

		// fmt.Printf("Certificate: %v\n", cert.Subject.CommonName)
	})
}
