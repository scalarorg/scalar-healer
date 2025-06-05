package evm

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/scalarorg/scalar-healer/pkg/db/sqlc"
	"github.com/scalarorg/scalar-healer/pkg/utils/funcs"
	"github.com/scalarorg/scalar-healer/pkg/utils/slices"
)

func CreateSwitchPhaseParams(cusID []byte, newPhase sqlc.RedeemPhase) []byte {
	uid := [32]byte{}
	copy(uid[:], cusID)
	return funcs.Must(SwitchPhaseArguments.Pack(newPhase.Uint8(), uid))
}

func PackArguments(chainID *big.Int, commandIDs []sqlc.CommandID, commands []sqlc.CommandType, commandParams [][]byte) ([]byte, error) {
	if len(commandIDs) != len(commands) || len(commandIDs) != len(commandParams) {
		return nil, fmt.Errorf("length mismatch for command arguments")
	}

	arguments := abi.Arguments{{Type: uint256Type}, {Type: bytes32ArrayType}, {Type: stringArrayType}, {Type: bytesArrayType}}
	result, err := arguments.Pack(
		chainID,
		commandIDs,
		slices.Map(commands, sqlc.CommandType.String),
		commandParams,
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}
