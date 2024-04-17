package ct

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"testing"

	golangutils "github.com/akakou/golang-utils"
	rawebscore "github.com/akakou/ra_webs/core"
	"github.com/akakou/ra_webs/ttp/core"
	"github.com/stretchr/testify/assert"
)

func exampleTTP(t *testing.T) *core.TTP {
	dbType := golangutils.GetEnv("DB_TYPE", "sqlite3")
	dbConfig := golangutils.GetEnv("DB_CONFIG", "file:ent?mode=memory&cache=shared&_fk=1")

	dbc := core.DBConfig{
		Type:   dbType,
		Config: dbConfig,
	}

	db, err := core.NewDB(&dbc)
	assert.NoError(t, err)

	ttp, err := core.NewTTP(db, nil, "")
	assert.NoError(t, err)

	return ttp
}

func TestAll(t *testing.T) {
	rawebscore.EnableDebug()

	t.Run("Pass", testPass)
	t.Run("FailTANoServer", testFailTANoServer)
	t.Run("FailByMissDomains", testFailByMissDomains)
}

func testPass(t *testing.T) {
	ttp := exampleTTP(t)
	defer ttp.DB.Close()

	ttp.DB.Client.TAServer.Create().SetDomain("example.com").SaveX(*ttp.DB.Ctx)
	ttp.DB.Client.TACode.Create().SetUniqueID([]byte{1, 2, 3}).SetRepository("").SetCommitID("").SaveX(*ttp.DB.Ctx)

	err := AuditOne(ttp, &x509.Certificate{
		DNSNames:  []string{"example.com"},
		PublicKey: []byte{7, 8, 9},
		Extensions: []pkix.Extension{
			{
				Id:       rawebscore.X509_EXTENSION_LABEL,
				Critical: false,
				Value:    []byte{},
			},
		},
	},
	)

	assert.NoError(t, err)
}

func testFailTANoServer(t *testing.T) {
	ttp := exampleTTP(t)
	defer ttp.DB.Close()

	ttp.DB.Client.TAServer.Create().SetDomain("example.com").SaveX(*ttp.DB.Ctx)
	ttp.DB.Client.TACode.Create().SetUniqueID([]byte{1, 2, 3}).SetRepository("").SetCommitID("").SaveX(*ttp.DB.Ctx)

	err := AuditOne(ttp, &x509.Certificate{
		DNSNames:  []string{"hoge.example.com"},
		PublicKey: []byte{7, 8, 9},
		Extensions: []pkix.Extension{
			{
				Id:       rawebscore.X509_EXTENSION_LABEL,
				Critical: false,
				Value:    []byte{},
			},
		},
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), ERROR_SELECT_SERVER)
}

func testFailByMissDomains(t *testing.T) {
	// TestMultipleDomain tests the revokeTAByDomains function
	// by revoking multiple TAs by their domain names.

	ttp := exampleTTP(t)
	defer ttp.DB.Close()

	cert := x509.Certificate{
		DNSNames:  []string{"example.com", "example.org"},
		PublicKey: []byte{7, 8, 9},
		Extensions: []pkix.Extension{
			{
				Id:       rawebscore.X509_EXTENSION_LABEL,
				Critical: false,
				Value:    []byte{},
			},
		},
	}

	err := AuditOne(ttp, &cert)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), ERROR_DOMAIN_INVALID)

	cert = x509.Certificate{
		DNSNames:  []string{"*.com"},
		PublicKey: []byte{7, 8, 9},
		Extensions: []pkix.Extension{
			{
				Id:       rawebscore.X509_EXTENSION_LABEL,
				Critical: false,
				Value:    []byte{},
			},
		},
	}

	err = AuditOne(ttp, &cert)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), ERROR_DOMAIN_INVALID)

}
