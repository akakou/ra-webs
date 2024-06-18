package audit

import (
	"fmt"
	"time"

	"github.com/akakou/ctstream"
	"github.com/akakou/ra_webs/ttp/core"
	"github.com/google/certificate-transparency-go/client"
	"github.com/google/certificate-transparency-go/x509"
)

var DefaultSleep = 1 * time.Second
var DefaultCTLogs = []string{
	"https://oak.ct.letsencrypt.org/2024h1/",
	"https://oak.ct.letsencrypt.org/2024h2/",
}

type Auditor struct {
	ctstream *ctstream.CTStream
}

func NewAuditor(sleep time.Duration, url []string) (*Auditor, error) {
	stream, err := ctstream.New(url, sleep)
	if err != nil {
		return nil, err
	}

	return &Auditor{
		ctstream: stream,
	}, nil
}

func DefaultAuditor() (*Auditor, error) {
	return NewAuditor(DefaultSleep, DefaultCTLogs)
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
