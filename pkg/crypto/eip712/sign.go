package eip712

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/scalarorg/scalar-healer/constants"
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

func VerifySignTypedData(typedData apitypes.TypedData, address common.Address, signature []byte) error {
	hash, err := HashTypedData(typedData)
	if err != nil {
		return fmt.Errorf("hash typed data: %w", err)
	}

	// recover public key
	publicKey, err := crypto.SigToPub(hash, signature)
	if err != nil {
		return fmt.Errorf("sig to pub: %w", err)
	}

	_addrress := crypto.PubkeyToAddress(*publicKey)

	if _addrress.Hex() != address.Hex() {
		return constants.ErrInvalidSignature
	}
	return nil
}

// HashTypedData generates the hash of EIP-712 typed data according to the specification
func HashTypedData(typedData apitypes.TypedData) ([]byte, error) {
	// Hash the domain separator
	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return nil, fmt.Errorf("failed to hash domain separator: %w", err)
	}

	// Hash the message
	messageHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return nil, fmt.Errorf("failed to hash message: %w", err)
	}

	// Encode the final hash according to EIP-712
	return crypto.Keccak256(
		[]byte{0x19, 0x01}, // EIP-712 prefix and version
		domainSeparator,    // Domain separator hash
		messageHash,        // Message hash
	), nil
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
