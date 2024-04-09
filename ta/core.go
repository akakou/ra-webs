package ta

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"os"
	"strconv"

	"golang.org/x/crypto/acme/autocert"
)

const TOKEN_NAME = "RA_WEBS_TOKEN"
const SERVER_ID_NAME = "RA_WEBS_SERVER_ID"

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
	token := os.Getenv(TOKEN_NAME)
	if token == "" {
		return nil, fmt.Errorf("token not found")
	}

	s := os.Getenv(SERVER_ID_NAME)
	serverId, err := strconv.Atoi(s)
	if err != nil {
		return nil, fmt.Errorf("failed to parse server id: %w", err)
	}

	return &Config{
		ServerId: serverId,
		Token:    token,
		ACMEUrl:  autocert.DefaultACMEDirectory,
	}, nil
}

func DefaultTA() (*TA, error) {
	config, err := DefaultConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("failed to generate rsa key: %w", err)
	}

	return &TA{
		Config:     *config,
		PrivateKey: privKey,
	}, nil
}
