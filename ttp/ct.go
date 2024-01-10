package ttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
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

const SSLMATE_API_URL = "https://api.certspotter.com/v1/issuances?domain=%v&include_subdomains=true&expand=dns_names&expand=cert_der"

func AuditCTLog(domain string, db *ttpDB) error {
	l := 1

	taInfo, err := db.selectTaInfoByDomain(domain)

	if err != nil {
		return fmt.Errorf("failed to select ta info by domain: %w", err)
	}

	for l != 0 {
		subdomain, _l := subDomain(domain)
		l = _l

		ctLogs, err := fetchCTLogs(subdomain)
		if err != nil {
			return fmt.Errorf("failed to fetch ct logs: %w", err)
		}

		// todo:
		if checkCTLogs(ctLogs, "taInfo.PublicKey") != nil {
			return fmt.Errorf("failed to check ct logs: %w", err)
		}

	}

	return nil
}

func fetchCTLogs(domain string) ([]CTLog, error) {
	ctLog := []CTLog{}

	url := fmt.Sprintf(SSLMATE_API_URL, domain)

	resp, _ := http.Get(url)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	err := json.Unmarshal(body, &ctLog)
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

func subDomain(domain string) (string, int) {
	splited := strings.Split(domain, ".")
	l := len(splited) - 1

	splited = splited[:l]

	return strings.Join(splited, "."), l
}
