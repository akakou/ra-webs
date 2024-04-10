package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/akakou/ra_webs/core"
	"github.com/akakou/ra_webs/ttp/api"

	goutils "github.com/akakou/go-utils"
	"github.com/labstack/echo/v4"
)

func (service *Service) post(path string, reqBody any) (string, error) {
	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	u, err := url.Parse(service.TTPBase)
	if err != nil {
		return "", fmt.Errorf("%v: %v", ERROR_TTP_BASE_PARSE, err)
	}

	u.Path = path
	fmt.Printf("url: %v", u.String())

	req, err := http.NewRequest(http.MethodPost, u.String(), strings.NewReader(string(body)))
	if err != nil {
		return "", fmt.Errorf("%v: %v", ERROR_REQUEST_FAILED, err)
	}

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", service.Token))

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

func (service *Service) PostCode(repository string) (string, error) {
	return service.post(api.PostCodeApi.Path, map[string]string{"repository": repository})
}

func (service *Service) PostServer(e *echo.Echo, domain, port string) (string, error) {
	nonce, _ := goutils.RandomHex(64)
	token := core.DomainToken(service.Token, nonce)

	service.DomainAuthServer(token, domain, port, e)

	return service.post(api.PostServerApi.Path, map[string]string{"domain": domain + port, "nonce": nonce})
}

func (service *Service) DomainAuthServer(token, domain, port string, e *echo.Echo) *echo.Echo {
	e.GET(api.DOMAIN_AUTH_PATH, func(c echo.Context) error {
		return c.String(http.StatusOK, token)
	})

	go e.Start(port)
	return e
}
