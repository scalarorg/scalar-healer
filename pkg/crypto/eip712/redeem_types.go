package eip712

import (
	"math/big"

	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

// RedeemRequestMessage represents the message data for EIP-712 signing
type RedeemRequestMessage struct {
	BaseMessage
	Symbol string   `json:"symbol"`
	Amount *big.Int `json:"amount"`
	Nonce  *big.Int `json:"nonce"`
}

var _ EIP712Message = &RedeemRequestMessage{}

// RedeemRequestTypes defines the EIP-712 type structure for redeem request signing
var RedeemRequestTypes = apitypes.Types{
	"EIP712Domain": []apitypes.Type{
		{Name: "name", Type: "string"},
		{Name: "version", Type: "string"},
		{Name: "chainId", Type: "uint256"},
		{Name: "verifyingContract", Type: "address"},
	},
	"RedeemRequest": []apitypes.Type{
		{Name: "symbol", Type: "string"},
		{Name: "amount", Type: "uint256"},
		{Name: "nonce", Type: "uint64"},
	},
}

// NewRedeemRequestMessage creates a new RedeemRequestMessage instance
func NewRedeemRequestMessage(symbol string, amount *big.Int, nonce uint64) *RedeemRequestMessage {
	msg := &RedeemRequestMessage{
		Symbol: symbol,
		Amount: amount,
		Nonce:  big.NewInt(int64(nonce)),
	}
	msg.BaseMessage = *NewBaseMessage(
		RedeemRequestTypes,
		"RedeemRequest",
		map[string]interface{}{
			"symbol": msg.Symbol,
			"amount": msg.Amount,
			"nonce":  msg.Nonce,
		},
	)
	return msg
}
