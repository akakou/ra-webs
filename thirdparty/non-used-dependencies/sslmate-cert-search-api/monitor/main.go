package monitor

import api "github.com/akakou/sslmate_cert_search_api"

var DefaultQuery = api.Query{
	IncludeSubdomains: true,
	MatchWildcards:    true,
	After:             "",
}

type Monitor struct {
	Domains string
	Api     *api.SSLMateSearchAPI
	Query   *api.Query
}

type Monitors []Monitor

func New(base *api.Query, api *api.SSLMateSearchAPI) *Monitor {
	return &Monitor{
		Query: base,
		Api:   api,
	}
}

func Default(domain string) *Monitor {
	query := DefaultQuery
	query.Domain = domain

	return New(
		&query,
		api.Default(),
	)
}
