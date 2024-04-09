package ta

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"

	"golang.org/x/crypto/acme/autocert"
)

type Config struct {
	ServerId int
	Token    string
	ACMEUrl  string
}

type TA struct {
	Config     Config
	PrivateKey *rsa.PrivateKey
}

func DefaultConfig() (*Config, error) {
	return &Config{
		ACMEUrl: autocert.DefaultACMEDirectory,
	}, nil
}

func DefaultTA() (*TA, error) {
	config, err := DefaultConfig()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ERROR_DEFAULT_CONFIG, err)
	}

	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ERROR_GENERATE_RSA_KEY, err)
	}

	return &TA{
		Config:     *config,
		PrivateKey: privKey,
	}, nil
}
