package ta

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"

	"github.com/labstack/echo/v4"
)

type AttestProxyConfig struct {
	ServerId int
	CodeId   int
	Token    string
}

type AttestProxy struct {
	Config     AttestProxyConfig
	Echo       *echo.Echo
	PrivateKey *rsa.PrivateKey
}

func DefaultAttestProxy(config AttestProxyConfig) (*AttestProxy, error) {
	echo := echo.New()

	privKey, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		return nil, fmt.Errorf("failed to generate rsa key: %w", err)
	}

	return &AttestProxy{
		Config:     config,
		PrivateKey: privKey,
		Echo:       echo,
	}, nil
}
