package ta

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
)

func decryptMiddleware(c echo.Context, provisioner scProvisioner) (*secureChannel, error) {
	r := c.Request()

	if r.URL.Path == "/" {
		return nil, nil
	}

	cipher, err := extractCipher(r)
	if err != nil {
		return nil, fmt.Errorf("extractCipher: %w", err)

	}

	sc, err := provisioner.provision(cipher.Key)
	if err != nil {
		return nil, fmt.Errorf("provision: %w", err)
	}

	plain, err := sc.decrypt(cipher)
	if err != nil {
		return nil, fmt.Errorf("decrypt: %w", err)
	}

	req, err := reqFromJson(plain, r)
	if err != nil {
		return nil, fmt.Errorf("reqFromJson: %w", err)
	}

	c.SetRequest(req)

	return sc, nil
}

func encryptMiddlware(c echo.Context, sc *secureChannel) {
	log.Println("after action")

	resp := c.Response()

	conn, rw, err := resp.Hijack()
	if err != nil {
		rw.Discard(int(resp.Size))
		return
	}
	defer conn.Close()

	jsonResp, err := respToJson(resp, rw)
	if err != nil {
		rw.Discard(int(resp.Size))
		return
	}

	cipher, err := sc.encrypt(jsonResp)
	if err != nil {
		rw.Discard(int(resp.Size))
		return
	}

	cipherText, err := json.Marshal(cipher)
	if err != nil {
		rw.Discard(int(resp.Size))
		return
	}

	rw.Discard(int(resp.Size))
	rw.Write(cipherText)
}
