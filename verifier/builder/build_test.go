package builder

import (
	"fmt"
	"testing"
)

const EXPECTED_COMMIT_ID = "b9d0bc9519799ae2b07ce033d5dd942d7c275beb"
const EXPECTED_UNIQUE_ID = "ea0e1b510d48a85981d5a0b3794bb08b75962a67bd44798cdeeaff58745d701a"

func TestMain(t *testing.T) {
	commitId, uniqueId, err := buildCode("1", "https://github.com/akakou-docs/ego-statistical-analysis")
	fmt.Printf("commit id: %v\nunique id: %v", commitId, uniqueId)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	if commitId != EXPECTED_COMMIT_ID {
		t.Errorf("Expected: %v, got: %v", EXPECTED_COMMIT_ID, commitId)
	}

	if uniqueId != EXPECTED_UNIQUE_ID {
		t.Errorf("Expected: %v, got: %v", EXPECTED_UNIQUE_ID, uniqueId)
	}
}
