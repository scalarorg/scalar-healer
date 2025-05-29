package eip712

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/scalarorg/scalar-healer/pkg/utils/chains"
)

// EIP712Message defines the interface for EIP-712 message types
type EIP712Message interface {
	// ToTypedData converts the message to EIP-712 typed data
	ToTypedData(contractAddress common.Address) (*apitypes.TypedData, error)
	Validate(ctx context.Context, contractAddress *common.Address) error
}

var _ EIP712Message = &BaseMessage{}

// BaseMessage provides common functionality for EIP-712 messages

type BaseRequest struct {
	Address   common.Address `json:"address"`
	Signature string         `json:"signature" validate:"hexadecimal"`
	Nonce     uint64         `json:"nonce" validate:"gte=0"`
	Chain     string         `json:"chain" validate:"required"`
}

type BaseMessage struct {
	*BaseRequest
	Types       apitypes.Types
	PrimaryType string
	Message     map[string]interface{}
}

// NewBaseMessage creates a new BaseMessage instance
func NewBaseMessage(types apitypes.Types, primaryType string, message map[string]interface{}, req *BaseRequest) *BaseMessage {
	return &BaseMessage{
		Types:       types,
		PrimaryType: primaryType,
		Message:     message,
		BaseRequest: req,
	}
}

// ToTypedData converts the BaseMessage to EIP-712 typed data
func (m *BaseMessage) ToTypedData(gatewayAddress common.Address) (*apitypes.TypedData, error) {
	if !chains.ChainName(m.Chain).IsEvmChain() {
		return nil, fmt.Errorf("invalid chain")
	}

	chainID, err := chains.ChainName(m.Chain).GetChainID()
	if err != nil {
		return nil, err
	}

	data := CreateTypedData(
		m.Types,
		m.PrimaryType,
		&TypedDataDomain{
			Name:              "ScalarGateway",
			Version:           "1",
			ChainId:           chainID.Uint64(),
			VerifyingContract: gatewayAddress,
		},
		m.Message,
	)

	return data, nil
}

func (m *BaseMessage) Validate(ctx context.Context, contractAddress *common.Address) error {
	if contractAddress == nil {
		return fmt.Errorf("contract address is nil")
	}
	// Create and convert message to EIP-712 typed data
	typedData, err := m.ToTypedData(*contractAddress)
	if err != nil {
		return err
	}

	// Verify the signature
	err = VerifySignTypedData(typedData, m.Address, common.FromHex(m.Signature))
	if err != nil {
		return err
	}
	return nil
}
