package sqlc

import (
	"math"
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


type CommandBatchStatus uint8

const (
	COMMAND_BATCH_STATUS_PENDING  CommandBatchStatus = 0
	COMMAND_BATCH_STATUS_EXECUTED CommandBatchStatus = 1
)