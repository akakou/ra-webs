package ttp

import (
	"crypto/rand"
	"encoding/hex"
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

func randomHexString(size int) string {
	buf := make([]byte, size)
	// then we can call rand.Read.
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}
	r := hex.EncodeToString(buf)

	return r
}
