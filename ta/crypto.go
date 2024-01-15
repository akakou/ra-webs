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
	IV   []byte `json:"iv"`
	Text []byte `json:"text"`
	Key  []byte `json:"key"`
}

type secureChannel struct {
	key []byte
	gcm cipher.AEAD
}

type scProvisioner struct {
	privateKey *rsa.PrivateKey
}

func newSCProvisioner(privKey *rsa.PrivateKey) *scProvisioner {
	return &scProvisioner{
		privateKey: privKey,
	}
}

func (provisioner *scProvisioner) decryptKey(pubkeyCipher []byte) ([]byte, error) {
	comKey, err := rsa.DecryptOAEP(sha256.New(), nil, provisioner.privateKey, pubkeyCipher, []byte{})
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt com key: %w", err)
	}

	return comKey, nil
}

func (provisioner *scProvisioner) provision(pubkeyCipher []byte) (*secureChannel, error) {
	key, err := provisioner.decryptKey(pubkeyCipher)
	if err != nil {
		return nil, err
	}

	secureChannel, err := newSecureChannel(key)
	if err != nil {
		return nil, err
	}

	return secureChannel, nil
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
	plaintext, err := sc.gcm.Open(nil, scCipher.IV, scCipher.Text, nil)
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
		IV:   iv,
		Text: cipher,
	}

	return &scCipher, nil
}
