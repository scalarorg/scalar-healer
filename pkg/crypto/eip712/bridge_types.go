package eip712

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

// BridgeRequestMessage represents the message data for EIP-712 signing
type BridgeRequestMessage struct {
	*BaseMessage
	TxHash common.Hash `json:"tx_hash"`
}

const BridgeRequestDomainName = "BridgeRequest"

// BridgeRequestTypes defines the EIP-712 type structure for bridge request signing
var BridgeRequestTypes = apitypes.Types{
	"EIP712Domain": []apitypes.Type{
		{Name: "name", Type: "string"},
		{Name: "version", Type: "string"},
		{Name: "chainId", Type: "uint256"},
		{Name: "verifyingContract", Type: "address"},
	},
	BridgeRequestDomainName: []apitypes.Type{
		{Name: "tx_hash", Type: "bytes32"},
		{Name: "nonce", Type: "uint64"},
	},
}

// NewBridgeRequestMessage creates a new BridgeRequestMessage instance
func NewBridgeRequestMessage(baseRequest *BaseRequest, txHash common.Hash) *BridgeRequestMessage {
	msg := &BridgeRequestMessage{
		TxHash: txHash,
	}
	msg.BaseMessage = NewBaseMessage(
		BridgeRequestTypes,
		BridgeRequestDomainName,
		map[string]interface{}{
			"tx_hash": msg.TxHash,
			"nonce":   big.NewInt(int64(baseRequest.Nonce)),
		},
		baseRequest,
	)
	return msg
}
