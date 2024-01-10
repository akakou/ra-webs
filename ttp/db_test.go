package ttp

import (
	"testing"

	"github.com/akakou/ra_webs/core"
	"github.com/go-playground/assert/v2"
)

func TestDB(t *testing.T) {
	expected := core.TAInfo{
		Domain:      "ta.example.com",
		PublicKey:   []byte("public key"),
		Attestation: "attestation",
	}

	dbConfig := DBConfig{
		Type:   "sqlite3",
		Config: "file:ent?mode=memory&cache=shared&_fk=1",
	}

	db, err := newttpDB(&dbConfig)

	if err != nil {
		panic(err)
	}

	err = db.store(&expected)

	if err != nil {
		panic(err)
	}

	actual, err := db.findByDomain("ta.example.com")

	if err != nil {
		panic(err)
	}

	assert.Equal(t, expected, *actual)
}
