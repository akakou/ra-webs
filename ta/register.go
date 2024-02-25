package ta

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"strconv"
)

var TTP_URL = ""
var TTP_REGISTER_TA = TTP_URL + "ta/"
var TTP_ISSUE_CERT = TTP_URL + "ta/%d/cert"

func (ap *TA) Register() (int, error) {
	publicKey := ap.PrivateKey.Public()
	publicKeyBuf := x509.MarshalPKCS1PublicKey(publicKey.(*rsa.PublicKey))

	body := map[string]any{
		"public_key": publicKeyBuf,
		"code_id":    ap.Config.CodeId,
		"server_id":  ap.Config.ServerId,
	}

	resp, err := ap.requestToTTP(TTP_REGISTER_TA, body)
	if err != nil {
		return 0, fmt.Errorf("failed to register: %w", err)
	}

	taId, err := strconv.Atoi(string(resp))
	if err != nil {
		return 0, fmt.Errorf("failed to convert ta id: %w", err)
	}

	return taId, nil
}
