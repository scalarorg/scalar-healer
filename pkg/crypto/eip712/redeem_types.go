package eip712

import (
	"math/big"

	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

// RedeemRequestMessage represents the message data for EIP-712 signing
type RedeemRequestMessage struct {
	*BaseMessage
	*BaseRequest
	Symbol string   `json:"symbol"`
	Amount *big.Int `json:"amount"`
}

const RedeemRequestDomainName = "RedeemRequest"

// RedeemRequestTypes defines the EIP-712 type structure for redeem request signing
var RedeemRequestTypes = apitypes.Types{
	"EIP712Domain": []apitypes.Type{
		{Name: "name", Type: "string"},
		{Name: "version", Type: "string"},
		{Name: "chainId", Type: "uint256"},
		{Name: "verifyingContract", Type: "address"},
	},
	RedeemRequestDomainName: []apitypes.Type{
		{Name: "symbol", Type: "string"},
		{Name: "amount", Type: "uint256"},
		{Name: "nonce", Type: "uint64"},
	},
}

// NewRedeemRequestMessage creates a new RedeemRequestMessage instance
func NewRedeemRequestMessage(baseRequest *BaseRequest, symbol string, amount *big.Int) *RedeemRequestMessage {
	msg := &RedeemRequestMessage{
		Symbol:      symbol,
		Amount:      amount,
		BaseRequest: baseRequest,
	}
	msg.BaseMessage = NewBaseMessage(
		RedeemRequestTypes,
		RedeemRequestDomainName,
		map[string]interface{}{
			"symbol": msg.Symbol,
			"amount": msg.Amount,
			"nonce":  big.NewInt(int64(baseRequest.Nonce)),
		},
		baseRequest,
	)
	return msg
}
