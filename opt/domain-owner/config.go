package domainowner

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

const CONFIG_PATH = "./config.yaml"

func LoadRecordConfig(zone string) (*DNSRecords, error) {
	records := DNSRecords{}
	recordList := map[string]string{}

	b, err := os.ReadFile(CONFIG_PATH)

	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %v", CONFIG_PATH, err)
	}

	err = yaml.Unmarshal(b, &recordList)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	records.FromDomains(recordList, zone)

	return &records, nil
}

func (s DNSServer) LoadConfig() error {
	records, err := LoadRecordConfig(s.Zone)

	if err != nil {
		return err
	}

	s.Records = records

	return nil
}
