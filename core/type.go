package core

type ProvisioningRequest struct {
	Attestation string `json:"attestation"`
	PublicKey   string `json:"public_key"`
	Domain      string `json:"domain"`
}
