package core

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"

	"github.com/akakou/ra-webs/log/api/interfacestruct"
)

type SignTarget struct {
	*interfacestruct.PostRequest `json:"post_request"`
	Domain                       string `json:"domain"`
}

var Sign = sign
var Verify = verify

func sign(log *Log, req *interfacestruct.PostRequest) ([]byte, error) {
	target := SignTarget{
		PostRequest: req,
		Domain:      log.Domain,
	}

	targetData, err := json.Marshal(target)
	if err != nil {
		return nil, err
	}

	hash := sha256.Sum256(targetData)

	signature, err := rsa.SignPKCS1v15(rand.Reader, log.SignKey, crypto.SHA256, hash[:])
	if err != nil {
		return nil, err
	}
	return signature, nil
}

func verify(signature []byte, domain string, req *interfacestruct.PostRequest, publicKey *rsa.PublicKey) ([]byte, error) {
	target := SignTarget{
		PostRequest: req,
		Domain:      domain,
	}

	targetData, err := json.Marshal(target)
	if err != nil {
		return nil, err
	}

	hash := sha256.Sum256(targetData)

	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash[:], signature)
	if err != nil {
		return nil, err
	}
	return signature, nil
}
