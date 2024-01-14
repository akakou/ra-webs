package ta

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
)

type reqJson struct {
	Method string              `json:"method"`
	Body   []byte              `json:"body"`
	Header map[string][]string `json:"header"`
	URL    string              `json:"url"`
}

type respJson struct {
	Body   []byte              `json:"body"`
	Header map[string][]string `json:"header"`
	Status int                 `json:"status_code"`
}

func reqFromJson(source []byte, req *http.Request) (*http.Request, error) {
	rj := reqJson{}
	err := json.Unmarshal(source, &rj)
	if err != nil {
		return nil, fmt.Errorf("aa ")
	}

	u, err := url.Parse(rj.URL)
	if err != nil {
		return nil, fmt.Errorf("aa ")
	}

	req.Method = rj.Method
	req.Header = rj.Header
	req.URL = u
	req.Body = io.NopCloser(bytes.NewBuffer(rj.Body))

	return req, nil
}

func respToJson(resp *echo.Response, reader *bufio.ReadWriter) ([]byte, error) {
	body, err := io.ReadAll(reader.Reader)
	if err != nil {
		return nil, err
	}

	rj := respJson{
		Body:   body,
		Header: resp.Header(),
		Status: resp.Status,
	}

	j, err := json.Marshal(rj)
	if err != nil {
		return nil, fmt.Errorf("aa")
	}

	return j, nil
}

func cipherToResp(cipher *scCipher, resp *echo.Response) (*echo.Response, error) {
	cj, err := json.Marshal(cipher)
	if err != nil {
		return nil, fmt.Errorf("aa")
	}

	rj := echo.Response{
		Status: 200,
	}

	rj.Write(cj)

	return &rj, nil
}

func extractCipher(r *http.Request) (*scCipher, error) {
	var cipher *scCipher

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, cipher)
	if err != nil {
		return nil, err
	}

	return cipher, nil
}
