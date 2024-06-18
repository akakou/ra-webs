package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"ta"

	goutils "github.com/akakou/go-utils"
	golangutils "github.com/akakou/golang-utils"
	"github.com/akakou/ra_webs/core"

	"github.com/akakou/ra_webs/ttp/api"
	"github.com/akakou/ra_webs/ttp/builder"
	ttpcore "github.com/akakou/ra_webs/ttp/core"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type CTForTest struct {
}

func (CTForTest) Setup(*ttpcore.TTP) error {
	return nil
}

func (CTForTest) Run(*ttpcore.TTP) {
}

func ttpForTest() (*ttpcore.TTP, error) {
	dbType := golangutils.GetEnv("DB_TYPE", "sqlite3")
	dbConfig := golangutils.GetEnv("DB_CONFIG", "file:ent?mode=memory&cache=shared&_fk=1")
	fmt.Printf("We use %s as database type and %s as database config\n", dbType, dbConfig)

	adminToken, err := goutils.RandomHex(ttpcore.RANDOM_SIZE)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ttpcore.ERROR_RANDOM_GENERATE, err)
	}

	fmt.Printf("Admin token generated: %s\n", adminToken)

	dbc := ttpcore.DBConfig{
		Type:   dbType,
		Config: dbConfig,
	}

	db, err := ttpcore.NewDB(&dbc)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ttpcore.ERROR_INIT_DB, err)
	}

	return ttpcore.NewTTP(db, &CTForTest{}, adminToken)
}
func TestAPI(t *testing.T) {
	core.EnableDebug()
	builder.EnableDebug()

	e := echo.New()
	e.Debug = true
	ttp, err := ttpForTest()
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
			&ta.TAConfig{
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
