package ta

import "crypto/tls"

type RAConfig struct {
	TTPDomain string
	Domain    string
}

func TLSConfig(config RAConfig) (*tls.Config, error) {
	return config.generateKeyPair()
}
