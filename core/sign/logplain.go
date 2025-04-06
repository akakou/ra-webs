package sign

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
)

var Sign = sign
var Verify = verify

func sign(target *LogPlain, signKey *rsa.PrivateKey) ([]byte, error) {
	targetData, err := json.Marshal(target)
	if err != nil {
		return nil, err
	}

	hash := sha256.Sum256(targetData)

	signature, err := rsa.SignPKCS1v15(rand.Reader, signKey, crypto.SHA256, hash[:])
	if err != nil {
		return nil, err
	}
	return signature, nil
}

func verify(signature []byte, target *LogPlain, publicKey *rsa.PublicKey) error {
	targetData, err := json.Marshal(target)
	if err != nil {
		return err
	}

	hash := sha256.Sum256(targetData)

	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash[:], signature)
	if err != nil {
		return err
	}
	return nil
}
