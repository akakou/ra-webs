package ta

import (
	"testing"
)

func TestRAConfiggenerateKeyPair(t *testing.T) {
	raConfig := RAConfig{
		TTPDomain: "ttp.example.com",
		Domain:    "ta.example.com",
	}

	got, _, err := raConfig.generateKeyPair()
	if err != nil {
		t.Errorf("RAConfig.generateKeyPair() got an unexpected error: %s", err)
	}
	if got == nil {
		t.Errorf("RAConfig.generateKeyPair() got: nil, want: non-nil")
	}
}
