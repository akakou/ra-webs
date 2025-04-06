package log

import (
	"crypto/rand"
	"crypto/rsa"
	"os"
)

const DB_FILE = "log.db"
const TOKEN_NAME = "RA_WEBS_LOG_TOKEN"
const DOMAIN_NAME = "RA_WEBS_LOG_DOMAIN"

type Log struct {
	SignKey   *rsa.PrivateKey
	VerifyKey *rsa.PublicKey
	Domain    string
	DB        *DB
	Token     string
}

func Default(privateKey *rsa.PrivateKey) (*Log, error) {
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

	return &Log{
		SignKey:   privateKey,
		VerifyKey: &privateKey.PublicKey,
		DB:        db,
		Domain:    domain,
		Token:     token,
	}, nil
}

func Random(privateKey *rsa.PrivateKey) (*Log, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	def, err := Default(key)
	if err != nil {
		return nil, err
	}

	return def, nil
}
