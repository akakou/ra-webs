package ttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type CTLogAudit struct {
	TADomain   string
	IsValid    bool
	LatestCTId string
}

type CTLog struct {
	Id           string `json:"id"`
	PubKeySha256 string `json:"pubkey_sha256"`
}

const SSLMATE_API_URL = "https://api.certspotter.com/v1/issuances?domain=%v&match_wildcards=true&expand=dns_names&expand=cert_der"

func AuditCTLog(domain string, db *ttpDB) error {
	taInfo, err := db.selectTaInfoByDomain(domain)
	if err != nil {
		return fmt.Errorf("failed to get ta info: %w", err)
	}

	ctLog := taInfo.Edges.CtLog
	if !ctLog.IsValid {
		return errors.New("ct log is not valid")
	}

	ctLogs, err := fetchCTLogs(domain, ctLog.LatestCtID)
	if err != nil {
		return fmt.Errorf("failed to fetch ct logs: %w", err)
	}

	logLen := len(ctLogs)
	ctLog.LatestCtID = ctLogs[logLen-1].Id

	if checkCTLogs(ctLogs, taInfo.PublicKeyHash) != nil {
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

func checkCTLogs(ctLogs []CTLog, hashedPublicKey string) error {
	for _, ctLog := range ctLogs {
		if ctLog.PubKeySha256 != hashedPublicKey {
			return errors.New("public key not match")
		}
	}

	return nil
}
