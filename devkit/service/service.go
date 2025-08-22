package service

import (
	"os"
)

const DB_FILE = "service.db?_fk=1"
const TOKEN_NAME = "RA_WEBS_SERVICE_TOKEN"
const TA_DOMAIN_NAME = "RA_WEBS_TA_DOMAIN"

type Service struct {
	Domain string
	DB     *DB
	Token  string
}

func Default() (*Service, error) {
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

	domain := os.Getenv(TA_DOMAIN_NAME)
	if domain == "" {
		return nil, ErrNoDomain
	}

	return &Service{
		DB:     db,
		Domain: domain,
		Token:  token,
	}, nil
}
