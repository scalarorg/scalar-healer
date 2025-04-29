package models

type TransferRequest struct {
	Address            []byte `json:"address" bson:"address"`
	Signature          []byte `json:"signature" bson:"signature"` // not need to validte length
	ChainID            uint64 `json:"chain_id" bson:"chain_id"`
	DestinationChain   string `json:"destination_chain" bson:"destination_chain"`
	DestinationAddress []byte `json:"destination_address" bson:"destination_address"`
	Symbol             string `json:"symbol" bson:"symbol"`
	Amount             string `json:"amount" bson:"amount"` // bigint format
	Nonce              uint64 `json:"nonce" bson:"nonce"`

	CreatedAt int64 `json:"created_at" bson:"created_at"`
	UpdatedAt int64 `json:"updated_at" bson:"updated_at"`
}
