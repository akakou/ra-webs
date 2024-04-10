package api

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	ttpcore "github.com/akakou/ra_webs/ttp/core"
	"github.com/akakou/ra_webs/ttp/ent"
	"github.com/akakou/ra_webs/ttp/ent/service"
	"github.com/labstack/echo/v4"
)

const (
	ERROR_AUTHENTICATE_SERVICE      = "failed to authenticate service"
	ERROR_AUTHENTICATE_ADMIN        = "failed to authenticate admin"
	ERROR_ACCESS_DOMAIN_AUTH_TARGET = "failed to access domain auth target"
	ERROR_DOMAIN_AUTH_TOKEN_INVALID = "domain auth token is invalid"
)

var SCHEME = "https"
var DOMAIN_AUTH_PATH = "/ra-webs"

func authenticateService(ttp *ttpcore.TTP, c echo.Context) (*ent.Service, error) {
	authorization := c.Request().Header["Authorization"][0]
	token := authorization[len("Bearer "):]

	service, err := ttp.DB.Client.Service.Query().Where(service.TokenEQ(token)).Only(*ttp.DB.Ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ERROR_AUTHENTICATE_SERVICE, err)
	}

	return service, nil
}

func authenticateAdmin(ttp *ttpcore.TTP, c echo.Context) error {
	authorization := c.Request().Header["Authorization"][0]
	token := authorization[len("Bearer "):]

	if token != ttp.AdminToken {
		return fmt.Errorf(ERROR_AUTHENTICATE_ADMIN)
	}

	return nil
}

func authenticateDomain(domain, serviceToken, nonce string) error {
	hashSource := []byte{}
	hashSource = append(hashSource, serviceToken...)
	hashSource = append(hashSource, nonce...)

	hash := sha256.Sum256(hashSource)
	expected := hex.EncodeToString(hash[:])

	u := url.URL{
		Scheme: SCHEME,
		Host:   domain,
		Path:   DOMAIN_AUTH_PATH,
	}

	resp, err := http.Get(u.String())
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		return errors.New(ERROR_ACCESS_DOMAIN_AUTH_TARGET)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "%v", err)
		return errors.New(ERROR_ACCESS_DOMAIN_AUTH_TARGET)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if string(b) != expected {
		return errors.New(ERROR_DOMAIN_AUTH_TOKEN_INVALID)
	}

	return nil
}
