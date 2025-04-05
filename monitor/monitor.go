package monitor

import (
	"fmt"

	golangutils "github.com/akakou/golang-utils"
)

type Monitor struct {
	Domain   string
	DB       *DB
	CT       CTMonitor
	Notifier Notifier
}

func New(db *DB) (*Monitor, error) {
	return &Monitor{
		DB: db,
	}, nil
}

func Default() (*Monitor, error) {
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

	return New(db)
}
