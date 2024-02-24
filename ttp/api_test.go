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

	t.Run("TestPostTACode", func(t *testing.T) {
		body := `{"repository":"github.com/ra-webs/ra_webs", "commit_id": "1111111111", "unique_id": "aGVsbG8K"}`

		req := httptest.NewRequest(http.MethodPost, postCodeApi.path, strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		err = postCodeApi.f(auditor)(c)
		assert.NoError(t, err)
		assert.Equal(t, rec.Result().StatusCode, http.StatusOK)

		assert.Equal(t, "1", rec.Body.String())
	})

	t.Run("TestPostTAServer", func(t *testing.T) {
		body := `{"ip":"0.0.0.0", "domain": "example.com"}`

		req := httptest.NewRequest(http.MethodPost, postServerApi.path, strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		err = postServerApi.f(auditor)(c)
		assert.NoError(t, err)
		assert.Equal(t, rec.Result().StatusCode, http.StatusOK)

		assert.Equal(t, "1", rec.Body.String())
	})

	t.Run("TestPostTA", func(t *testing.T) {
		body, err := json.Marshal(map[string]interface{}{
			"public_key": publicKeyBuf,
			"code_id":    1,
			"server_id":  1,
		})
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, postTAApi.path, strings.NewReader(string(body)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		c.SetPath(postTAApi.path)
		err = postTAApi.f(auditor)(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
		fmt.Printf("%v", rec.Body.String())
	})

}
