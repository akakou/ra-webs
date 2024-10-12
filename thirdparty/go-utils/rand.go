package goutils

import (
	"crypto/rand"
	"encoding/hex"
)

func RandomHex(size int) (string, error) {
	buf := make([]byte, size)

	_, err := rand.Read(buf)
	if err != nil {
		return "", nil
	}

	r := hex.EncodeToString(buf)

	return r, nil
}
