package eip712

type BaseRequest struct {
	Address   string `json:"address" validate:"eth_addr"`
	Signature string `json:"signature" validate:"hexadecimal"`
	Nonce     uint64 `json:"nonce" validate:"gte=0"`
	ChainID   uint64 `json:"chain_id" validate:"required,gte=0"`
}
