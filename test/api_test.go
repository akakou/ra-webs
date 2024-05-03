package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"ta"

	"github.com/akakou/ra_webs/core"

	"github.com/akakou/ra_webs/ttp/api"
	"github.com/akakou/ra_webs/ttp/builder"
	ttpcore "github.com/akakou/ra_webs/ttp/core"
	"github.com/akakou/ra_webs/ttp/ct"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAPI(t *testing.T) {
	core.EnableDebug()
	ct.EnableDebug()
	builder.EnableDebug()

	e := echo.New()
	e.Debug = true
	ttp, err := ttpcore.DefaultTTP()
	assert.NoError(t, err)

	api.Route(e, ttp)

	go e.Start(core.TTPPort)

	var token = ""

	t.Run("TestPostServiceByAdmin", func(t *testing.T) {
		fmt.Println("TestPostServiceByAdmin")
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

	t.Run("TestRegister", func(t *testing.T) {
		_ta, _ := ta.InitTA(
			&ta.Config{
				Token:      token,
				Domain:     "localhost",
				TTP:        "http://localhost" + core.TTPPort,
				Repository: "https://github.com/akakou-docs/ego-statistical-analysis",
			},
		)

		_, err := _ta.Register()
		assert.NoError(t, err)
	})
}
