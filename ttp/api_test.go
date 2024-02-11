package ttp

import (
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

	t.Run("TestRegisterTAApi", func(t *testing.T) {
		body := `{"domain":"example.com","public_key":"public_key", "ip":"0.0.0.0", "git":"github.com/ra-webs/ra_webs"}`

		req := httptest.NewRequest(http.MethodPost, registerTAApi.path, strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		err = registerTAApi.f(auditor)(c)
		assert.NoError(t, err)

		assert.Equal(t, "1", rec.Body.String())
	})

	t.Run("TestUpdateTAApi", func(t *testing.T) {
		path := fmt.Sprintf("/ta/%d/update", 1)
		fmt.Printf("path: %v\n", path)

		req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(""))
		req.Header.Set(echo.HeaderContentType, echo.MIMETextPlain)

		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		c.SetPath(path)
		c.SetParamNames("id")
		c.SetParamValues("1")

		err = updateTAApi.f(auditor)(c)
		assert.NoError(t, err)

		assert.Equal(t, "ok", rec.Body.String())
	})

}
