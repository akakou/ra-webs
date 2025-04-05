package monitor

import (
	"fmt"
)

type publicKey interface{}

func (monitor *Monitor) Monitor(ctPublicKey publicKey, id int) {
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

	err = CheckPublicKey(ctPublicKey, publicKey)
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
