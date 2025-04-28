package eip712

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

// EIP712Message defines the interface for EIP-712 message types
type EIP712Message interface {
	// ToTypedData converts the message to EIP-712 typed data
	ToTypedData(gatewayAddress common.Address, chainId uint64) apitypes.TypedData
}

var _ EIP712Message = &BaseMessage{}

// BaseMessage provides common functionality for EIP-712 messages
type BaseMessage struct {
	Types       apitypes.Types
	PrimaryType string
	Message     map[string]interface{}
}

// NewBaseMessage creates a new BaseMessage instance
func NewBaseMessage(types apitypes.Types, primaryType string, message map[string]interface{}) *BaseMessage {
	return &BaseMessage{
		Types:       types,
		PrimaryType: primaryType,
		Message:     message,
	}
}

// ToTypedData converts the BaseMessage to EIP-712 typed data
func (m *BaseMessage) ToTypedData(gatewayAddress common.Address, chainId uint64) apitypes.TypedData {
	return CreateTypedData(
		m.Types,
		m.PrimaryType,
		&TypedDataDomain{
			Name:              "ScalarGateway",
			Version:           "1",
			ChainId:           chainId,
			VerifyingContract: gatewayAddress,
		},
		m.Message,
	)
}
