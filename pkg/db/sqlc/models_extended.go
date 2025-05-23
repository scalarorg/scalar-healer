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

func (s *RedeemSession) Cmp(other *RedeemSession) int64 {
	if other == nil {
		return math.MaxInt64
	}
	var diffSeq, diffPhase int64
	if s.Sequence >= other.Sequence {
		diffSeq = s.Sequence - other.Sequence
	} else {
		diffSeq = -int64(other.Sequence - s.Sequence)
	}

	if s.CurrentPhase.Uint8() >= other.CurrentPhase.Uint8() {
		diffPhase = int64(s.CurrentPhase.Uint8() - other.CurrentPhase.Uint8())

	} else {
		diffPhase = -int64(other.CurrentPhase.Uint8() - s.CurrentPhase.Uint8())
	}

	return diffSeq*2 + diffPhase
}
