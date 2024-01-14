package ta

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type reqCore struct {
	Method string              `json:"method"`
	Body   []byte              `json:"body"`
	Header map[string][]string `json:"header"`
	URL    string              `json:"url"`
}

func reqFromJson(source []byte, req *http.Request) (*http.Request, error) {
	rc := reqCore{}
	err := json.Unmarshal(source, &rc)
	if err != nil {
		return nil, fmt.Errorf("aa ")
	}

	u, err := url.Parse(rc.URL)
	if err != nil {
		return nil, fmt.Errorf("aa ")
	}

	req.Method = rc.Method
	req.Header = rc.Header
	req.URL = u
	req.Body = io.NopCloser(bytes.NewBuffer(rc.Body))

	return req, nil
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

// func decryptCipher(c echo.Context, ) (*http.Request, error) {
// 	key, err := decryptor.receive(cipher.key)
// 	if err != nil {
// 		return nil, fmt.Errorf("aa")
// 	}

// 	sc, err := newAESSecureChannel(key)
// 	if err != nil {
// 		return nil, fmt.Errorf("aa")
// 	}

// 	plainText, err := sc.decrypt(cipher.content)
// 	if err != nil {
// 		return nil, fmt.Errorf("aa")
// 	}

// 	err = json.Unmarshal(plainText, &req)
// 	if err != nil {
// 		return nil, fmt.Errorf("aa")
// 	}

// 	u, err := url.Parse(core.URL)
// 	if err != nil {
// 		return nil, fmt.Errorf("aa")
// 	}

// 	body := io.NopCloser(bytes.NewReader(core.Body))

// 	req.Method = core.Method
// 	req.URL = u
// 	req.Body = body
// 	req.Header = core.Header
// 	req.ContentLength = int64(len(core.Body))

// 	return req, nil
// }
