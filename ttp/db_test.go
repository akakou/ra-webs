package ttp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeTestDB(t *testing.T) *auditDB {
	dbConfig := DBConfig{
		Type:   "sqlite3",
		Config: "file:ent?mode=memory&cache=shared&_fk=1",
	}

	db, err := newAuditDB(&dbConfig)

	if err != nil {
		panic(err)
	}

	assert.NoError(t, err)

	return db

}

func TestTAInfoDB(t *testing.T) {
	TADomain := "ta.example.com"

	expected := TAInfo{
		Domain:        "ta.example.com",
		PublicKeyHash: "public key hash",
		Attestation:   "attestation",
	}

	db := makeTestDB(t)
	defer db.close()

	e := db.toEntTaInfo(&expected)
	e.SaveX(*db.ctx)

	taInfo, err := db.selectTaInfoByDomain(TADomain)

	assert.NoError(t, err)

	actual := db.toCoreTaInfo(taInfo)

	assert.Equal(t, expected, *actual)
}

// func TestCTLogAuditDB(t *testing.T) {
// 	TADomain := "ta.example.com"

// 	expected := CTLogAudit{
// 		IsValid:    true,
// 		LatestCTId: "",
// 		TADomain:   TADomain,
// 	}

// 	taInfo := TAInfo{
// 		Attestation:   "",
// 		PublicKeyHash: "hash",
// 		Domain:        TADomain,
// 	}

// 	db := makeTestDB(t)
// 	defer db.close()

// 	entTaInfoCreate := db.toEntTaInfo(&taInfo)
// 	entTaInfo := entTaInfoCreate.SaveX(*db.ctx)

// 	entCTLogAuditCreate, err := db.toEntCTLogAudit(&expected, entTaInfo)
// 	assert.NoError(t, err)

// 	entCTLogAuditCreate.SaveX(*db.ctx)
// 	ctLogAudit, err := db.selectCTLogByDomain(TADomain)
// 	assert.NoError(t, err)

// 	actual := db.toCoreCTLogAudit(ctLogAudit)
// 	assert.Equal(t, expected, *actual)
// }
