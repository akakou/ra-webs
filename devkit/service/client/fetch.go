package serviceclient

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	sc "github.com/akakou/ra-webs/monitor/serviceclient"
)

func (serviceClient *ServiceClient) Fetch(publicKey []byte) (*sc.EvidenceEntry, error) {
	u, err := url.Parse(serviceClient.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)

	}

	u.Path = "/api/ta"

	q := u.Query()
	encodedPublicKey := base64.URLEncoding.EncodeToString(publicKey)
	q.Set("public_key", encodedPublicKey)
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result *sc.EvidenceEntry
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result, nil
}
