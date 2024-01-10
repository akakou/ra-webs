package ttp

// func TestProvision(t *testing.T) {
// 	dbConfig := DBConfig{
// 		Type:   "sqlite3",
// 		Config: "file:ent?mode=memory&cache=shared&_fk=1",
// 	}

// 	router := NewTTPServer(&dbConfig, "views/*.html")

// 	postBody := core.ProvisionRequest{
// 		Attestation: "attestation",
// 		Domain:      "domain",
// 	}

// 	body, _ := json.Marshal(postBody)

// 	req := httptest.NewRequest("POST", "/provision", bytes.NewReader(body))
// 	req.Header.Set("Content-Type", "application/json")

// 	rec := httptest.NewRecorder()

// 	router.ServeHTTP(rec, req)

// 	assert.Equal(t, http.StatusOK, rec.Code)
// }
