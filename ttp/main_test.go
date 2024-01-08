package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/akakou/ra_webs/core"
	"github.com/go-playground/assert/v2"
)

func TestRegister(t *testing.T) {
	router := NewRouter()

	postBody := core.ProvisioningRequest{
		Attestation: "attestation",
		PublicKey:   "public_key",
		Domain:      "domain",
	}

	body, _ := json.Marshal(postBody)

	req := httptest.NewRequest("POST", "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}
