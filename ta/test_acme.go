package ta

import (
	"crypto/rand"
	"crypto/rsa"
	"log"
	"os"
	"testing"
)

func TestAcmeAll(t *testing.T) {
	email := os.Getenv("RAWEBS_TEST_EMAIL")
	domain := os.Getenv("RAWEBS_TEST_DOMAIN")

	user, err := CreateUser(email)
	if err != nil {
		t.Fatalf("failed to create acme user: %v\n", err)
	}

	acmeClient, err := NewStagingAcmeClient(user)
	if err != nil {
		t.Fatalf("failed to create acme client: %v\n", err)
	}

	err = SetupProvider(acmeClient)
	if err != nil {
		t.Fatalf("failed to setup provider: %v\n", err)
	}

	err = user.Register(acmeClient)
	if err != nil {
		t.Fatalf("failed to register acme user: %v\n", err)
	}

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate private key: %v\n", err)
	}

	cert, err := user.ObtainCertificate(acmeClient, priv, []string{domain})
	if err != nil {
		t.Fatalf("failed to obtain certificates: %v\n", err)
	}

	log.Printf("Test_AcmeAll: Obtained certificate for %v\n", cert.Domain)
}

// func Test_AcmeAllWithPrivateStorage(t *testing.T) {
// 	storage, err := CreatePrivateStorage(testDomain)
// 	if err != nil {
// 		t.Fatalf("creating private storage: %v\n", err)
// 	}

// 	user := storage.AcmeUser(EMAIL)

// 	acmeClient, err := NewStagingAcmeClient(user)
// 	if err != nil {
// 		t.Fatalf("failed to create acme client: %v\n", err)
// 	}

// 	err = SetupProvider(acmeClient)
// 	if err != nil {
// 		t.Fatalf("failed to setup provider: %v\n", err)
// 	}

// 	err = user.Register(acmeClient)
// 	if err != nil {
// 		t.Fatalf("failed to register acme user: %v\n", err)
// 	}

// 	cert, err := user.ObtainCertificate(acmeClient, storage.PrivateKey, []string{DOMAIN2})
// 	if err != nil {
// 		t.Fatalf("failed to obtain certificates: %v\n", err)
// 	}

// 	log.Printf("Test_AcmeAll: Obtained certificate for %v\n", cert.Domain)
// }

// func Test_IssueCertificate(t *testing.T) {
// 	storage, err := CreatePrivateStorage(testDomain)
// 	if err != nil {
// 		t.Fatalf("creating private storage: %v\n", err)
// 	}

// 	user := storage.AcmeUser(EMAIL)

// 	cert, err := IssueCertiifcate(user, storage.PrivateKey, []string{DOMAIN2}, false)
// 	if err != nil {
// 		t.Errorf("failed to issue certificate: %v", err)
// 	}

// 	t.Logf("Test_AcmeAll: Obtained certificate for %v\n", cert.Domain)
// }
