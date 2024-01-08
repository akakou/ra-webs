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
	db, err := newtTAInfoDB("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	assert.Equal(t, nil, err)

	router := NewRouter(db)

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
