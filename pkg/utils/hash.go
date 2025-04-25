package utils

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// CreateRedeemMsgHash generates a keccak256 hash of the redeem message parameters
// following EVM's encoding standards. The parameters are packed in the following order:
// amount (uint256), symbol (string), nonce (uint64)
func CreateRedeemMsgHash(amount string, symbol string, nonce uint64) ([]byte, error) {
	// Convert amount string to big.Int
	amountBig := new(big.Int)
	_, success := amountBig.SetString(amount, 10)
	if !success {
		return nil, fmt.Errorf("invalid amount format: %s", amount)
	}

	// Pack the parameters according to Solidity's tight packing rules
	// 1. amount (uint256) - pad to 32 bytes
	amountBytes := common.LeftPadBytes(amountBig.Bytes(), 32)

	// 2. symbol (string) - encode as bytes
	symbolBytes := []byte(symbol)

	// 3. nonce (uint64) - convert to bytes
	nonceBytes := make([]byte, 8)
	big.NewInt(int64(nonce)).FillBytes(nonceBytes)

	// Concatenate all parameters
	message := append(amountBytes, symbolBytes...)
	message = append(message, nonceBytes...)

	// Calculate keccak256 hash
	hash := crypto.Keccak256(message)
	return hash, nil
}
