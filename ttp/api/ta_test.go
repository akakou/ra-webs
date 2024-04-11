package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	ttpcore "github.com/akakou/ra_webs/ttp/core"
	"github.com/akakou/ra_webs/ttp/ent"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestTAFromDomainAPI(t *testing.T) {
	e := echo.New()
	e.Debug = true
	ttp, err := ttpcore.DefaultTTP()
	assert.NoError(t, err)

	domain := "example.com"

	servicer := ttp.DB.Client.Service.Create().
		SetName("").
		SetToken("").
		SetIsActive(true).
		SaveX(*ttp.DB.Ctx)

	server := ttp.DB.Client.TAServer.Create().
		SetDomain(domain).
		SetIsActive(true).
		SetService(servicer).
		SaveX(*ttp.DB.Ctx)

	code := ttp.DB.Client.TACode.Create().
		SetUniqueID([]byte("1234")).
		SetCommitID("1234").
		SetRepository("https://example.com").
		SetIsActive(true).
		SaveX(*ttp.DB.Ctx)

	ta := ttp.DB.Client.TA.Create().
		SetCode(code).
		SetServer(server).
		SetPublicKey([]byte("1234")).
		SetQuote([]byte("1234")).
		SetIsValid(true).
		SaveX(*ttp.DB.Ctx)

	req := httptest.NewRequest(http.MethodGet, "/ta/"+domain, strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/ta/:domain")
	c.SetParamNames("domain")
	c.SetParamValues(domain)

	err = GetTAFromDomainApi.F(ttp)(c)

	assert.NoError(t, err)
	assert.Equal(t, 200, c.Response().Status)

	bytes := rec.Body.Bytes()
	assert.NoError(t, err)

	respTa := ent.TA{}
	err = json.Unmarshal(bytes, &respTa)
	assert.NoError(t, err)

	assert.Equal(t, ta.ID, respTa.ID)
	assert.Equal(t, ta.PublicKey, respTa.PublicKey)
	assert.Equal(t, ta.Quote, respTa.Quote)
	assert.Equal(t, ta.IsValid, respTa.IsValid)

	assert.Equal(t, code.ID, respTa.Edges.Code.ID)
	assert.Equal(t, code.UniqueID, respTa.Edges.Code.UniqueID)
	assert.Equal(t, code.CommitID, respTa.Edges.Code.CommitID)
	assert.Equal(t, code.Repository, respTa.Edges.Code.Repository)
	assert.Equal(t, code.IsActive, respTa.Edges.Code.IsActive)

	t.Errorf("%v", string(bytes))

}
