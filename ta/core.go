package ta

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
)

type TAConfig struct {
	ServerId int
	CodeId   int
	Token    string
	Type     REGISTER_TYPE
}

type TA struct {
	Config     TAConfig
	PrivateKey *rsa.PrivateKey
}

func DefaultTA(config TAConfig) (*TA, error) {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		return nil, fmt.Errorf("failed to generate rsa key: %w", err)
	}

	return &TA{
		Config:     config,
		PrivateKey: privKey,
	}, nil
}
