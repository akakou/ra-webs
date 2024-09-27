package ta

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"os"
	"strings"

	"github.com/akakou/ra_webs/core"
	"golang.org/x/crypto/acme/autocert"
)

var AcmeURL = autocert.DefaultACMEDirectory

type TAConfig struct {
	Verifiers  []string
	Token      string
	Repository string
	Domain     string
	Email      string
}

type TA struct {
	config     TAConfig
	privateKey *rsa.PrivateKey
}

func DefaultConfig() (*TAConfig, error) {
	token := os.Getenv("RA_WEBS_SERVICE_TOKEN")
	repository := os.Getenv("RA_WEBS_TA_REPOSITORY")
	domain := os.Getenv("RA_WEBS_TA_DOMAIN")
	email := os.Getenv("RA_WEBS_SERVICE_EMAIL")
	verifierBaseEnv := os.Getenv("RA_WEBS_VERIFIER_BASES")

	if token == "" {
		return nil, fmt.Errorf("%v", ERROR_TOKEN_NOT_SET)
	}

	if repository == "" {
		return nil, fmt.Errorf("%v", ERROR_REPOSITORY_NOT_SET)
	}

	if email == "" {
		return nil, fmt.Errorf("%v", ERROR_EMAIL_NOT_SET)
	}

	if verifierBaseEnv == "" {
		verifierBaseEnv = "http://localhost" + core.VerifierPort
		fmt.Printf("RA_WEBS_VERIFIER_BASES is not set: so use %v\n", verifierBaseEnv)
	}
	verifiers := strings.Split(verifierBaseEnv, ",")

	return &TAConfig{
		Token:      token,
		Repository: repository,
		Domain:     domain,
		Verifiers:  verifiers,
	}, nil
}

func DefaultTA() (*TA, error) {
	config, err := DefaultConfig()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ERROR_DEFAULT_CONFIG, err)
	}

	return InitTA(config)
}

func InitTA(config *TAConfig) (*TA, error) {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ERROR_GENERATE_RSA_KEY, err)
	}

	return &TA{
		config:     *config,
		privateKey: privKey,
	}, nil
}
