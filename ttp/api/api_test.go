package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	ttpcore "github.com/akakou/ra_webs/ttp/core"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

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
		body := `{"ip":"0.0.0.0", "domain": "example.com"}`

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

	t.Run("TestActivateServer", func(t *testing.T) {
		path := "/server/1/activate"

		req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(""))
		req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", ttp.AdminToken))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		c.SetPath(path)
		c.SetParamNames("id")
		c.SetParamValues("1")
		err = postActivateServerApi.F(ttp)(c)
		assert.NoError(t, err)
		assert.Equal(t, rec.Result().StatusCode, http.StatusOK)

		assert.Equal(t, "1", rec.Body.String())
	})

}
