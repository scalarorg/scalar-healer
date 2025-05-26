package evm

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/rs/zerolog/log"
	contracts "github.com/scalarorg/scalar-healer/pkg/evm/contracts/generated"
)

var (
	addressType, _      = abi.NewType("address", "address", nil)
	stringType, _       = abi.NewType("string", "string", nil)
	bytesType, _        = abi.NewType("bytes", "bytes", nil)
	bytes32Type, _      = abi.NewType("bytes32", "bytes32", nil)
	uint8Type, _        = abi.NewType("uint8", "uint8", nil)
	uint256Type, _      = abi.NewType("uint256", "uint256", nil)
	uint64Type, _       = abi.NewType("uint64", "uint64", nil)
	bytesArrayType, _   = abi.NewType("bytes[]", "bytes[]", nil)
	uint256ArrayType, _ = abi.NewType("uint256[]", "uint256[]", nil)
	uint32ArrayType, _  = abi.NewType("uint32[]", "uint32[]", nil)
	uint64ArrayType, _  = abi.NewType("uint64[]", "uint64[]", nil)
	bytes32ArrayType, _ = abi.NewType("bytes32[]", "bytes32[]", nil)
	stringArrayType, _  = abi.NewType("string[]", "string[]", nil)
	addressArrayType, _ = abi.NewType("address[]", "address[]", nil)

	RedeemTokenPayloadArguments = abi.Arguments{{Type: uint64Type}, {Type: bytesType}, {Type: stringArrayType}, {Type: uint32ArrayType}, {Type: uint64ArrayType}, {Type: bytes32Type}}
	RedeemTokenArguments        = abi.Arguments{{Type: stringType}, {Type: stringType}, {Type: bytesType}, {Type: stringType}, {Type: uint256Type}, {Type: bytes32Type}, {Type: uint64Type}}

	SwitchPhaseArguments = abi.Arguments{{Type: uint8Type}, {Type: bytes32Type}}
)

var (
	scalarGatewayAbi *abi.ABI
	mapEvents        = map[string]*abi.Event{}
)

func init() {
	log.Info().Msg("Initializing ABI")
	scalarGatewayAbi, _ = contracts.IScalarGatewayMetaData.GetAbi()
	for _, event := range scalarGatewayAbi.Events {
		mapEvents[event.Name] = &event
	}
}

func GetScalarGatewayAbi() (*abi.ABI, error) {
	if scalarGatewayAbi == nil {
		var err error
		scalarGatewayAbi, err = contracts.IScalarGatewayMetaData.GetAbi()
		if err != nil {
			return nil, err
		}
	}
	return scalarGatewayAbi, nil
}

func GetEventIndexedArguments(eventName string) (abi.Arguments, error) {
	gatewayAbi, err := GetScalarGatewayAbi()
	if err != nil {
		return nil, err
	}
	var args abi.Arguments
	if event, ok := gatewayAbi.Events[eventName]; ok {
		for _, arg := range event.Inputs {
			if arg.Indexed {
				//Cast to non-indexed
				args = append(args, abi.Argument{
					Name: arg.Name,
					Type: arg.Type,
				})
			}
		}
	}
	return args, nil
}

func GetMapEvents() map[string]*abi.Event {
	return mapEvents
}

func GetEventByName(name string) (*abi.Event, bool) {
	event, ok := mapEvents[name]
	return event, ok
}

func AbiUnpack(data []byte, types ...string) ([]interface{}, error) {
	var arguments ethabi.Arguments
	for _, t := range types {
		typ, err := ethabi.NewType(t, t, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create type: %w", err)
		}
		arguments = append(arguments, ethabi.Argument{Type: typ})
	}
	args, err := arguments.Unpack(data)
	if err != nil {
		return nil, fmt.Errorf("failed to get arguments: %w", err)
	}
	return args, nil
}

func AbiUnpackIntoMap(v map[string]interface{}, data []byte, types ...byte) error {
	var arguments ethabi.Arguments
	for _, t := range types {
		arguments = append(arguments, ethabi.Argument{Type: ethabi.Type{T: t}})
	}
	err := arguments.UnpackIntoMap(v, data)
	if err != nil {
		return fmt.Errorf("failed to get arguments: %w", err)
	}
	return nil
}
