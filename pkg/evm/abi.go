package evm

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/rs/zerolog/log"
	contracts "github.com/scalarorg/scalar-healer/pkg/evm/contracts/generated"
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

func getScalarGatewayAbi() (*abi.ABI, error) {
	if scalarGatewayAbi == nil {
		var err error
		scalarGatewayAbi, err = contracts.IScalarGatewayMetaData.GetAbi()
		if err != nil {
			return nil, err
		}
	}
	return scalarGatewayAbi, nil
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
