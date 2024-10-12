package monitor

import (
	"crypto/x509"

	api "github.com/akakou/sslmate_cert_search_api"
)

type Callback func([]x509.Certificate, error)

func (monitor *Monitor) Next() ([]x509.Certificate, *api.Index, error) {
	result := []x509.Certificate{}

	certs, index, err := monitor.Api.Search(monitor.Query)
	if err != nil {
		return nil, index, err
	}

	result = append(result, certs...)

	return result, index, nil
}

func (monitor *Monitor) run(callback Callback) {
	certs, _, err := monitor.Next()
	callback(certs, err)
}

func (monitor *Monitor) Run(callback Callback) {
	for {
		monitor.run(callback)
	}
}

func (monitors Monitors) Run(callback Callback) {
	for _, monitor := range monitors {
		monitor.Run(callback)
	}
}
