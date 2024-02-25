package core

import (
	"fmt"
	"os"

	goutils "github.com/akakou/go-utils"
	golangutils "github.com/akakou/golang-utils"
	metact "github.com/akakou/meta-ct"
	"github.com/akakou/ra_webs/ttp/db"
	simplecertify "github.com/akakou/simple-certify"
)

var ATTEST_PROXY_UNIQUE_ID = []byte{}

type TTP struct {
	DB         *db.DB
	CA         *simplecertify.Certifier
	CT         *metact.MetaCT
	AdminToken string
}

func NewTTP(db *db.DB, ca *simplecertify.Certifier, ct *metact.MetaCT, adminToken string) (*TTP, error) {
	return &TTP{
		DB:         db,
		CA:         ca,
		CT:         ct,
		AdminToken: adminToken,
	}, nil
}

func DefaultTTP() (*TTP, error) {
	dbType := golangutils.GetEnv("DB_TYPE", "sqlite3")
	dbConfig := golangutils.GetEnv("DB_CONFIG", "file:ent?mode=memory&cache=shared&_fk=1")
	fmt.Printf("We use %s as database type and %s as database config\n", dbType, dbConfig)

	metaAppId := os.Getenv("META_APP_ID")
	metaAccessToken := os.Getenv("META_ACCESS_TOKEN")

	adminToken, err := goutils.RandomHex(32)
	if err != nil {
		return nil, fmt.Errorf("failed to generate admin token: %w", err)
	}

	fmt.Printf("Admin token generated: %s\n", adminToken)

	caTempl := simplecertify.CATemplate()

	dbc := db.DBConfig{
		Type:   dbType,
		Config: dbConfig,
	}

	db, err := db.NewDB(&dbc)
	if err != nil {
		return nil, fmt.Errorf("failed to init db: %w", err)
	}

	ct := metact.NewCT(metaAppId, metaAccessToken)

	ca, err := simplecertify.LoadOrInit(&caTempl, &caTempl)

	if err != nil {
		return nil, fmt.Errorf("failed to init ca: %w", err)

	}

	return NewTTP(db, ca, ct, adminToken)

}
