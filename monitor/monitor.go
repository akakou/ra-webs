package monitor

import (
	"fmt"
	"os"
	"time"

	ctcore "github.com/akakou/ctstream/core"
	golangutils "github.com/akakou/golang-utils"
	devkitserviceclient "github.com/akakou/ra-webs/devkit/service/client"
	"github.com/akakou/ra-webs/monitor/serviceclient"
	"github.com/cockroachdb/errors"
)

var Sleep = time.Minute

type Monitor struct {
	TADomain      string
	DB            *DB
	CT            ctcore.CtClient
	Notifier      Notifier
	ServiceClient serviceclient.ServiceClient
}

func New(taDomain string, db *DB, ct ctcore.CtClient, notifier Notifier, serviceClient serviceclient.ServiceClient) (*Monitor, error) {
	return &Monitor{
		TADomain:      taDomain,
		DB:            db,
		CT:            ct,
		Notifier:      notifier,
		ServiceClient: serviceClient,
	}, nil
}

func Default(ct ctcore.CtClient, notifier Notifier) (*Monitor, error) {
	taDomain := os.Getenv("RA_WEBS_TA_DOMAIN")
	if taDomain == "" {
		return nil, errors.Wrap(errDomainEnvironmentVariableIsEmpty, "RA_WEBS_TA_DOMAIN not found")
	}

	serviceBase := os.Getenv("RA_WEBS_SERVICE_BASE")
	if serviceBase == "" {
		return nil, errors.Wrap(errDomainEnvironmentVariableIsEmpty, "RA_WEBS_SERVICE_BASE not found")
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

	serviceClient, err := devkitserviceclient.New(serviceBase)
	if err != nil {
		return nil, err
	}

	return New(taDomain, db, ct, notifier, serviceClient)
}

func (monitor *Monitor) Run() {

	for {
		monitor.CT.Init()
		monitor.CT.Next(monitor.Monitor)
		time.Sleep(Sleep)
	}
}

func (m *Monitor) Close() {
	m.DB.Close()
}
