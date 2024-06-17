package audit

import (
	"fmt"
	"time"

	"github.com/akakou/ctstream"
	"github.com/akakou/ra_webs/ttp/core"
	"github.com/google/certificate-transparency-go/client"
	"github.com/google/certificate-transparency-go/x509"
)

var DefaultSleep = 5 * time.Second
var DefaultCTLogs = []string{
	"https://ct.googleapis.com/logs/us1/argon2024/",
	"https://ct.googleapis.com/logs/us1/argon2025h1/",
	"https://ct.googleapis.com/logs/us1/argon2025h2/",
	"https://ct.googleapis.com/logs/eu1/xenon2024/",
	"https://ct.googleapis.com/logs/eu1/xenon2025h1/",
	"https://ct.googleapis.com/logs/eu1/xenon2025h2/",
	"https://ct.cloudflare.com/logs/nimbus2024/",
	"https://ct.cloudflare.com/logs/nimbus2025/",
	"https://yeti2024.ct.digicert.com/log/",
	"https://yeti2025.ct.digicert.com/log/",
	"https://nessie2024.ct.digicert.com/log/",
	"https://nessie2025.ct.digicert.com/log/",
	"https://sabre.ct.comodo.com/",
	"https://sabre2024h1.ct.sectigo.com/",
	"https://sabre2024h2.ct.sectigo.com/",
	"https://sabre2025h1.ct.sectigo.com/",
	"https://sabre2025h2.ct.sectigo.com/",
	"https://mammoth2024h1.ct.sectigo.com/",
	"https://mammoth2024h1b.ct.sectigo.com/",
	"https://mammoth2024h2.ct.sectigo.com/",
	"https://mammoth2025h1.ct.sectigo.com/",
	"https://mammoth2025h2.ct.sectigo.com/",
	"https://oak.ct.letsencrypt.org/2024h1/",
	"https://oak.ct.letsencrypt.org/2024h2/",
	"https://oak.ct.letsencrypt.org/2025h1/",
	"https://oak.ct.letsencrypt.org/2025h2/",
	"https://ct2024.trustasia.com/log2024/",
	"https://ct2025-a.trustasia.com/log2025a/",
	"https://ct2025-b.trustasia.com/log2025b/",
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
