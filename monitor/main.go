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
	taLog, skip := monitor.MonitorCTLog(ctLog)
	if skip {
		return
	}

	taEntry, err := monitor.ServiceClient.Fetch(*taLog.PublicKey)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	monitor.MonitorEvidence(taEntry, taLog)
}

func (monitor *Monitor) MonitorEvidence(taEntry *serviceclient.EvidenceEntry, taLog *ent.TA) {
	var err error

	report, err := CheckEvidence(taEntry.Evidence)
	if err != nil {
		fmt.Printf("Failed to Check Evidence: %v\n", err)
		return
	}

	if !reflect.DeepEqual(report.Data, *taLog.PublicKey) {
		fmt.Printf("Failed to Check Public Key: %v\n", err)
		return
	}

	uniqueId := report.UniqueID

	err = CheckSourceHash(taEntry, uniqueId)
	if err != nil {
		fmt.Printf("Failed to Check Source Hash: %v\n", err)
		return
	}

	atLog := monitor.RegisterATLog(uniqueId, taEntry, taLog, true)
	fmt.Printf("Inserted: %v\n", atLog)
}

func (monitor *Monitor) MonitorCTLog(entry crtsh.CertificateEntry) (*ent.TA, bool) {
	skip := monitor.DB.Client.CTLog.Query().
		Where(ctlog.MonitorLogIDEQ(entry.ID)).
		ExistX(*monitor.DB.Ctx)

	if skip {
		return nil, skip
	}

	unmarshaledPublicKey, isRSA := entry.Certificate.PublicKey.(*rsa.PublicKey)
	publicKeyBuf := []byte("no public key")

	if isRSA {
		publicKeyBuf = x509.MarshalPKCS1PublicKey(unmarshaledPublicKey)
	}

	ta := monitor.RegisterTA(publicKeyBuf)
	monitor.RegisterCTLog(entry.ID, ta, true)

	if isRSA {
		NotifyUpdateX(monitor)
	} else {
		skip = true
		fmt.Printf("Violation: %v\n", errPublicKeyIsNotRSA)
		NotifyViolation(monitor)
	}

	return ta, skip
}
