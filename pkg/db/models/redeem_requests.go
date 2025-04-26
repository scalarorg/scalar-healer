package models

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
