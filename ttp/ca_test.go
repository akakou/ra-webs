package ttp

import (
	"crypto/x509/pkix"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadAndSaveCA(t *testing.T) {
	config_path = "ca_test.json"

	ca, err := initCA(
		pkix.Name{
			Country:            []string{"US"},
			Organization:       []string{"Google"},
			OrganizationalUnit: []string{"Google"},
			Locality:           []string{"Mountain View"},
			Province:           []string{"California"},
			StreetAddress:      []string{"1600 Amphitheatre Parkway"},
			PostalCode:         []string{"94043"},
		},
	)

	assert.NoError(t, err)

	err = saveCA(ca)
	assert.NoError(t, err)

	ca2, err := loadCA()
	assert.NoError(t, err)

	assert.Equal(t, ca.PrivateKey, ca2.PrivateKey)
	assert.Equal(t, ca.Certificate, ca2.Certificate)
}
