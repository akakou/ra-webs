package log

import (
	"os"
)

const DB_FILE = "log.db?_fk=1"
const TOKEN_NAME = "RA_WEBS_LOG_TOKEN"
const DOMAIN_NAME = "RA_WEBS_LOG_DOMAIN"
const AT_PRIVATE_KEY_NAME = "RA_WEBS_AT_PRIVATE_KEY"

type Log struct {
	Domain string
	DB     *DB
	Token  string
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

	return &Log{
		DB:     db,
		Domain: domain,
		Token:  token,
	}, nil
}
