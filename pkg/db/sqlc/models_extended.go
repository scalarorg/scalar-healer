package sqlc

import (
	"math"

	"github.com/ethereum/go-ethereum/crypto"
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

type CommandStatus uint8

const (
	COMMAND_STATUS_PENDING  CommandStatus = 0
	COMMAND_STATUS_QUEUED   CommandStatus = 1
	COMMAND_STATUS_EXECUTED CommandStatus = 2
)

func (s CommandStatus) Int32() int32 {
	return int32(s)
}

type CommandBatchStatus uint8

const (
	COMMAND_BATCH_STATUS_PENDING  CommandBatchStatus = 0
	COMMAND_BATCH_STATUS_EXECUTED CommandBatchStatus = 1
)

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
