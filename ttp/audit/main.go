package audit

import (
	"context"
	"fmt"

	"github.com/akakou/ctstream"
	"github.com/akakou/ra_webs/ttp/core"
	"github.com/google/certificate-transparency-go/client"
	"github.com/google/certificate-transparency-go/x509"
)

var DefaultCTLogs = []string{
	"https://oak.ct.letsencrypt.org/2024h1/",
	"https://oak.ct.letsencrypt.org/2024h2/",
}

type Auditor struct {
	ctstream *ctstream.CTsStream
}

func NewAuditor(url []string, ctx context.Context) (*Auditor, error) {
	stream, err := ctstream.DefaultCTsStream(url, ctx)
	if err != nil {
		return nil, err
	}

	return &Auditor{
		ctstream: stream,
	}, nil
}

func DefaultAuditor() (*Auditor, error) {
	ctx := context.Background()
	return NewAuditor(DefaultCTLogs, ctx)
}

func (a *Auditor) Setup(ttp *core.TTP) error {
	return a.ctstream.Init()
}

func (a *Auditor) Run(ttp *core.TTP) {
	a.ctstream.Run(func(cert *x509.Certificate, li ctstream.LogID, lc *client.LogClient, err error) {
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}

		err = Audit(ttp, cert)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}

		fmt.Printf("Certificate: %v\n", cert.Subject.CommonName)
	})
}
