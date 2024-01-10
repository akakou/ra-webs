package ttp

import (
	"testing"

	"github.com/akakou/ra_webs/core"
	"github.com/go-playground/assert/v2"
)

func makeTestDB() *ttpDB {
	dbConfig := DBConfig{
		Type:   "sqlite3",
		Config: "file:ent?mode=memory&cache=shared&_fk=1",
	}

	db, err := newTtpDB(&dbConfig)

	if err != nil {
		panic(err)
	}

	return db

}

func TestTAInfoDB(t *testing.T) {
	TADomain := "ta.example.com"

	expected := core.TAInfo{
		Domain:        "ta.example.com",
		PublicKeyHash: "public key hash",
		Attestation:   "attestation",
	}

	db := makeTestDB()

	e := db.toEntTaInfo(&expected)
	e.SaveX(*db.ctx)

	taInfo, err := db.selectTaInfoByDomain(TADomain)

	if err != nil {
		panic(err)
	}

	actual := db.toCoreTaInfo(taInfo)

	assert.Equal(t, expected, *actual)
}

func TestCTLogAuditDB(t *testing.T) {
	TADomain := "ta.example.com"

	expected := CTLogAudit{
		IsValid:    true,
		LatestCTId: "",
		TADomain:   TADomain,
	}

	taInfo := core.TAInfo{
		Attestation:   "",
		PublicKeyHash: "hash",
		Domain:        TADomain,
	}

	db := makeTestDB()

	entTaInfoCreate := db.toEntTaInfo(&taInfo)
	entTaInfo := entTaInfoCreate.SaveX(*db.ctx)

	entCTLogAuditCreate, err := db.toEntCTLogAudit(&expected, entTaInfo)
	if err != nil {
		panic(err)
	}

	entCTLogAuditCreate.SaveX(*db.ctx)
	ctLogAudit, err := db.selectCTLogByDomain(TADomain)
	if err != nil {
		panic(err)
	}

	actual := db.toCoreCTLogAudit(ctLogAudit)
	assert.Equal(t, expected, actual)
}
