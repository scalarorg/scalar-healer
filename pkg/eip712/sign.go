package eip712

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
)

// SignTypedData signs EIP-712 typed data with the provided private key
func SignTypedData(typedData apitypes.TypedData, privateKey *ecdsa.PrivateKey) ([]byte, error) {
	hash, err := HashTypedData(typedData)
	if err != nil {
		return nil, fmt.Errorf("hash typed data: %w", err)
	}

	signature, err := crypto.Sign(hash, privateKey)
	if err != nil {
		return nil, fmt.Errorf("sign hash: %w", err)
	}

	return signature, nil
}

// HashTypedData generates the hash of EIP-712 typed data
func HashTypedData(typedData apitypes.TypedData) ([]byte, error) {
	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return nil, fmt.Errorf("hash domain: %w", err)
	}

	typed, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return nil, fmt.Errorf("hash message: %w", err)
	}

	return crypto.Keccak256([]byte{
		0x19, // prefix
		0x01, // version
	}, domainSeparator, typed), nil
}

// CreateTypedData creates an EIP-712 typed data structure for signing
func CreateTypedData(types apitypes.Types, primaryType string, domain *TypedDataDomain, message map[string]interface{}) apitypes.TypedData {
	return apitypes.TypedData{
		Types:       types,
		PrimaryType: primaryType,
		Domain: apitypes.TypedDataDomain{
			Name:              domain.Name,
			Version:           domain.Version,
			ChainId:           math.NewHexOrDecimal256(int64(domain.ChainId)),
			VerifyingContract: domain.VerifyingContract.Hex(),
		},
		Message: message,
	}
}
