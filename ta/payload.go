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
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}

	u, err := url.Parse(rj.URL)
	if err != nil {
		return nil, fmt.Errorf("url.Parse: %w", err)
	}

	req.Method = rj.Method
	req.Header = rj.Header
	req.URL = u
	req.Body = io.NopCloser(bytes.NewBuffer(rj.Body))

	return req, nil
}

func respToJson(resp *echo.Response, reader *bufio.ReadWriter) ([]byte, error) {
	fmt.Printf("2-1 ")

	size := reader.Reader.Buffered()

	body := make([]byte, size)
	_, err := io.ReadAtLeast(reader, body, size)

	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: %w", err)
	}
	fmt.Printf("2-2 ")

	rj := respJson{
		Body:   body,
		Header: resp.Header(),
		Status: resp.Status,
	}
	fmt.Printf("2-3 ")

	j, err := json.Marshal(rj)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %w", err)
	}

	return j, nil
}

func cipherToResp(cipher *scCipher, rw *bufio.ReadWriter, resp *echo.Response) error {
	cj, err := json.Marshal(cipher)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	// rw.Discard(int(resp.Size))
	_, err = rw.Write(cj)
	if err != nil {
		return fmt.Errorf("rw.Write: %w", err)
	}

	return nil
}

func extractCipher(r *http.Request) (*scCipher, error) {
	var cipher scCipher

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: %w", err)
	}

	err = json.Unmarshal(body, &cipher)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %w", err)
	}

	return &cipher, nil
}
