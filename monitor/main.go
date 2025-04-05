package monitor

import (
	"fmt"

	"github.com/akakou/crtsh"
	"github.com/akakou/ra-webs/monitor/ent/taserver"
)

type publicKey interface{}

func (monitor *Monitor) MonitorAll(entries []crtsh.CertificateEntry) {
	for _, entry := range entries {
		monitor.MonitorOne(entry)
	}
}

func (monitor *Monitor) MonitorOne(entry crtsh.CertificateEntry) {
	exist, err := monitor.DB.Client.TAServer.Query().
		Where(taserver.MonitorLogIDEQ(entry.ID)).
		Exist(*monitor.DB.Ctx)

	if err != nil {
		panic(err)
	}

	if exist {
		return
	}

	revoked := false

	log, err := monitor.LogClient.Fetch()
	if err != nil {
		fmt.Printf("Violation: %v\n", err)
		revoked = true
	}

	report, err := CheckEvidence(string(log.Evidence))
	if err != nil {
		fmt.Printf("Violation: %v\n", err)
		revoked = true
	}

	publicKey := report.Data
	uniqueId := report.Data

	err = CheckSourceHash(log, uniqueId)
	if err != nil {
		fmt.Printf("Violation: %v\n", err)
		revoked = true
	}

	err = CheckPublicKey(entry.Certificate.PublicKey, publicKey)
	if err != nil {
		fmt.Printf("Violation: %v\n", err)
		revoked = true
	}

	code, err := monitor.RegisterCode(uniqueId, log)
	if err != nil {
		panic(err)
	}

	server, err := monitor.RegisterServer(string(log.Evidence), publicKey, code)
	if err != nil {
		panic(err)
	}

	if revoked {
		monitor.Revoke(server)
	} else {
		err = NotifyUpdate(monitor.Domain, monitor)
		if err != nil {
			panic(err)
		}
	}
}
