package verifier

import (
	"context"
	"fmt"

	goutils "github.com/akakou/go-utils"
	golangutils "github.com/akakou/golang-utils"
	"github.com/akakou/ra_webs/verifier/core"
	"github.com/akakou/ra_webs/verifier/monitor"
	notifier "github.com/akakou/ra_webs/verifier/notifier"
)

func DefaultVerifier() (*core.Verifier, error) {
	dbType := golangutils.GetEnv("DB_TYPE", "sqlite3")
	dbConfig := golangutils.GetEnv("DB_CONFIG", "file:ent?mode=memory&cache=shared&_fk=1")
	fmt.Printf("We use %s as database type and %s as database config\n", dbType, dbConfig)

	adminToken, err := goutils.RandomHex(core.RANDOM_SIZE)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", core.ERROR_RANDOM_GENERATE, err)
	}
	adminToken = golangutils.GetEnv("ADMIN_TOKEN", adminToken)

	fmt.Printf("Admin token generated: %s\n", adminToken)

	dbc := core.DBConfig{
		Type:   dbType,
		Config: dbConfig,
	}

	db, err := core.NewDB(&dbc)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", core.ERROR_INIT_DB, err)
	}

	monitor := monitor.DefaultSSLMateMonitor(context.Background())

	notifier, err := notifier.DefaultBrowserNotifier()
	if err != nil {
		return nil, err
	}

	return core.NewVerifier(db, monitor, notifier, adminToken)
}
