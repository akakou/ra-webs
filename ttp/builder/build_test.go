package builder

import "testing"

const EXPECTED_COMMIT_ID = "b9d0bc9519799ae2b07ce033d5dd942d7c275beb"
const EXPECTED_UNIQUE_ID = "d239da309e3568497cded1caa9583ef67e9acf76587218315a913709c2db05c0"

func TestMain(t *testing.T) {
	commitId, uniqueId, err := build("1", "https://github.com/akakou-docs/ego-statistical-analysis")
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
