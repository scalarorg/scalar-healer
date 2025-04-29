package eip712

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

// TransferRequestMessage represents the message data for EIP-712 signing
type TransferRequestMessage struct {
	*BaseMessage
	*BaseRequest
	DestinationChain   string          `json:"destinationChain"`
	DestinationAddress *common.Address `json:"destinationAddress"`
	Symbol             string          `json:"symbol"`
	Amount             *big.Int        `json:"amount"`
}

const TransferRequestDomainName = "TransferRequest"

// TransferRequestTypes defines the EIP-712 type structure for Transfer request signing
var TransferRequestTypes = apitypes.Types{
	"EIP712Domain": []apitypes.Type{
		{Name: "name", Type: "string"},
		{Name: "version", Type: "string"},
		{Name: "chainId", Type: "uint256"},
		{Name: "verifyingContract", Type: "address"},
	},
	TransferRequestDomainName: []apitypes.Type{
		{Name: "destinationChain", Type: "string"},
		{Name: "destinationAddress", Type: "address"},
		{Name: "symbol", Type: "string"},
		{Name: "amount", Type: "uint256"},
		{Name: "nonce", Type: "uint64"},
	},
}

// NewTransferRequestMessage creates a new TransferRequestMessage instance
func NewTransferRequestMessage(baseRequest *BaseRequest, destChain string, destAddress *common.Address, symbol string, amount *big.Int) *TransferRequestMessage {
	msg := &TransferRequestMessage{
		DestinationChain:   destChain,
		DestinationAddress: destAddress,
		Symbol:             symbol,
		Amount:             amount,
		BaseRequest:        baseRequest,
	}
	msg.BaseMessage = NewBaseMessage(
		TransferRequestTypes,
		TransferRequestDomainName,
		map[string]interface{}{
			"destinationChain":   destChain,
			"destinationAddress": destAddress.Bytes(),
			"symbol":             symbol,
			"amount":             amount,
			"nonce":              big.NewInt(int64(baseRequest.Nonce)),
		},
		baseRequest,
	)
	return msg
}
