package api

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	ttpcore "github.com/akakou/ra_webs/ttp/core"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func setupMockServer(handler http.HandlerFunc, t *testing.T) (*httptest.Server, *url.URL) {
	h := http.HandlerFunc(handler)

	ts := httptest.NewServer(h)
	u, err := url.Parse(ts.URL)
	assert.NoError(t, err)

	return ts, u
}

func TestAPI(t *testing.T) {
	e := echo.New()
	e.Debug = true
	ttp, err := ttpcore.DefaultTTP()
	assert.NoError(t, err)

	token := ""

	t.Run("TestPostServiceByAdmin", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, postServiceByAdmin.Path, strings.NewReader(""))
		req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", ttp.AdminToken))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		err = postServiceByAdmin.F(ttp)(c)
		assert.NoError(t, err)
		assert.Equal(t, rec.Result().StatusCode, http.StatusOK)

		token = rec.Body.String()
	})

	t.Run("TestPostTACode", func(t *testing.T) {
		body := `{"repository":"https://github.com/akakou-docs/ego-statistical-analysis"}`

		req := httptest.NewRequest(http.MethodPost, postCodeApi.Path, strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", token))

		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		err = postCodeApi.F(ttp)(c)
		assert.NoError(t, err)
		assert.Equal(t, rec.Result().StatusCode, http.StatusOK)

		assert.Equal(t, "c585bdc800065a3de7c372c8fc6f1259154a85ef9d76b671caf80551082679ba", rec.Body.String())
	})

	t.Run("TestPostTAServer", func(t *testing.T) {
		nonce := []byte("aaaaa")
		SCHEME = "http"

		hashSource := []byte{}
		hashSource = append(hashSource, []byte(token)...)
		hashSource = append(hashSource, []byte(nonce)...)

		hash := sha256.Sum256(hashSource)
		serverToken := hex.EncodeToString(hash[:])

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, serverToken)
		})

		_, u := setupMockServer(handler, t)

		body := fmt.Sprintf(`{"domain": "%s", "nonce": "aaaaa"}`, u.Host)

		req := httptest.NewRequest(http.MethodPost, postServerApi.Path, strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", token))

		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		err = postServerApi.F(ttp)(c)
		assert.NoError(t, err)
		assert.Equal(t, rec.Result().StatusCode, http.StatusOK)

		assert.NoError(t, err)

		assert.Equal(t, "1", rec.Body.String())
	})

}
