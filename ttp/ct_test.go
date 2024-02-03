package ttp

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var test = `
{
  "entry": [
    {
      "changes": [
        {
          "field": "phishing",
          "value": {
            "ct_cert": {
              "id": "123",
              "certificate_pem": "-----BEGIN CERTIFICATE----- ... -----END CERTIFICATE-----\n",
              "cert_hash_sha256": "f2297..."
            },
            "phishing_domains": [
              "facebook.com.evil.com",
              "xn—facbook-9gg.ml"
            ],
            "phished_domain": "facebook.com"
          }
        }
      ],
      "id": "123",
      "time": 1524762838
    }
  ],
  "object": "certificate_transparency"
}
`

func TestFetchCTLogs(t *testing.T) {
	const TEST_SITE = "hoge.ochano.co"
	const TEST_PUBLIC_KEY_HASH = "e214ae98f05708cbf4487415fb4a13cd0865d78c512ca80aad06a4ee0e39a2f4"

	var testCT = os.Getenv("RAWEBS_TEST_CT")
	if testCT != "1" {
		t.Skip("RAWEBS_TEST_CT is not set")
	}

	result, err := fetchCTLogs(TEST_SITE, "")
	assert.NoError(t, err)
	assert.NotEmpty(t, result)

	for _, ct := range result {
		fmt.Printf("fetchCTLogs() \ngot: %s\nwant: %s\n\n", ct.PubKeySha256, TEST_PUBLIC_KEY_HASH)
	}
}
