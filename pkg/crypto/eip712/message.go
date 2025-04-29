package eip712

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/scalarorg/scalar-healer/pkg/db"
)

// EIP712Message defines the interface for EIP-712 message types
type EIP712Message interface {
	// ToTypedData converts the message to EIP-712 typed data
	ToTypedData(contractAddress common.Address) apitypes.TypedData
	Validate(ctx context.Context, db db.DbAdapter, contractAddress *common.Address) error
}

var _ EIP712Message = &BaseMessage{}

// BaseMessage provides common functionality for EIP-712 messages

type BaseRequest struct {
	Address   string `json:"address" validate:"eth_addr"`
	Signature string `json:"signature" validate:"hexadecimal"`
	Nonce     uint64 `json:"nonce" validate:"gte=0"`
	ChainID   uint64 `json:"chain_id" validate:"required,gte=0"`
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
func (m *BaseMessage) ToTypedData(gatewayAddress common.Address) apitypes.TypedData {
	return CreateTypedData(
		m.Types,
		m.PrimaryType,
		&TypedDataDomain{
			Name:              "ScalarGateway",
			Version:           "1",
			ChainId:           m.ChainID,
			VerifyingContract: gatewayAddress,
		},
		m.Message,
	)
}

func (m *BaseMessage) Validate(ctx context.Context, db db.DbAdapter, contractAddress *common.Address) error {
	address := common.HexToAddress(m.Address)
	nonce := db.GetNonce(ctx, address)
	if nonce != m.Nonce {
		return fmt.Errorf("invalid nonce")
	}

	// Create and convert message to EIP-712 typed data
	typedData := m.ToTypedData(*contractAddress)

	// Verify the signature
	err := VerifySignTypedData(typedData, address, common.FromHex(m.Signature))
	if err != nil {
		return err
	}
	return nil
}
