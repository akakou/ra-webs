package ct

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"testing"
	"time"

	"github.com/akakou/ra_webs/ttp/core"
)

func TestSSLMate(t *testing.T) {
	ct := DefaultSSLMateCT("")
	ct.MaxSleep = time.Second * 20

	db, _ := core.NewDB(&core.DBConfig{
		Type:   "sqlite3",
		Config: "file:ent?mode=memory&cache=shared&_fk=1",
	})
	defer db.Close()

	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	keyBuf := x509.MarshalPKCS1PublicKey(&priv.PublicKey)

	ttp, _ := core.NewTTP(db, ct, "admin_token")

	service := ttp.DB.Client.Service.Create().SetName("sslmate").SetToken("aaaa").SetIsActive(true).SaveX(*ttp.DB.Ctx)
	ttp.DB.Client.TACode.Create().SetUniqueID([]byte{1, 2, 3}).SetRepository("").SetCommitID("").SaveX(*ttp.DB.Ctx)
	ttp.DB.Client.TAServer.Create().SetDomain("example.com").SetPublicKey(keyBuf).SetQuote("1").SetHasActivated(false).SetService(service).SaveX(*ttp.DB.Ctx)

	err := ct.Setup(nil, ttp)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	for {
	}
}
