package ta

import (
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
)

func decryptMiddleware(c echo.Context, provisioner scProvisioner) (*secureChannel, error) {
	r := c.Request()
	fmt.Printf("before method: %v", c.Request().Method)

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

	fmt.Printf("\nplain: %v\n", string(plain))

	req, err := reqFromJson(plain, r)
	if err != nil {
		return nil, fmt.Errorf("reqFromJson: %w", err)
	}

	c.SetRequest(req)
	fmt.Printf("after method: %v", c.Request().Method)

	return sc, nil
}

func encryptMiddlware(c echo.Context, sc *secureChannel) error {
	log.Println("after action")

	resp := c.Response()

	conn, rw, err := resp.Hijack()
	if err != nil {
		return fmt.Errorf("resp.Hijack: %w", err)
	}
	fmt.Printf("1 ")

	defer conn.Close()

	fmt.Printf("2 ")

	jsonResp, err := respToJson(resp, rw)
	if err != nil {
		rw.Discard(int(resp.Size))
		return fmt.Errorf("respToJson: %w", err)
	}
	fmt.Printf("3 ")

	cipher, err := sc.encrypt(jsonResp)
	if err != nil {
		rw.Discard(int(resp.Size))
		return fmt.Errorf("encrypt: %w", err)
	}
	fmt.Printf("4 ")

	// rw.Discard(int(resp.Size))

	cipherToResp(cipher, rw, resp)
	fmt.Printf("5 ")
	rw.Flush()

	fmt.Printf("after action end")
	return nil
}
