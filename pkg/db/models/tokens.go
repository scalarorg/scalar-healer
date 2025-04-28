package models

type Protocol struct {
	Name               string   `json:"name" bson:"name"`
	CustodianGroupName string   `json:"custodian_group_name" bson:"custodian_group_name"`
	CustodianGroupUid  [32]byte `json:"custodian_group_uid" bson:"custodian_group_uid"`
	Tag                string   `json:"tag" bson:"tag"`
	LiquidityModel     string   `json:"liquidity_model" bson:"liquidity_model"`
	Asset              string   `json:"asset" bson:"asset"`
	Symbol             string   `json:"symbol" bson:"symbol"`
	Decimals           uint8    `json:"decimals" bson:"decimals"`
	Capacity           uint64   `json:"capacity" bson:"capacity"`
	DailyMintLimit     uint64   `json:"daily_mint_limit" bson:"daily_mint_limit"`
	Avatar             string   `json:"avatar" bson:"avatar"`
}
type Token struct {
	Protocol  string `json:"protocol" bson:"protocol"`
	Symbol    string `json:"symbol" bson:"symbol"`
	ChainID   uint64 `json:"chain_id" bson:"chain_id"`
	Active    bool   `json:"active" bson:"active"`
	Address   []byte `json:"address" bson:"address"`
	Decimal   uint64 `json:"decimal" bson:"decimal"`
	Name      string `json:"name" bson:"name"`
	Avatar    string `json:"avatar" bson:"avatar"`
	CreatedAt uint64 `json:"created_at" bson:"created_at"`
	UpdatedAt uint64 `json:"updated_at" bson:"updated_at"`
}
