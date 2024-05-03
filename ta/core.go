package ta

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"os"

	"github.com/akakou/ra_webs/core"
	"golang.org/x/crypto/acme/autocert"
)

var AcmeURL = autocert.DefaultACMEDirectory

type Config struct {
	ServerId   int
	TTP        string
	Token      string
	Repository string
	Domain     string
}

type TA struct {
	config     Config
	privateKey *rsa.PrivateKey
}

func DefaultConfig() (*Config, error) {
	token := os.Getenv("RA_WEBS_SERVICE_TOKEN")
	repository := os.Getenv("RA_WEBS_REPOSITORY")
	domain := os.Getenv("RA_WEBS_DOMAIN")

	if token == "" {
		panic("RA_WEBS_SERVICE_TOKEN is not set")
	}

	if token == "" {
		panic("Reporsitory is not set")
	}

	return &Config{
		Token:      token,
		Repository: repository,
		Domain:     domain,
		TTP:        "http://localhost" + core.TTPPort,
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
		config:     *config,
		privateKey: privKey,
	}, nil
}
