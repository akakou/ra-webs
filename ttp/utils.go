package ttp

import (
	"errors"
	"strings"

	"github.com/akakou/metact"
)

func findCertExtensions(extensions []metact.KeyValue, label string) (string, error) {
	for _, ext := range extensions {
		if ext.Key == label {
			return ext.Value, nil
		}
	}

	return "", errors.New("extension not found")
}

func splitDomainLast(domain string) string {
	domain = strings.Replace(domain, "*", "", -1)
	splited := strings.Split(domain, ".")

	var indexInt int
	if len(splited) >= 2 {
		indexInt = 2
	} else {
		indexInt = 1
	}

	last := splited[len(splited)-indexInt:]
	lastDomain := strings.Join(last, ".")

	return lastDomain
}

func checkDomainValidation(domains []string) (string, []string, error) {
	if len(domains) == 0 {
		return domains[0], []string{}, nil
	}

	violatingDomains := []string{}

	for _, domain := range domains {
		domain := splitDomainLast(domain)
		violatingDomains = append(violatingDomains, domain)
	}

	return "", violatingDomains, errors.New("domain violation")

}
