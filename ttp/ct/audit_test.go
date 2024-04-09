package ct

import (
	"crypto/x509"
	"testing"

	golangutils "github.com/akakou/golang-utils"
	"github.com/akakou/ra_webs/ttp/core"
	"github.com/akakou/ra_webs/ttp/db"
	"github.com/edgelesssys/ego/attestation"
	"github.com/stretchr/testify/assert"
)

func exampleTTP(t *testing.T) *core.TTP {
	dbType := golangutils.GetEnv("DB_TYPE", "sqlite3")
	dbConfig := golangutils.GetEnv("DB_CONFIG", "file:ent?mode=memory&cache=shared&_fk=1")

	dbc := db.DBConfig{
		Type:   dbType,
		Config: dbConfig,
	}

	db, err := db.NewDB(&dbc)
	assert.NoError(t, err)

	ttp, err := core.NewTTP(db, nil, "")
	assert.NoError(t, err)

	return ttp
}

func TestAll(t *testing.T) {
	t.Run("Pass", testPass)
	t.Run("FailTANoCode", testFailTANoCode)
	t.Run("FailTANoServer", testFailTANoServer)
	t.Run("FailByMissDomains", testFailByMissDomains)
}

func testPass(t *testing.T) {
	ttp := exampleTTP(t)
	defer ttp.DB.Close()

	ttp.DB.Client.TAServer.Create().SetDomain("example.com").SaveX(*ttp.DB.Ctx)
	ttp.DB.Client.TACode.Create().SetUniqueID([]byte{1, 2, 3}).SetRepository("").SetCommitID("").SaveX(*ttp.DB.Ctx)

	validateAttestation = func(_ *x509.Certificate) (*attestation.Report, error) {
		return &attestation.Report{
			UniqueID: []byte{1, 2, 3},
			Data:     []byte{4, 5, 6},
		}, nil
	}

	err := AuditOne(ttp, &x509.Certificate{
		DNSNames:  []string{"example.com"},
		PublicKey: []byte{7, 8, 9},
	})

	assert.NoError(t, err)
}

func testFailTANoCode(t *testing.T) {
	ttp := exampleTTP(t)
	defer ttp.DB.Close()

	ttp.DB.Client.TAServer.Create().SetDomain("example.com").SaveX(*ttp.DB.Ctx)
	ttp.DB.Client.TACode.Create().SetUniqueID([]byte{7, 8, 9}).SetRepository("").SetCommitID("").SaveX(*ttp.DB.Ctx)

	validateAttestation = func(_ *x509.Certificate) (*attestation.Report, error) {
		return &attestation.Report{
			UniqueID: []byte{1, 2, 3},
			Data:     []byte{4, 5, 6},
		}, nil
	}

	err := AuditOne(ttp, &x509.Certificate{
		DNSNames:  []string{"example.com"},
		PublicKey: []byte{7, 8, 9},
	})

	assert.Error(t, err)
}

func testFailTANoServer(t *testing.T) {
	ttp := exampleTTP(t)
	defer ttp.DB.Close()

	ttp.DB.Client.TAServer.Create().SetDomain("example.com").SaveX(*ttp.DB.Ctx)
	ttp.DB.Client.TACode.Create().SetUniqueID([]byte{1, 2, 3}).SetRepository("").SetCommitID("").SaveX(*ttp.DB.Ctx)

	validateAttestation = func(_ *x509.Certificate) (*attestation.Report, error) {
		return &attestation.Report{
			UniqueID: []byte{1, 2, 3},
			Data:     []byte{4, 5, 6},
		}, nil
	}

	err := AuditOne(ttp, &x509.Certificate{
		DNSNames:  []string{"hoge.example.com"},
		PublicKey: []byte{7, 8, 9},
	})

	assert.Error(t, err)
}

func testFailByMissDomains(t *testing.T) {
	// TestMultipleDomain tests the revokeTAByDomains function
	// by revoking multiple TAs by their domain names.

	ttp := exampleTTP(t)
	defer ttp.DB.Close()

	cert := x509.Certificate{
		DNSNames: []string{"example.com", "example.org"},
	}

	err := AuditOne(ttp, &cert)
	assert.Error(t, err)

	cert = x509.Certificate{
		DNSNames:  []string{"*.com"},
		PublicKey: []byte{7, 8, 9},
	}

	err = AuditOne(ttp, &cert)
	assert.Error(t, err)
}
