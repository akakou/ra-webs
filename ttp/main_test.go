package ttp

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
	dbConfig := DBConfig{
		Type:   "sqlite3",
		Config: "file:ent?mode=memory&cache=shared&_fk=1",
	}

	router := NewTTPServer(&dbConfig, "views/*.html")

	postBody := core.TAInfo{
		Attestation: "attestation",
		PublicKey:   []byte("public_key"),
		Domain:      "domain",
	}

	body, _ := json.Marshal(postBody)

	req := httptest.NewRequest("POST", "/provision", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}
