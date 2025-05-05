package db

import (
	"math"
)

type Phase uint8

const (
	Preparing Phase = iota // EnumIndex = 0
	Executing              // EnumIndex = 1
)

type Session struct {
	Sequence uint64
	Phase    Phase
}

func (s *Session) Cmp(other *Session) int64 {
	if other == nil {
		return math.MaxInt64
	}
	var diffSeq, diffPhase int64
	if s.Sequence >= other.Sequence {
		diffSeq = int64(s.Sequence - other.Sequence)
	} else {
		diffSeq = -int64(other.Sequence - s.Sequence)
	}

	if s.Phase >= other.Phase {
		diffPhase = int64(s.Phase - other.Phase)
	} else {
		diffPhase = -int64(other.Phase - s.Phase)
	}

	return diffSeq*2 + diffPhase
}

type RedeemSession struct {
}

type RedeemRequest struct {
	Address   []byte `json:"address" bson:"address"`
	Signature []byte `json:"signature" bson:"signature"` // not need to validte length
	ChainID   uint64 `json:"chain_id" bson:"chain_id"`
	Symbol    string `json:"symbol" bson:"symbol"`
	Amount    string `json:"amount" bson:"amount"` // bigint format
	Nonce     uint64 `json:"nonce" bson:"nonce"`

	CreatedAt int64 `json:"created_at" bson:"created_at"`
	UpdatedAt int64 `json:"updated_at" bson:"updated_at"`
}

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

type Reservation struct {
	Request string `json:"request" bson:"request"`
	Amount  uint64 `json:"amount" bson:"amount"`
}

type UTXO struct {
	TxID         []byte         `json:"txid" bson:"txid"` // in reverse-order
	Vout         uint32         `json:"vout" bson:"vout"`
	ScriptPubkey []byte         `json:"script_pubkey" bson:"script_pubkey"`
	AmountInSats uint64         `json:"amount_in_sats" bson:"amount_in_sats"`
	Reservations []*Reservation `json:"reservations" bson:"reservations"`
}

type UTXOSnapshot struct {
	CustodianGroupUID []byte  `json:"custodian_group_uid" bson:"custodian_group_uid"`
	BlockHeight       uint64  `json:"block_height" bson:"block_height"`
	UTXOs             []*UTXO `json:"utxos" bson:"utxos"`
}

type GatewayAddress struct {
	Address []byte `json:"address" bson:"address"`
	ChainID uint64 `json:"chain_id" bson:"chain_id"` // TODO: index this

	CreatedAt uint64 `json:"created_at" bson:"created_at"`
	UpdatedAt uint64 `json:"updated_at" bson:"updated_at"`
}

type CustodianGroup struct {
	UID           []byte   `json:"uid" bson:"uid"`
	Name          string   `json:"name" bson:"name"`
	BitcoinPubkey []byte   `json:"bitcoin_pubkey" bson:"bitcoin_pubkey"`
	Quorum        uint32   `json:"quorum" bson:"quorum"`
	Custodians    [][]byte `json:"custodians" bson:"custodians"`
}

type BridgeRequest struct {
	Address   []byte `json:"address" bson:"address"`
	Signature []byte `json:"signature" bson:"signature"`
	ChainID   uint64 `json:"chain_id" bson:"chain_id"`
	TxHash    []byte `json:"tx_hash" bson:"tx_hash"`
	Nonce     uint64 `json:"nonce" bson:"nonce"`

	CreatedAt int64 `json:"created_at" bson:"created_at"`
	UpdatedAt int64 `json:"updated_at" bson:"updated_at"`
}
