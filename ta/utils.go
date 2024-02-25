package ta

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/labstack/echo"
)

func (ap *AttestProxy) requestToTTP(url, reqBody string) ([]byte, error) {
	req := httptest.NewRequest(http.MethodPost, url, strings.NewReader(string(reqBody)))
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", ap.Config.Token))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to register: %s", resp.Status)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}
