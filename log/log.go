package log

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

const DB_FILE = "log.db?_fk=1"
const TOKEN_NAME = "RA_WEBS_LOG_TOKEN"
const DOMAIN_NAME = "RA_WEBS_LOG_DOMAIN"
const AT_PRIVATE_KEY_NAME = "RA_WEBS_AT_PRIVATE_KEY"

type Log struct {
	SignKey   *rsa.PrivateKey
	VerifyKey *rsa.PublicKey
	Domain    string
	DB        *DB
	Token     string
}

func Default() (*Log, error) {
	db, err := NewDB(&DBConfig{
		Type:   "sqlite3",
		Config: DB_FILE,
	})

	if err != nil {
		return nil, err
	}

	token := os.Getenv(TOKEN_NAME)
	if token == "" {
		return nil, ErrNoToken
	}

	domain := os.Getenv(DOMAIN_NAME)
	if domain == "" {
		return nil, ErrNoDomain
	}

	privateKey := os.Getenv(AT_PRIVATE_KEY_NAME)

	if privateKey == "" {
		return nil, ErrNoDomain
	}

	pemPrivateKey, _ := pem.Decode([]byte(privateKey))
	rsaPrivateKey, err := x509.ParsePKCS1PrivateKey(pemPrivateKey.Bytes)
	if err != nil {
		return nil, err
	}

	return &Log{
		SignKey:   rsaPrivateKey,
		VerifyKey: &rsaPrivateKey.PublicKey,
		DB:        db,
		Domain:    domain,
		Token:     token,
	}, nil
}
