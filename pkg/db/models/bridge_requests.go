package models

type BridgeRequest struct {
	Address   []byte `json:"address" bson:"address"`
	Signature []byte `json:"signature" bson:"signature"`
	ChainID   uint64 `json:"chain_id" bson:"chain_id"`
	TxHash    []byte `json:"tx_hash" bson:"tx_hash"`
	Nonce     uint64 `json:"nonce" bson:"nonce"`

	CreatedAt int64 `json:"created_at" bson:"created_at"`
	UpdatedAt int64 `json:"updated_at" bson:"updated_at"`
}
