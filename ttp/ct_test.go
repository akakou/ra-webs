package ttp

var test = `
{
  "data": [
    {
      "cert_hash_sha256": "152151b387412a95ab90fd4664f2d98880408e9e43913124e4c91226a483386b",
      "domains": [
        "*.facebook.com",
        "*.facebook.net",
        "*.fb.com",
        "*.fbcdn.net",
        "*.fbsbx.com",
        "*.m.facebook.com",
        "*.messenger.com",
        "*.xx.fbcdn.net",
        "*.xy.fbcdn.net",
        "*.xz.fbcdn.net",
        "facebook.com",
        "fb.com",
        "messenger.com"
      ],
      "issuer_name": "/C=US/O=DigiCert Inc/OU=www.digicert.com/CN=DigiCert SHA2 High Assurance Server CA",
      "certificate_pem": "-----BEGIN CERTIFICATE-----
MIIH5DCCBsygAwIBAgIQDACZt9eJyfZmJjF+vOp8HDANBgkqhkiG9w0BAQsFADBw
MQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3
d3cuZGlnaWNlcnQuY29tMS8wLQYDVQQDEyZEaWdpQ2VydCBTSEEyIEhpZ2ggQXNz
dXJhbmNlIFNlcnZlciBDQTAeFw0xNjEyMDkwMDAwMDBaFw0xODAxMjUxMjAwMDBa
MGkhMRMwEQYDVQQHEwpNZW5sbyBQYXJrMRcwFQYDVQQKEw5GYWNlYm9vaywgSW5j
LjEXMBUGA1UEAwwOKi5mHwYDVR0jBBgwFoAUUWj/kK8CB3U8zNllZGKiErhZcjsw
YWNlYm9vay5jb20wWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAASg8YyvpzmIaFsT
Vg4VFbSnRe8bx+WFPCsE1GWKMTEi6qOS7WSdumWB47YSdtizC0Xx/wooFJxP3HOp
s0ktoHbTo4IFSjCCBUYwHwYDVR0jBBgwFoAUUWj/kK8CB3U8zNllZGKiErhZcjsw
HQYDVR0OBBYEFMuYKIyhcufiMqmaPfINoYFWoRqLMIHHBgNVHREEgb8wgbyCDiou
ZmFjZWJvb2suY29tgg4qLmZhY2Vib29rLm5ldIIIKi5mYi5jb22CCyouZmJjZG4u
bmV0ggsqLmZic2J4LmNvbYIQKi5tLmZ
-----END CERTIFICATE-----",
      "id": "1662768163744657"
    }
  ]
}
`

// func TestFetchCTLogs(t *testing.T) {
// 	const TEST_SITE = "hoge.ochano.co"
// 	const TEST_PUBLIC_KEY_HASH = "e214ae98f05708cbf4487415fb4a13cd0865d78c512ca80aad06a4ee0e39a2f4"

// 	var testCT = os.Getenv("RAWEBS_TEST_CT")
// 	if testCT != "1" {
// 		t.Skip("RAWEBS_TEST_CT is not set")
// 	}

// 	result, err := fetchCTLogs(TEST_SITE, "")
// 	assert.NoError(t, err)
// 	assert.NotEmpty(t, result)

// 	for _, ct := range result {
// 		fmt.Printf("fetchCTLogs() \ngot: %s\nwant: %s\n\n", TEST_PUBLIC_KEY_HASH)
// 	}
// }
