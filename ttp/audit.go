package ttp

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/akakou/ra_webs/core"
	"github.com/akakou/ra_webs/ttp/ent/tainfo"
)

type CTLogAudit struct {
	TADomain   string
	IsValid    bool
	LatestCTId string
}

type CTLog struct {
	Id          string `json:"id"`
	Certificate string `json:"certificate_pem"`
}

func AuditCTLog(domain string, db *ttpDB) error {
	taInfo, err := db.client.TAInfo.
		Query().
		Where(tainfo.DomainEQ(domain)).
		WithCtLog().
		Only(*db.ctx)

	if err != nil {
		return fmt.Errorf("failed to get ta info: %w", err)
	}

	ctLog := taInfo.Edges.CtLog
	if !ctLog.IsValid {
		return errors.New("ct log is not valid")
	}

	taCode := taInfo.Edges.TaCode[len(taInfo.Edges.TaCode)-1]

	ctLogs, err := fetchCTLogs(domain, ctLog.LatestCtID)
	if err != nil {
		return fmt.Errorf("failed to fetch ct logs: %w", err)
	}

	logLen := len(ctLogs)
	ctLog.LatestCtID = ctLogs[logLen-1].Id

	if checkCTLogs(ctLogs, taCode.ProductID) != nil {
		ctLog.IsValid = false
		ctLog.Update().Save(*db.ctx)
		return fmt.Errorf("failed to check ct logs: %w", err)
	}

	ctLog.Update().Save(*db.ctx)

	return nil
}

func fetchCTLogs(domain string, after string) ([]CTLog, error) {
	ctLog := []CTLog{}

	url := fmt.Sprintf(SSLMATE_API_URL, domain)

	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.New("failed to get ct logs")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get ct logs")
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	err = json.Unmarshal(body, &ctLog)
	if err != nil {
		return nil, errors.New("failed to unmarshal json")
	}

	return ctLog, nil
}

func checkCTLogs(ctLogs []CTLog, productId uint16) error {
	for _, ctLog := range ctLogs {
		err := checkCTLog(ctLog, productId)

		if err != nil {
			return err
		}
	}
	return nil
}

func checkCTLog(ctLog CTLog, productId uint16) error {
	cert, err := x509.ParseCertificate([]byte(ctLog.Certificate))
	if err != nil {
		return fmt.Errorf("failed to parse certificate: %v", err)
	}

	extension, err := findCertExtensions(cert.Extensions, core.X509_EXTENSION_LABEL)
	if err != nil {
		return fmt.Errorf("extension not found: %v", err)
	}

	publicKey := cert.PublicKey.(*rsa.PublicKey)

	_, err = core.VerifyByAzure(string(extension.Value), productId, publicKey)
	if err != nil {
		return fmt.Errorf("failed to verify attestation: %v", err)
	}
	return nil
}
