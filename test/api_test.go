package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"service"

	"github.com/akakou/ra_webs/ttp/api"
	ttpcore "github.com/akakou/ra_webs/ttp/core"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var PORT = ":11111"

func TestAPI(t *testing.T) {
	e := echo.New()
	e.Debug = true
	ttp, err := ttpcore.DefaultTTP()
	assert.NoError(t, err)

	api.Route(e, ttp)

	go e.Start(PORT)

	token := ""

	t.Run("TestPostServiceByAdmin", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, api.PostServiceByAdmin.Path, strings.NewReader(""))
		req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", ttp.AdminToken))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		err = api.PostServiceByAdmin.F(ttp)(c)
		assert.NoError(t, err)
		assert.Equal(t, rec.Result().StatusCode, http.StatusOK)

		token = rec.Body.String()
	})

	svc := service.NewService(token, "http://localhost"+PORT+"/")
	api.SCHEME = "http"

	t.Run("TestPostTACode", func(t *testing.T) {
		repo := "https://github.com/akakou-docs/ego-statistical-analysis"
		uniqueId, err := svc.PostCode(repo)
		assert.NoError(t, err)
		assert.Equal(t, "c585bdc800065a3de7c372c8fc6f1259154a85ef9d76b671caf80551082679ba", uniqueId)
	})

	t.Run("TestPostTAServer", func(t *testing.T) {
		e := echo.New()
		e.Debug = true

		res, err := svc.PostServer(e, "localhost", ":8080")
		assert.NoError(t, err)
		assert.Equal(t, "1", res)
	})
}
