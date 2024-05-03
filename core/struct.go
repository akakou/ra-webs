package core

type CodeRequest struct {
	Repository string `json:"repository"`
}

type ServerRequest struct {
	PublicKey []byte `json:"public_key"`
	Domain    string `json:"domain"`
	Quote     string `json:"quote"`
}

type RegisterRequest struct {
	CodeRequest   `json:"code"`
	ServerRequest `json:"server"`
}
