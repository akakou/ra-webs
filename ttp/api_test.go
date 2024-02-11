package ttp

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAPI(t *testing.T) {
	e := echo.New()
	e.Debug = true
	auditor, err := DefaultAuditor()
	assert.NoError(t, err)

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)
	publicKey := privateKey.PublicKey

	publicKeyBuf := x509.MarshalPKCS1PublicKey(&publicKey)

	t.Run("TestRegisterTAApi", func(t *testing.T) {
		body := `{"domain":"example.com","public_key":"public_key", "ip":"0.0.0.0", "git":"github.com/ra-webs/ra_webs"}`

		req := httptest.NewRequest(http.MethodPost, registerTAApi.path, strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		err = registerTAApi.f(auditor)(c)
		assert.NoError(t, err)
		assert.Equal(t, rec.Result().StatusCode, http.StatusOK)

		assert.Equal(t, "1", rec.Body.String())
	})

	t.Run("TestUpdateTAApi", func(t *testing.T) {
		body, err := json.Marshal(map[string]interface{}{"public_key": publicKeyBuf})
		assert.NoError(t, err)

		path := fmt.Sprintf("/ta/%d/update", 1)
		fmt.Printf("path: %v\n", path)

		req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(string(body)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		c.SetPath(path)
		c.SetParamNames("id")
		c.SetParamValues("1")

		err = updateTAApi.f(auditor)(c)
		assert.NoError(t, err)

		assert.Equal(t, rec.Result().StatusCode, http.StatusOK)
		fmt.Printf("%v", rec.Body.String())

	})

}
