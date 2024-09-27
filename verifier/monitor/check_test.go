package monitor

import (
	"crypto/rand"
	"crypto/rsa"

	"github.com/google/certificate-transparency-go/x509"

	"testing"

	golangutils "github.com/akakou/golang-utils"
	rawebscore "github.com/akakou/ra_webs/core"
	"github.com/akakou/ra_webs/verifier/core"
	"github.com/stretchr/testify/assert"
)

func exampleVerifier(t *testing.T) *core.Verifier {
	dbType := golangutils.GetEnv("DB_TYPE", "sqlite3")
	dbConfig := golangutils.GetEnv("DB_CONFIG", "file:ent?mode=memory&cache=shared&_fk=1")

	dbc := core.DBConfig{
		Type:   dbType,
		Config: dbConfig,
	}

	db, err := core.NewDB(&dbc)
	assert.NoError(t, err)

	verifier, err := core.NewVerifier(db, nil, nil, "")
	assert.NoError(t, err)

	return verifier
}

func TestAll(t *testing.T) {
	rawebscore.EnableDebug()

	t.Run("Pass", testPass)
}

func testPass(t *testing.T) {
	verifier := exampleVerifier(t)
	defer verifier.DB.Close()

	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	keyBuf := x509.MarshalPKCS1PublicKey(&priv.PublicKey)

	server := verifier.DB.Client.TAServer.Create().SetDomain("example.com").SetPublicKey(keyBuf).SetQuote("1").SetIsActive(false).SaveX(*verifier.DB.Ctx)
	verifier.DB.Client.TACode.Create().SetUniqueID([]byte{1, 2, 3}).SetRepository("").SetCommitID("").SaveX(*verifier.DB.Ctx)

	err := Check(&priv.PublicKey, server)

	assert.NoError(t, err)
}
