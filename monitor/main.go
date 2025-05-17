package monitor

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"reflect"

	"github.com/akakou/crtsh"
	"github.com/akakou/ra-webs/monitor/ent"
	"github.com/akakou/ra-webs/monitor/ent/ctlog"
	"github.com/akakou/ra-webs/monitor/serviceclient"
)

func (monitor *Monitor) Monitor(ctLogs []crtsh.CertificateEntry) {
	for _, entry := range ctLogs {
		monitor.MonitorOne(entry)
	}
}

func (monitor *Monitor) MonitorOne(ctLog crtsh.CertificateEntry) {
	taLog, skip, err := monitor.MonitorCTLog(ctLog)
	if skip {
		return
	}

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		NotifyViolationX(monitor)
		return
	}

	taEntry, err := monitor.ServiceClient.Fetch(*taLog.PublicKey)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		NotifyViolationX(monitor)
		return
	}

	err = monitor.MonitorEvidence(taEntry, taLog)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		NotifyViolationX(monitor)
		return
	}

	NotifyUpdateX(monitor)
}

func (monitor *Monitor) MonitorEvidence(taEntry *serviceclient.EvidenceEntry, taLog *ent.TA) error {
	var err error

	report, err := CheckEvidence(taEntry.Evidence)
	if err != nil {
		return fmt.Errorf("Failed to Check Evidence: %v\n", err)
	}

	if !reflect.DeepEqual(report.Data, *taLog.PublicKey) {
		return fmt.Errorf("Failed to Check Public Key: %v\n", err)
	}

	uniqueId := report.UniqueID

	err = CheckSourceHash(taEntry, uniqueId)
	if err != nil {
		return fmt.Errorf("Failed to Check Source Hash: %v\n", err)
	}

	evidenceLog := monitor.RegisterATLog(uniqueId, taEntry, taLog)
	fmt.Printf("Inserted: %v\n", evidenceLog)

	return nil
}

func (monitor *Monitor) MonitorCTLog(entry crtsh.CertificateEntry) (*ent.TA, bool, error) {
	var err error

	skip := monitor.DB.Client.CTLog.Query().
		Where(ctlog.MonitorLogIDEQ(entry.ID)).
		ExistX(*monitor.DB.Ctx)

	if skip {
		return nil, true, nil
	}

	unmarshaledPublicKey, isRSA := entry.Certificate.PublicKey.(*rsa.PublicKey)
	publicKeyBuf := []byte("no public key")

	if isRSA {
		publicKeyBuf = x509.MarshalPKCS1PublicKey(unmarshaledPublicKey)
	} else {
		err = fmt.Errorf("Violation: %v\n", errPublicKeyIsNotRSA)
	}

	ta := monitor.RegisterTA(publicKeyBuf)
	monitor.RegisterCTLog(entry.ID, ta)

	return ta, false, err
}
