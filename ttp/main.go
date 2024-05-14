package ttp

import (
	"fmt"
	"os"

	goutils "github.com/akakou/go-utils"
	golangutils "github.com/akakou/golang-utils"
	"github.com/akakou/ra_webs/ttp/core"
	"github.com/akakou/ra_webs/ttp/ct"
)

func DefaultTTP() (*core.TTP, error) {
	dbType := golangutils.GetEnv("DB_TYPE", "sqlite3")
	dbConfig := golangutils.GetEnv("DB_CONFIG", "file:ent?mode=memory&cache=shared&_fk=1")
	fmt.Printf("We use %s as database type and %s as database config\n", dbType, dbConfig)

	token := os.Getenv("SSLMATE_TOKEN")

	adminToken, err := goutils.RandomHex(core.RANDOM_SIZE)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", core.ERROR_RANDOM_GENERATE, err)
	}

	fmt.Printf("Admin token generated: %s\n", adminToken)

	dbc := core.DBConfig{
		Type:   dbType,
		Config: dbConfig,
	}

	db, err := core.NewDB(&dbc)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", core.ERROR_INIT_DB, err)
	}

	_ct := ct.NewSSLMateCT(token)

	return core.NewTTP(db, _ct, adminToken)
}
