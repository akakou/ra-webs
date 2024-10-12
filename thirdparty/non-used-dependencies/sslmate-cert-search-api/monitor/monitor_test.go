package monitor

import "testing"

func TestNext(t *testing.T) {
	// TestNext tests the Next method of the Monitor struct.
	certs, _, err := Default("example.com").Next()
	if err != nil {
		t.Fatalf("failed to get certs: %v", err)
	}

	t.Logf("%v", certs)
}
