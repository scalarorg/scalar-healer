package evm

import (
	"github.com/scalarorg/scalar-healer/pkg/db/sqlc"
	"github.com/scalarorg/scalar-healer/pkg/utils/funcs"
)

func CreateSwitchPhaseParams(cusID []byte, newPhase sqlc.RedeemPhase) []byte {
	uid := [32]byte{}
	copy(uid[:], cusID)
	return funcs.Must(SwitchPhaseArguments.Pack(newPhase.Uint8(), uid))
}
