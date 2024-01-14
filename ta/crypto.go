package ta

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
)

type scCipher struct {
	iv         []byte
	cipherText []byte
}

type secureChannel struct {
	key []byte
	gcm cipher.AEAD
}

type scKeyDecryptor struct {
	privateKey *rsa.PrivateKey
}

func newscKeyDecryptor(privKey *rsa.PrivateKey) *scKeyDecryptor {
	return &scKeyDecryptor{
		privateKey: privKey,
	}
}

func (receiver *scKeyDecryptor) decrypt(pubkeyCipher []byte) ([]byte, error) {
	comKey, err := rsa.DecryptOAEP(sha256.New(), nil, receiver.privateKey, pubkeyCipher, []byte{})
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt com key: %w", err)
	}

	return comKey, nil
}

func newSecureChannel(key []byte) (*secureChannel, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create gcm: %w", err)
	}

	sc := secureChannel{
		key: key,
		gcm: gcm,
	}

	return &sc, nil
}

func (sc *secureChannel) decrypt(scCipher *scCipher) ([]byte, error) {
	plaintext, err := sc.gcm.Open(nil, scCipher.iv, scCipher.cipherText, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open gcm: %w", err)
	}
	return plaintext, nil
}

func (sc *secureChannel) encrypt(plainText []byte) (*scCipher, error) {
	iv := make([]byte, sc.gcm.NonceSize())
	_, err := rand.Read(iv)
	if err != nil {
		return nil, fmt.Errorf("failed to read random: %w", err)
	}

	cipher := sc.gcm.Seal(nil, iv, plainText, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open gcm: %w", err)
	}

	scCipher := scCipher{
		iv:         iv,
		cipherText: cipher,
	}

	return &scCipher, nil
}
