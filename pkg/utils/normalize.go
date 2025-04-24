package utils

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/scalarorg/bitcoin-vault/go-utils/btc"
	"github.com/scalarorg/bitcoin-vault/go-utils/chain"
	"github.com/scalarorg/bitcoin-vault/go-utils/types"
)

func NormalizeAddress(address string, chainType types.ChainType) string {
	switch chainType {
	case types.ChainTypeEVM:
		address = strings.ToLower(address)
		if !strings.HasPrefix(address, "0x") {
			address = "0x" + address
		}
		return address
	default:
		return address
	}
}

func NormalizeHash(hash string) string {
	return strings.ToLower(strings.TrimPrefix(hash, "0x"))
}

func CalculateDestinationAddress(payload []byte, chainInfoBytes *chain.ChainInfoBytes) (destinationAddress string, err error) {
	if chainInfoBytes.ChainType() != types.ChainTypeBitcoin {
		return "", nil
	}

	decodedPayload, err := DecodeContractCallWithTokenPayload(payload)
	if err != nil || decodedPayload == nil {
		return "", fmt.Errorf("invalid payload: %v", decodedPayload)
	}

	params := btc.BtcChainsRecords().GetChainParamsByID(chainInfoBytes.ChainID())

	if params == nil {
		return "", fmt.Errorf("invalid destination chain: %d", chainInfoBytes.ChainID())
	}

	if decodedPayload.CustodianOnly != nil {
		identifier := decodedPayload.CustodianOnly.RecipientChainIdentifier
		address, err := btc.ScriptPubKeyToAddress(identifier, params.Name)
		if err != nil {
			return "", fmt.Errorf("failed to convert script pubkey %s to address with params name %s: %w",
				hex.EncodeToString(identifier), params.Name, err)
		}
		return address.EncodeAddress(), nil
	} else if decodedPayload.UPC != nil && decodedPayload.UPC.Psbt != nil {
		packet, err := psbt.NewFromRawBytes(
			bytes.NewReader(decodedPayload.UPC.Psbt), false,
		)

		if err != nil {
			return "", fmt.Errorf("failed to create psbt packet: %w", err)
		}

		identifier := packet.UnsignedTx.TxOut[1].PkScript
		address, err := btc.ScriptPubKeyToAddress(identifier, params.Name)
		if err != nil {
			return "", fmt.Errorf("failed to convert script pubkey to address: %w", err)
		}
		return address.EncodeAddress(), nil
	}

	return "", fmt.Errorf("invalid payload: %v", decodedPayload)
}
