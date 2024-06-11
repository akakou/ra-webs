package ct

import (
	"crypto/x509"
	"io"
	"os"

	ct "github.com/google/certificate-transparency-go"
)

func readFile(name string) (string, error) {
	lastFile, err := os.Open(name)
	if err != nil {
		return "", err
	}

	defer lastFile.Close()

	b, err := io.ReadAll(lastFile)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func writeFile(name, body string) error {
	lastFile, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer lastFile.Close()

	_, err = io.WriteString(lastFile, body)
	if err != nil {
		return err
	}

	return nil
}

func removeDeplication(src []string) []string {
	midMap := make(map[string]bool)
	dest := []string{}

	for _, id := range src {
		if !midMap[id] {
			midMap[id] = true
			dest = append(dest, id)
		}
	}

	return dest
}

func parseEntries(entries []ct.LogEntry) ([]x509.Certificate, int64, error) {
	last := int64(0)

	certs := []x509.Certificate{}
	for _, e := range entries {
		cert, err := x509.ParseCertificate(e.X509Cert.Raw)

		if err != nil {
			return nil, 0, err
		}

		certs = append(certs, *cert)
		last = e.Index
	}

	return certs, last, nil
}
