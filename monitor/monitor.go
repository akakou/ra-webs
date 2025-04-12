package monitor

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	golangutils "github.com/akakou/golang-utils"
	"github.com/akakou/ra-webs/monitor/logclient"
	"github.com/cockroachdb/errors"
)

type Monitor struct {
	Domain      string
	DB          *DB
	CT          CT
	Notifier    Notifier
	LogClient   *logclient.LogClient
	ATPublicKey *rsa.PublicKey
}

func New(domain string, db *DB, ct CT, atPublicKey *rsa.PublicKey, notifier Notifier, logclient *logclient.LogClient) (*Monitor, error) {
	return &Monitor{
		Domain:      domain,
		DB:          db,
		CT:          ct,
		ATPublicKey: atPublicKey,
		Notifier:    notifier,
		LogClient:   logclient,
	}, nil
}

func Default(ct CT, notifier Notifier) (*Monitor, error) {
	domain := os.Getenv("DOMAIN")
	if domain == "" {
		return nil, errors.Wrap(errDomainEnvironmentVariableIsEmpty, "DOMAIN not found")
	}

	atPublicKey := os.Getenv("RA_WEBS_AT_PUBLIC_KEY")
	pemATPublicKey, _ := pem.Decode([]byte(atPublicKey))
	rsaATPublicKey, err := x509.ParsePKIXPublicKey(pemATPublicKey.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse AT public key: %w", err)
	}

	atDomain := os.Getenv("RA_WEBS_AT_DOMAIN")
	if domain == "" {
		return nil, errors.Wrap(errDomainEnvironmentVariableIsEmpty, "RA_WEBS_AT_DOMAIN not found")
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

	logclient, err := logclient.New(atDomain)
	if err != nil {
		return nil, err
	}

	return New(domain, db, ct, rsaATPublicKey.(*rsa.PublicKey), notifier, logclient)
}

func (monitor *Monitor) Run() {
	monitor.CT.Run(monitor)
}

func (m *Monitor) Close() {
	m.DB.Close()
}
