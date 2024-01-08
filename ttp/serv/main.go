package main

import (
	"flag"

	"github.com/akakou/ra_webs/ttp"
)

const PORT = ":1323"

func main() {
	dbType := flag.String("db_type", "sqlite3", "database type")
	dbConfig := flag.String("db_config", "file:ent?mode=memory&cache=shared&_fk=1", "database config")

	DBConfig := ttp.DBConfig{
		Type:   *dbType,
		Config: *dbConfig,
	}

	e := ttp.NewTTPServer(&DBConfig)
	e.Logger.Fatal(e.Start(PORT))
}
