package eip712

import (
	"github.com/ethereum/go-ethereum/common"
)

// TypedDataDomain represents the domain separator parameters for EIP-712
type TypedDataDomain struct {
	Name              string         `json:"name"`
	Version           string         `json:"version"`
	ChainId           uint64         `json:"chainId"`
	VerifyingContract common.Address `json:"verifyingContract"`
}
