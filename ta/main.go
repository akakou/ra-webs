package ta

import "crypto/tls"

type RAConfig struct {
	TTPDomain string
	Domain    string
}


func TLSConfig(config RAConfig) (*tls.Config, error) {
	tlsConfig, _, err := config.generateKeyPair()
	if err != nil {
		return nil, err
	}

	return tlsConfig, nil

}
