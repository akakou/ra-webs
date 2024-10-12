package sslmate_cert_search_api

import (
	"testing"
)

func TestSearchJson(t *testing.T) {
	api := Default()

	query := Query{
		Domain:            "example.com",
		IncludeSubdomains: true,
		MatchWildcards:    true,
		After:             "",
		Expand:            "",
	}

	certs, err := api.SearchJsonCerts(&query)
	if err != nil {
		t.Fatalf("failed to search json certs: %v", err)
	}

	t.Logf("%v", certs)
}

func TestSearch(t *testing.T) {
	api := Default()

	query := Query{
		Domain:            "example.com",
		IncludeSubdomains: true,
		MatchWildcards:    true,
		After:             "",
		Expand:            "",
	}

	certs, _, err := api.Search(&query)
	if err != nil {
		t.Fatalf("failed to search certs: %v", err)
	}

	t.Logf("%v", certs)
}
