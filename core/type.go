package core

type ProvisioningRequest struct {
	Attestation string `json:"attestation"`
	PublicKey   []byte `json:"public_key"`
	Domain      string `json:"domain"`
}
