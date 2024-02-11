package ttp

import (
	"errors"
	"strings"

	metact "github.com/akakou/meta-ct"
)

func findCertExtensions(extensions []metact.KeyValue, label string) (string, error) {
	for _, ext := range extensions {
		if ext.Key == label {
			return ext.Value, nil
		}
	}

	return "", errors.New("extension not found")
}

func extractDomainLast(domain string) string {
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
