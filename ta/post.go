package ta

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/akakou/ra_webs/core"
)

const WAIT = 3
const REGISTER_PATH = "/api/ta"

func (ta *TA) Register() (string, error) {
	publicKey := ta.privateKey.Public()
	keyBin := x509.MarshalPKCS1PublicKey(publicKey.(*rsa.PublicKey))

	quote, err := core.AttestServer(keyBin, ta.config.Token)
	if err != nil {
		return "", fmt.Errorf("%s: %w", ERROR_ATTEST_PUBLIC_KEY, err)
	}

	reqBody := core.RegisterRequest{
		CodeRequest: core.CodeRequest{
			Repository: ta.config.Repository,
		},
		ServerRequest: core.ServerRequest{
			PublicKey: keyBin,
			Domain:    ta.config.Domain,
			Quote:     quote,
		},
	}

	return ta.post(REGISTER_PATH, reqBody)
}

func (ta *TA) post(path string, reqBody any) (string, error) {
	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	u, err := url.Parse(ta.config.Verifier)
	if err != nil {
		return "", fmt.Errorf("%v: %v", ERROR_VERIFIER_BASE_PARSE, err)
	}

	u.Path = path

	req, err := http.NewRequest(http.MethodPost, u.String(), strings.NewReader(string(body)))
	if err != nil {
		return "", fmt.Errorf("%v: %v", ERROR_REQUEST_FAILED, err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ta.config.Token))
	req.Close = true

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("%v: %v", ERROR_REQUEST_FAILED, err)
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("%v: %v", ERROR_READ_BODY, err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%v, %v: %v(%v)", ERROR_STATUS_NOT_OK, err, string(respBody), resp.StatusCode)
	}

	return string(respBody), nil
}
