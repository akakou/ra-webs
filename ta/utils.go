package ta

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/akakou/ra_webs/core"
	"github.com/labstack/echo"
)

func (ap *TA) requestToTTP(url string, reqBodyJson map[string]any) ([]byte, error) {
	reqBody, err := json.Marshal(reqBodyJson)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

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

func attestPublicKey(ap *TA) (string, error) {
	publicKey := ap.PrivateKey.Public()
	publicKeyBuf := x509.MarshalPKCS1PublicKey(publicKey.(*rsa.PublicKey))

	quote, err := core.AttestByAzure(publicKeyBuf)
	return quote, err
}
