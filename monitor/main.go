package monitor

import (
	"fmt"

	"github.com/akakou/crtsh"
	"github.com/akakou/ra-webs/monitor/ent/ctlog"
)

type publicKey interface{}

func (monitor *Monitor) MonitorAll(entries []crtsh.CertificateEntry) {
	for _, entry := range entries {
		monitor.MonitorOne(entry)
	}
}

func (monitor *Monitor) MonitorOne(entry crtsh.CertificateEntry) {
	exist, err := monitor.DB.Client.CTLog.Query().
		Where(ctlog.MonitorLogIDEQ(entry.ID)).
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

	report, err := CheckEvidence(log.Evidence)
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

	ctLog, err := monitor.RegisterCTLog(log.Evidence, publicKey)
	if err != nil {
		panic(err)
	}

	_, err = monitor.RegisterATLog(uniqueId, log, ctLog)
	if err != nil {
		panic(err)
	}

	if revoked {
		monitor.Revoke(ctLog)
	} else {
		err = NotifyUpdate(monitor.Domain, monitor)
		if err != nil {
			panic(err)
		}
	}
}
