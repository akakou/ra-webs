package core

import (
	"fmt"
	"os"

	goutils "github.com/akakou/go-utils"
	golangutils "github.com/akakou/golang-utils"
	metact "github.com/akakou/meta-ct"
	"github.com/akakou/ra_webs/ttp/db"
)

type TTP struct {
	DB         *db.DB
	CT         *metact.MetaCT
	AdminToken string
}

func NewTTP(db *db.DB, ct *metact.MetaCT, adminToken string) (*TTP, error) {
	return &TTP{
		DB:         db,
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
		return nil, fmt.Errorf("%s: %w", ERROR_RANDOM_GENERATE, err)
	}

	fmt.Printf("Admin token generated: %s\n", adminToken)

	dbc := db.DBConfig{
		Type:   dbType,
		Config: dbConfig,
	}

	db, err := db.NewDB(&dbc)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ERROR_INIT_DB, err)
	}

	ct := metact.NewCT(metaAppId, metaAccessToken)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", ERROR_INIT_CT, err)

	}

	return NewTTP(db, ct, adminToken)

}
