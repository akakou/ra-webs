package monitor

import (
	"fmt"
	"os"

	golangutils "github.com/akakou/golang-utils"
	"github.com/akakou/ra-webs/monitor/logclient"
)

type Monitor struct {
	Domain    string
	DB        *DB
	CT        CT
	Notifier  Notifier
	LogClient *logclient.LogClient
}

func New(domain string, db *DB, ct CT, notifier Notifier) (*Monitor, error) {
	return &Monitor{
		Domain:   domain,
		DB:       db,
		CT:       ct,
		Notifier: notifier,
	}, nil
}

func Default(ct CT, notifier Notifier) (*Monitor, error) {
	domain := os.Getenv("DOMAIN")
	if domain == "" {
		return nil, errDomainEnvironmentVariableIsEmpty
	}

	dbType := golangutils.GetEnv("DB_TYPE", "sqlite3")
	dbConfig := golangutils.GetEnv("DB_CONFIG", "file:ent?mode=memory&cache=shared&_fk=1")
	fmt.Printf("We use %s as database type and %s as database config\n", dbType, dbConfig)

	dbc := DBConfig{
		Type:   dbType,
		Config: dbConfig,
	}

	db, err := NewDB(&dbc)
	if err != nil {
		return nil, err
	}

	return New(domain, db, ct, notifier)
}

func (monitor *Monitor) Run() {
	monitor.CT.Run(monitor)
}

func (m *Monitor) Close() {
	m.DB.Close()
}
