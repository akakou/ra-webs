package core

type TAInfo struct {
	Attestation string `json:"attestation"`
	PublicKey   []byte `json:"public_key"`
	Domain      string `json:"domain"`
}
