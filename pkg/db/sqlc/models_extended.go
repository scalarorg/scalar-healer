package sqlc

import (
	"encoding/json"
	"math"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/jackc/pgx/v5/pgtype"
)

func (e RedeemPhase) Uint8() uint8 {
	switch e {
	case RedeemPhasePREPARING:
		return 0
	case RedeemPhaseEXECUTING:
		return 1
	default:
		return math.MaxUint8
	}
}

func (e RedeemPhase) Bytes() byte {
	return byte(e.Uint8())
}

func PhaseFromUint8(phase uint8) RedeemPhase {
	switch phase {
	case 0:
		return RedeemPhasePREPARING
	case 1:
		return RedeemPhaseEXECUTING
	default:
		panic("invalid phase")
	}
}

func (s *ChainRedeemSession) Cmp(other *ChainRedeemSession) int64 {
	if other == nil {
		return math.MaxInt64
	}
	diffSeq := s.Sequence - other.Sequence

	if diffSeq != 0 {
		return diffSeq
	}

	return int64(s.CurrentPhase.Uint8()) - int64(other.CurrentPhase.Uint8())
}

func (c CommandType) String() string {
	return string(c)
}

type CommandStatus uint8

const (
	COMMAND_STATUS_PENDING  CommandStatus = 0
	COMMAND_STATUS_QUEUED   CommandStatus = 1
	COMMAND_STATUS_EXECUTED CommandStatus = 2
)

func (s CommandStatus) Int32() int32 {
	return int32(s)
}

func (c CommandStatus) ToPgType() pgtype.Int4 {
	return pgtype.Int4{
		Int32: c.Int32(),
		Valid: true,
	}
}

type CommandBatchStatus uint8

const (
	COMMAND_BATCH_STATUS_PENDING  CommandBatchStatus = 0
	COMMAND_BATCH_STATUS_EXECUTED CommandBatchStatus = 1
)

func (c CommandBatchStatus) ToPgType() pgtype.Int4 {
	return pgtype.Int4{
		Int32: c.Int32(),
		Valid: true,
	}
}

func (s CommandBatchStatus) Int32() int32 {
	return int32(s)
}

type ChainRedeemSessionUpdate struct {
	Chain             string
	CustodianGroupUid []byte
	Sequence          int64
	CurrentPhase      RedeemPhase
	NewPhase          RedeemPhase
}

const commandIDSize = 32

type CommandID [commandIDSize]byte

func NewCommandID(data []byte, chainID string) CommandID {
	var commandID CommandID
	copy(commandID[:], crypto.Keccak256(append(data, []byte(chainID)...))[:commandIDSize])

	return commandID
}

func (c CommandID) Bytes() []byte {
	return c[:]
}

func (c *CommandBatch) GetExtraData() ([][]byte, error) {
	var extraData [][]byte
	err := json.Unmarshal(c.ExtraData, &extraData)
	if err != nil {
		return nil, err
	}
	return extraData, nil
}

type TokenDetails struct {
	Address []byte `json:"address"`
	ChainID int64  `json:"chain_id"`
}

type ProtocolWithTokenDetails struct {
	*Protocol
	Tokens []TokenDetails `json:"token_details"`
}
