package api

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/akakou/ra-webs/core/sign"
	"github.com/akakou/ra-webs/log"
	core "github.com/akakou/ra-webs/log"
	"github.com/akakou/ra-webs/log/api/io"
	"github.com/akakou/ra-webs/log/ent"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var storedData = []ent.TA{
	ent.TA{
		ID:         1,
		Repository: "test",
		CommitID:   "test",
		Evidence:   "test",
		Signature:  []byte("hello"),
	},
	ent.TA{
		ID:         2,
		Repository: "test",
		CommitID:   "test",
		Evidence:   "test",
		Signature:  []byte("hello"),
	},
}

var reqData = []io.PostRequest{
	&sign.LogPlain{
		Repository: storedData[0].Repository,
		CommitId:   storedData[0].CommitID,
		Evidence:   storedData[0].Evidence,
	},
	&sign.LogPlain{
		Repository: storedData[1].Repository,
		CommitId:   storedData[1].CommitID,
		Evidence:   storedData[1].Evidence,
	},
}

var max = 200

func TestAll(t *testing.T) {
	db, err := log.NewDB(&log.DBConfig{
		Type:   "sqlite3",
		Config: ":memory:?_fk=1",
	})

	assert.NoError(t, err, "DB initialization failed")

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err, "Key generation failed")

	log := &log.Log{
		DB:        db,
		Domain:    "localhost",
		VerifyKey: &privateKey.PublicKey,
		SignKey:   privateKey,
		Token:     "token",
	}

	e := echo.New()
	g := e.Group("/")

	sign.Sign = func(target *sign.LogPlain, signKey *rsa.PrivateKey) ([]byte, error) {
		return []byte("hello"), nil
	}

	GetApi.Set(g, log)
	PostApi.Set(g, log)

	testPost(t, 1, reqData[0], log, e)
	testPost(t, 2, reqData[1], log, e)

	testGet(t, storedData, log, e)

	for i := 0; i < max; i++ {
		log.DB.Client.TA.Create().
			SetRepository("test").
			SetCommitID("test").
			SetEvidence("test").
			SetSignature([]byte("hello")).
			SaveX(*log.DB.Ctx)
	}

	testGetWithStart(t, log, e)
	testGetWithStartAndEnd(t, log, e)

	defer log.DB.Close()
}

func testPost(t *testing.T, counter int, data *sign.LogPlain, log *log.Log, e *echo.Echo) {
	reqJson, err := json.Marshal(data)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(reqJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	req.Header.Set("Authorization", "Bearer "+log.Token)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = PostApi.F(log)(c)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)

	id := rec.Body.String()
	assert.Equal(t, fmt.Sprintf("%d", counter), id)
}

func testGet(t *testing.T, data []ent.TA, log *core.Log, e *echo.Echo) {
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte{}))

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := GetApi.F(log)(c)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)

	expected, err := json.Marshal(data)
	assert.NoError(t, err)

	actual := rec.Body.String()
	assert.Equal(t, string(expected)+"\n", actual)
}

func testGetWithStart(t *testing.T, log *core.Log, e *echo.Echo) {
	var data io.GetResponse

	rec := httptest.NewRecorder()

	q := make(url.Values)
	q.Set("start", "25")
	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)

	c := e.NewContext(req, rec)

	err := GetApi.F(log)(c)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)

	actual := rec.Body.String()
	json.Unmarshal([]byte(actual), &data)

	fmt.Printf("data: %v\n", data)

	assert.Equal(t, 25, data[0].ID)
	assert.Equal(t, 24+100, data[len(data)-1].ID)
	assert.Equal(t, 100, len(data))
}

func testGetWithStartAndEnd(t *testing.T, log *core.Log, e *echo.Echo) {
	var data io.GetResponse

	rec := httptest.NewRecorder()

	q := make(url.Values)
	q.Set("start", "25")
	q.Set("end", "100")
	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)

	c := e.NewContext(req, rec)

	err := GetApi.F(log)(c)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)

	actual := rec.Body.String()
	json.Unmarshal([]byte(actual), &data)

	fmt.Printf("data: %v\n", data)

	assert.Equal(t, 25, data[0].ID)
	assert.Equal(t, 100, data[len(data)-1].ID)
	assert.Equal(t, 100-24, len(data))
}
