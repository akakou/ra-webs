package builder

import (
	"fmt"
	"testing"
)

const REPOSITORY = "https://github.com/akakou/ra-webs"
const COMMIT_ID = "8bc46f9bf7569a0d3c21f37bdeca94c54f504806"
const EXPECTED_UNIQUE_ID = "4759f05537868a6dbfbd2bf1109b8805c49a40ef4cd37ed4e5c446743523d4e5"

func TestMain(t *testing.T) {
	uniqueId, err := buildCode("1", REPOSITORY, COMMIT_ID)
	fmt.Printf("commit id: %v\nunique id: %v", COMMIT_ID, uniqueId)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if uniqueId != EXPECTED_UNIQUE_ID {
		t.Errorf("Expected: %v, got: %v", EXPECTED_UNIQUE_ID, uniqueId)
	}
}
