package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
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
		PublicKey:  []byte("publickey_test_1"),
	},
	ent.TA{
		ID:         2,
		Repository: "test",
		CommitID:   "test",
		Evidence:   "test",
		PublicKey:  []byte("publickey_test_2"),
	},
}

var reqData = []io.PostRequest{
	&sign.LogPlain{
		Repository: storedData[0].Repository,
		CommitId:   storedData[0].CommitID,
		Evidence:   storedData[0].Evidence,
		PublicKey:  []byte("publickey_test_1"),
	},
	&sign.LogPlain{
		Repository: storedData[1].Repository,
		CommitId:   storedData[1].CommitID,
		Evidence:   storedData[1].Evidence,
		PublicKey:  []byte("publickey_test_2"),
	},
}

var max = 200

func TestAll(t *testing.T) {
	db, err := log.NewDB(&log.DBConfig{
		Type:   "sqlite3",
		Config: ":memory:?_fk=1",
	})

	assert.NoError(t, err, "DB initialization failed")

	log := &log.Log{
		DB:     db,
		Domain: "localhost",
		Token:  "token",
	}

	e := echo.New()
	g := e.Group("/")

	GetApi.Set(g, log)
	PostApi.Set(g, log)

	testPost(t, 1, reqData[0], log, e)
	testPost(t, 2, reqData[1], log, e)

	testGet(t, storedData, log, e)

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
	encPK := base64.URLEncoding.EncodeToString(data[0].PublicKey)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte{}))
	q := req.URL.Query()
	q.Set("publicKey", encPK)

	req.URL.RawQuery = q.Encode()

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := GetApi.F(log)(c)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)

	expected, err := json.Marshal(data[0])
	assert.NoError(t, err)

	actual := rec.Body.String()
	assert.Equal(t, string(expected)+"\n", actual)
}
