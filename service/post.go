package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"core"

	goutils "github.com/akakou/go-utils"
	"github.com/labstack/echo/v4"
)

func (service *Service) post(path string, reqBody any) (string, error) {
	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, service.TTP_BASE+path, strings.NewReader(string(body)))
	if err != nil {
		return "", err
	}

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", service.Token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to post code: %v", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(respBody), nil
}

func (service *Service) PostCode(repository string) (string, error) {
	return service.post("/code", map[string]string{"repository": repository})
}

func (service *Service) PostServer(e *echo.Echo, domain string) (string, error) {
	service.DomainAuthServer(e, domain)
	go e.Start(PORT)

	return service.post("/server", map[string]string{"domain": domain})
}

func (service *Service) DomainAuthServer(e *echo.Echo, domain string) *echo.Echo {
	nonce, _ := goutils.RandomHex(64)
	token := core.DomainToken(service.Token, nonce)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, token)
	})

	go e.Start(PORT)
	return e
}
