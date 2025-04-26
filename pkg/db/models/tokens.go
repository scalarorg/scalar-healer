package models

type Token struct {
	Symbol  string `json:"symbol" bson:"symbol"`
	ChainID uint64 `json:"chain_id" bson:"chain_id"`
	Active  bool   `json:"active" bson:"active"`
	Address []byte `json:"address" bson:"address"`
	Decimal uint64 `json:"decimal" bson:"decimal"`
	Name    string `json:"name" bson:"name"`

	CreatedAt uint64 `json:"created_at" bson:"created_at"`
	UpdatedAt uint64 `json:"updated_at" bson:"updated_at"`
}
