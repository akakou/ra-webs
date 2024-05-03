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

	code := ttp.DB.Client.TACode.Create().
		SetUniqueID([]byte("1234")).
		SetCommitID("1234").
		SetRepository("https://example.com").
		SetIsActive(true).
		SaveX(*ttp.DB.Ctx)

	server := ttp.DB.Client.TAServer.Create().
		SetDomain(domain).
		SetService(servicer).
		SetCode(code).
		SetPublicKey([]byte("1234")).
		SetQuote([]byte("1234")).
		SetHasActivated(true).
		SaveX(*ttp.DB.Ctx)

	violation := ttp.DB.Client.TAViolation.Create().
		SetServer(server).
		SaveX(*ttp.DB.Ctx)

	req := httptest.NewRequest(http.MethodGet, "/server/"+domain, strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/server/:domain")
	c.SetParamNames("domain")
	c.SetParamValues(domain)

	err = GetServerFromDomainApi.F(ttp)(c)

	assert.NoError(t, err)
	assert.Equal(t, 200, c.Response().Status)

	bytes := rec.Body.Bytes()
	assert.NoError(t, err)

	respTa := []ent.TAServer{}
	err = json.Unmarshal(bytes, &respTa)
	assert.NoError(t, err)

	assert.Equal(t, server.ID, respTa[0].ID)
	assert.Equal(t, server.PublicKey, respTa[0].PublicKey)
	assert.Equal(t, server.Quote, respTa[0].Quote)

	assert.Equal(t, code.ID, respTa[0].Edges.Code.ID)
	assert.Equal(t, code.UniqueID, respTa[0].Edges.Code.UniqueID)
	assert.Equal(t, code.CommitID, respTa[0].Edges.Code.CommitID)
	assert.Equal(t, code.Repository, respTa[0].Edges.Code.Repository)
	assert.Equal(t, code.IsActive, respTa[0].Edges.Code.IsActive)

	assert.Equal(t, violation.ID, respTa[0].Edges.Violation[0].ID)

	t.Errorf("%v", string(bytes))

}
