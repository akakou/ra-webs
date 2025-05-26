package monitor

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"reflect"
	"time"

	"github.com/akakou/ra-webs/monitor/ent"
	"github.com/akakou/ra-webs/monitor/ent/ctlog"
	"github.com/akakou/ra-webs/monitor/serviceclient"
	ctx509 "github.com/google/certificate-transparency-go/x509"
)

func (monitor *Monitor) Monitor(cert *ctx509.Certificate, id int, params any, err error) {
	fmt.Printf("Monitor: %v %v %v\n", cert, id, params)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		time.Sleep(time.Minute * 5)
		return
	}

	taLog, skip, err := monitor.MonitorCTLog(cert, id)
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

	evidenceLog := monitor.RegisterEvidenceLog(uniqueId, taEntry, taLog)
	fmt.Printf("Inserted: %v\n", evidenceLog)

	return nil
}

func (monitor *Monitor) MonitorCTLog(cert *ctx509.Certificate, id int) (*ent.TA, bool, error) {
	var err error
	skip := monitor.DB.Client.CTLog.Query().
		Where(ctlog.MonitorLogIDEQ(id)).
		ExistX(*monitor.DB.Ctx)

	if skip {
		return nil, true, nil
	}

	unmarshaledPublicKey, isRSA := cert.PublicKey.(*rsa.PublicKey)
	publicKeyBuf := []byte("no public key")

	if isRSA {
		publicKeyBuf = x509.MarshalPKCS1PublicKey(unmarshaledPublicKey)
	} else {
		err = fmt.Errorf("Violation: %v\n", errPublicKeyIsNotRSA)
	}

	ta := monitor.RegisterTA(publicKeyBuf)
	monitor.RegisterCTLog(id, ta)

	return ta, false, err
}
