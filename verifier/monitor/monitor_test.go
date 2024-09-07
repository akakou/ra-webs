package monitor

import (
	"crypto/rand"
	"crypto/rsa"

	"github.com/google/certificate-transparency-go/x509"
	"github.com/google/certificate-transparency-go/x509/pkix"

	"encoding/pem"
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
	t.Run("FailTANoServer", testFailTANoServer)
	t.Run("FailByMissDomains", testFailByMissDomains)
}

func testPass(t *testing.T) {
	verifier := exampleVerifier(t)
	defer verifier.DB.Close()

	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	keyBuf := x509.MarshalPKCS1PublicKey(&priv.PublicKey)

	verifier.DB.Client.TAServer.Create().SetDomain("example.com").SetPublicKey(keyBuf).SetQuote("1").SetHasActivated(false).SaveX(*verifier.DB.Ctx)
	verifier.DB.Client.TACode.Create().SetUniqueID([]byte{1, 2, 3}).SetRepository("").SetCommitID("").SaveX(*verifier.DB.Ctx)

	err := Monitor(verifier, &x509.Certificate{
		DNSNames:  []string{"example.com"},
		Subject:   pkix.Name{CommonName: "example.com"},
		PublicKey: &priv.PublicKey,
	})

	assert.NoError(t, err)
}

func testFailTANoServer(t *testing.T) {
	verifier := exampleVerifier(t)
	defer verifier.DB.Close()

	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	keyBuf := x509.MarshalPKCS1PublicKey(&priv.PublicKey)

	verifier.DB.Client.TAServer.Create().SetDomain("example.com").SetPublicKey(keyBuf).SetQuote("1").SetHasActivated(true).SaveX(*verifier.DB.Ctx)
	verifier.DB.Client.TACode.Create().SetUniqueID([]byte{1, 2, 3}).SetRepository("").SetCommitID("").SaveX(*verifier.DB.Ctx)

	err := Monitor(verifier, &x509.Certificate{
		DNSNames:  []string{"hoge.example.com"},
		Subject:   pkix.Name{CommonName: "hoge.example.com"},
		PublicKey: &priv.PublicKey,
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), ERROR_CERTIFICATE_NOT_FOUND)
}

func testFailByMissDomains(t *testing.T) {
	// TestMultipleDomain tests the revokeTAByDomains function
	// by revoking multiple TAs by their domain names.

	verifier := exampleVerifier(t)
	defer verifier.DB.Close()

	cert := x509.Certificate{
		DNSNames:  []string{"example.com", "example.org"},
		PublicKey: []byte{7, 8, 9},
	}

	err := Monitor(verifier, &cert)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), ERROR_DOMAIN_INVALID)

	cert = x509.Certificate{
		DNSNames:  []string{"*.com"},
		PublicKey: []byte{7, 8, 9},
	}

	err = Monitor(verifier, &cert)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), ERROR_DOMAIN_INVALID)

}

func TestRealCert(t *testing.T) {
	str := `-----BEGIN CERTIFICATE-----
MIIDLTCCAhWgAwIBAgISA0j8iAuL3Swtl7wnUpn2glJuMA0GCSqGSIb3DQEBCwUA
MDIxCzAJBgNVBAYTAlVTMRYwFAYDVQQKEw1MZXQncyBFbmNyeXB0MQswCQYDVQQD
EwJSMzAeFw0yNDA1MDUxOTEzMDdaFw0yNDA4MDMxOTEzMDZaMBoxGDAWBgNVBAMT
D3Rlc3QyLm9jaGFuby5jbzBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABD5o+Zg0
4dgfFn4WCdGiIvMKGv6mezTbrbNbe5D/IM/3jmtAf01ynrIpGadUx+jrvcRudqyO
0Yk+7xBFWrov4X+jggEeMIIBGjAOBgNVHQ8BAf8EBAMCB4AwHQYDVR0lBBYwFAYI
KwYBBQUHAwEGCCsGAQUFBwMCMAwGA1UdEwEB/wQCMAAwHQYDVR0OBBYEFIw0jiKg
ImmIKTV3MA3G2ERyG8WzMB8GA1UdIwQYMBaAFBQusxe3WFbLrlAJQOYfr52LFMLG
MFUGCCsGAQUFBwEBBEkwRzAhBggrBgEFBQcwAYYVaHR0cDovL3IzLm8ubGVuY3Iu
b3JnMCIGCCsGAQUFBzAChhZodHRwOi8vcjMuaS5sZW5jci5vcmcvMBoGA1UdEQQT
MBGCD3Rlc3QyLm9jaGFuby5jbzATBgNVHSAEDDAKMAgGBmeBDAECATATBgorBgEE
AdZ5AgQDAQH/BAIFADANBgkqhkiG9w0BAQsFAAOCAQEASilPljBjCnU6I6akocU8
2iBlk5N7sccyE6C0iE5WBiNp49EfEsxvh5EFqqtYWGZ0I82gLy1V38TNp2Y+xCHw
5gVNQ4IhpGC58yRkKgCREtx6vHbdb6OfZfWEqdiqFrW+xrbJl9qNUKYwO7WFcDa+
q8VR55lAxm81XMIilEKdxigsm+ZTmtxgqTC/BIicET9G6Eaqj5os1+UvRMxLErED
oxPKuqFmwydi8Av5EI6bk33GOTlOxtJaOFpkmJ7X6sUiakSLQSpffHzsvsjFYFaW
6DAIckSIckj41e6jGG4XtiP2j7Rb+4YXtciEJrT8owgV2xgiyw4yV799/8t63d1V
ww==
-----END CERTIFICATE-----`

	c, _ := pem.Decode([]byte(str))
	cert, err := x509.ParseCertificate(c.Bytes)
	assert.NoError(t, err)

	// verifier := exampleVerifier(t)
	_, err = validateDomains(cert)
	assert.NoError(t, err)
}
