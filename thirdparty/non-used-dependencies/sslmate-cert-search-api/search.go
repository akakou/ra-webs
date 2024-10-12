package sslmate_cert_search_api

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func search[T any](any T, query *Query, api *SSLMateSearchAPI) error {
	url, _ := url.Parse(SSLMATE_API_URL)
	urlQuery := url.Query()

	urlQuery.Add("domain", query.Domain)
	urlQuery.Add("include_subdomains", fmt.Sprintf("%t", query.IncludeSubdomains))
	urlQuery.Add("match_wildcards", fmt.Sprintf("%t", query.MatchWildcards))
	urlQuery.Add("after", query.After)
	urlQuery.Add("expand", query.Expand)
	url.RawQuery = urlQuery.Encode()

	fmt.Printf("url: %v\n", url.String())

	req := http.Request{
		Method: http.MethodGet,
		URL:    url,
		Header: http.Header{
			// "Authorization": []string{fmt.Sprintf("Bearer %s", api.Token)},
			"Content-Type": []string{"application/json"},
		},
	}

	if api.Token != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", api.Token))
	}

	client := new(http.Client)
	resp, err := client.Do(&req)
	if err != nil {
		return fmt.Errorf("%v: %v", ERROR_FAILED_TO_FETCH, err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%v: %v", ERROR_SERVER_RETURN_NOT_200, resp.StatusCode)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%v: %v", ERROR_FAILED_TO_READ_BODY, err)
	}

	err = json.Unmarshal(body, any)
	if err != nil {
		return fmt.Errorf("%v: %v", ERROR_FAILED_TO_MARSHAL_JSON, err)
	}
	return err
}

func (api *SSLMateSearchAPI) SearchJsonCerts(query *Query) (Json, error) {
	var result Json
	err := search(&result, query, api)
	return result, err
}

func (api *SSLMateSearchAPI) SearchSSLMateCerts(query *Query) ([]SSLMateCert, error) {
	query.Expand = "cert_der"

	var result []SSLMateCert
	err := search(&result, query, api)
	return result, err
}

func (api *SSLMateSearchAPI) Search(query *Query) ([]x509.Certificate, *Index, error) {
	var result []x509.Certificate
	var sslmateCerts []SSLMateCert

	sslmateCerts, err := api.SearchSSLMateCerts(query)
	if err != nil {
		return nil, nil, err
	}

	for _, c := range sslmateCerts {
		raw, _ := base64.StdEncoding.DecodeString(c.RawCert)
		x509Cert, err := x509.ParseCertificate(raw)
		if err != nil {
			return nil, nil, fmt.Errorf("%v: %v", ERROR_FAILED_TO_PARSE_CERT, err)
		}
		result = append(result, *x509Cert)
	}

	index := Index{
		First: sslmateCerts[0].Id,
		Last:  sslmateCerts[len(sslmateCerts)-1].Id,
	}

	return result, &index, nil
}
