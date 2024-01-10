package ta

import "crypto/tls"

type RAConfig struct {
	TTPDomain string
	Domain    string
}

type RA struct {
	config *RAConfig
}

func NewRA(config *RAConfig) *RA {
	return &RA{
		config: config,
	}
}

func TLSConfig(config *RA) (*tls.Config, error) {
	tlsConfig, _, err := config.generateKeyPair()
	if err != nil {
		return nil, err
	}

	return tlsConfig, nil

}
