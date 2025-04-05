package builder

import (
	"fmt"
	"testing"
)

const REPOSITORY = "https://github.com/akakou-docs/ego-statistical-analysis"
const COMMIT_ID = "b9d0bc9519799ae2b07ce033d5dd942d7c275beb"
const EXPECTED_UNIQUE_ID = "ea0e1b510d48a85981d5a0b3794bb08b75962a67bd44798cdeeaff58745d701a"

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
