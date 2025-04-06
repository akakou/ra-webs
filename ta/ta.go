package ta

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"

	"golang.org/x/crypto/acme/autocert"
)

var AcmeURL = autocert.DefaultACMEDirectory

type TA struct {
	config     TAConfig
	privateKey *rsa.PrivateKey
}

func Default() (*TA, error) {
	config, err := DefaultConfig()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ERROR_DEFAULT_CONFIG, err)
	}

	return Init(config)
}

func Init(config *TAConfig) (*TA, error) {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ERROR_GENERATE_RSA_KEY, err)
	}

	return &TA{
		config:     *config,
		privateKey: privKey,
	}, nil
}
