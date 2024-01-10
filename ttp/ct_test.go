package ttp

import (
	"fmt"
	"os"
	"testing"
)

func TestFetchCTLogs(t *testing.T) {
	const TEST_SITE = "hoge.ochano.co"
	const TEST_PUBLIC_KEY_HASH = "e214ae98f05708cbf4487415fb4a13cd0865d78c512ca80aad06a4ee0e39a2f4"

	var testCT = os.Getenv("RAWEBS_TEST_CT")
	if testCT != "1" {
		t.Skip("RAWEBS_TEST_CT is not set")
	}

	result, err := fetchCTLogs(TEST_SITE, "")
	if err != nil {
		t.Fatalf("fetchCTLogs() got an unexpected error: %s", err)
	}

	if len(result) == 0 {
		t.Fatal("fetchCTLogs() got: empty slice, want: non-empty slice")
	}

	for _, ct := range result {
		fmt.Printf("fetchCTLogs() \ngot: %s\nwant: %s\n\n", ct.PubKeySha256, TEST_PUBLIC_KEY_HASH)
	}
}
