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

type SSLMateStream = ctcore.ConcurrentCTsStream[*ctcore.CTStream[*sslmate.SSLMateCTClient]]

type SSLMateMonitor struct {
	ctstream *SSLMateStream
	ctx      context.Context
}

func (a *SSLMateMonitor) Setup(verifier *core.Verifier) error {
	return a.ctstream.Init()
}

func (a *SSLMateMonitor) Register(domain string, verifier *core.Verifier) error {
	err := sslmate.AddByDomain(domain, context.Background(), a.ctstream)
	return err
}

func (a *SSLMateMonitor) Run(verifier *core.Verifier) {
	a.ctstream.Run(func(cert *ctx509.Certificate, i int, params any, err error) {
		if err == nil {
		} else if err.Error() == direct.ERROR_FAILED_TO_NEW {
			return
		} else {
			fmt.Printf("Error: %v\n", err)
			return
		}

		index := params.(sslmate.SSLMateCTParams)

		err = Check(cert, index.Last, verifier)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}

		// fmt.Printf("Certificate: %v\n", cert.Subject.CommonName)
	})
}
