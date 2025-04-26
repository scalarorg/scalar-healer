package models

type GatewayAddress struct {
	Address []byte `json:"address" bson:"address"`
	ChainID uint64 `json:"chain_id" bson:"chain_id"` // TODO: index this

	CreatedAt uint64 `json:"created_at" bson:"created_at"`
	UpdatedAt uint64 `json:"updated_at" bson:"updated_at"`
}
