package logclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	logio "github.com/akakou/ra-webs/log/api/io"
)

var schema = "http://"

func (logclient *LogClient) Fetch() ([]*logio.TA, error) {
	resp, err := http.Get(schema + logclient.Domain + "/api/ta")
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

	var result []*logio.TA
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result, nil
}
