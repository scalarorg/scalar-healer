package eip712

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

// CreateBridgeMessage represents the message data for EIP-712 signing
type CreateBridgeMessage struct {
	BaseMessage
	ChainID *big.Int    `json:"chain_id"`
	TxHash  common.Hash `json:"tx_hash"`
}

// CreateBridgeTypes defines the EIP-712 type structure for bridge request signing
var CreateBridgeTypes = apitypes.Types{
	"EIP712Domain": []apitypes.Type{
		{Name: "name", Type: "string"},
		{Name: "version", Type: "string"},
		{Name: "chainId", Type: "uint256"},
		{Name: "verifyingContract", Type: "address"},
	},
	"CreateBridgeRequest": []apitypes.Type{
		{Name: "chain_id", Type: "uint64"},
		{Name: "tx_hash", Type: "bytes32"},
	},
}

// NewCreateBridgeMessage creates a new CreateBridgeMessage instance
func NewCreateBridgeMessage(chainID *big.Int, txHash common.Hash) *CreateBridgeMessage {
	msg := &CreateBridgeMessage{
		ChainID: chainID,
		TxHash:  txHash,
	}
	msg.BaseMessage = *NewBaseMessage(
		CreateBridgeTypes,
		"CreateBridgeRequest",
		map[string]interface{}{
			"chain_id": msg.ChainID,
			"tx_hash":  msg.TxHash,
		},
	)
	return msg
}
