package monitor

import (
	"fmt"
	"os"

	golangutils "github.com/akakou/golang-utils"
	devkitserviceclient "github.com/akakou/ra-webs/devkit/service/client"
	"github.com/akakou/ra-webs/monitor/serviceclient"
	"github.com/cockroachdb/errors"
)

type Monitor struct {
	TADomain      string
	DB            *DB
	CT            CT
	Notifier      Notifier
	ServiceClient serviceclient.ServiceClient
}

func New(taDomain string, db *DB, ct CT, notifier Notifier, serviceClient serviceclient.ServiceClient) (*Monitor, error) {
	return &Monitor{
		TADomain:      taDomain,
		DB:            db,
		CT:            ct,
		Notifier:      notifier,
		ServiceClient: serviceClient,
	}, nil
}

func Default(ct CT, notifier Notifier) (*Monitor, error) {
	taDomain := os.Getenv("RA_WEBS_TA_DOMAIN")
	if taDomain == "" {
		return nil, errors.Wrap(errDomainEnvironmentVariableIsEmpty, "RA_WEBS_TA_DOMAIN not found")
	}

	serviceDomain := os.Getenv("RA_WEBS_SERVICE_DOMAIN")
	if serviceDomain == "" {
		return nil, errors.Wrap(errDomainEnvironmentVariableIsEmpty, "RA_WEBS_SERVICE_DOMAIN not found")
	}

	dbType := golangutils.GetEnv("DB_TYPE", "sqlite3")
	dbConfig := golangutils.GetEnv("RA_WEBS_DB_CONFIG", "file:ent?mode=memory&cache=shared&_fk=1")
	fmt.Printf("We use %s as database type and %s as database config\n", dbType, dbConfig)

	dbc := DBConfig{
		Type:   dbType,
		Config: dbConfig,
	}

	db, err := NewDB(&dbc)
	if err != nil {
		return nil, err
	}

	serviceClient, err := devkitserviceclient.New(serviceDomain)
	if err != nil {
		return nil, err
	}

	return New(taDomain, db, ct, notifier, serviceClient)
}

func (monitor *Monitor) Run() {
	monitor.CT.Run(monitor)
}

func (m *Monitor) Close() {
	m.DB.Close()
}
