package ta

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
)

type secureChannelCipher struct {
	iv         []byte
	cipherText []byte
}

type secureChannel struct {
	key []byte
	gcm cipher.AEAD
}

type secureKeyReceiver struct {
	privateKey *rsa.PrivateKey
}

func newSecureKeyReceiver(privKey *rsa.PrivateKey) *secureKeyReceiver {
	return &secureKeyReceiver{
		privateKey: privKey,
	}
}

func (receiver *secureKeyReceiver) run(pubkeyCipher []byte) ([]byte, error) {
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

func (sc *secureChannel) decrypt(scCipher *secureChannelCipher) ([]byte, error) {
	plaintext, err := sc.gcm.Open(nil, scCipher.iv, scCipher.cipherText, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open gcm: %w", err)
	}
	return plaintext, nil
}

func (sc *secureChannel) encrypt(plainText []byte) (*secureChannelCipher, error) {
	iv := make([]byte, sc.gcm.NonceSize())
	_, err := rand.Read(iv)
	if err != nil {
		return nil, fmt.Errorf("failed to read random: %w", err)
	}

	cipherText := sc.gcm.Seal(nil, iv, plainText, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open gcm: %w", err)
	}

	scCipher := secureChannelCipher{
		iv:         iv,
		cipherText: cipherText,
	}

	return &scCipher, nil
}
