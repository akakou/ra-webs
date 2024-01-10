package core

type TAInfo struct {
	Attestation   string `json:"attestation"`
	PublicKeyHash string `json:"public_key_hash"`
	Domain        string `json:"domain"`
}
