package monitor

import (
	"crypto/rand"
	"crypto/rsa"

	"github.com/google/certificate-transparency-go/x509"

	"testing"

	golangutils "github.com/akakou/golang-utils"
	"github.com/akakou/ra_webs/monitor"
	"github.com/stretchr/testify/assert"

	core "github.com/akakou/ra_webs/core"
)

func exampleVerifier(t *testing.T) *monitor.Monitor {
	dbType := golangutils.GetEnv("DB_TYPE", "sqlite3")
	dbConfig := golangutils.GetEnv("DB_CONFIG", "file:ent?mode=memory&cache=shared&_fk=1")

	dbc := monitor.DBConfig{
		Type:   dbType,
		Config: dbConfig,
	}

	db, err := monitor.NewDB(&dbc)
	assert.NoError(t, err)

	monitor, err := monitor.New(db)
	assert.NoError(t, err)

	return monitor
}

func TestAll(t *testing.T) {
	core.EnableDebug()

	t.Run("Pass", testPass)
}

func testPass(t *testing.T) {
	monitor := exampleVerifier(t)
	defer monitor.DB.Close()

	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	keyBuf := x509.MarshalPKCS1PublicKey(&priv.PublicKey)

	server := monitor.DB.Client.TAServer.Create().SetPublicKey(keyBuf).SetQuote("1").SetMonitorLogID(0).SetIsActive(false).SaveX(*monitor.DB.Ctx)
	monitor.DB.Client.TACode.Create().SetUniqueID([]byte{1, 2, 3}).SetRepository("").SetCommitID("").SaveX(*monitor.DB.Ctx)

	err := Check(&priv.PublicKey, server)

	assert.NoError(t, err)
}
