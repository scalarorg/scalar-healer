package eip4361_test

import (
	"encoding/hex"
	"fmt"
	"log"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/scalarorg/scalar-healer/pkg/crypto/eip4361"
	"github.com/stretchr/testify/assert"
)

func TestNewSiweMessage(t *testing.T) {
	domain := "example.com"
	address := common.HexToAddress("0x71C7656EC7ab88b098defB751B7401B5f6d8976F")

	msg := eip4361.NewSiweMessage(domain, address)

	assert.Equal(t, domain, msg.GetDomain())
	assert.Equal(t, address.Hex(), msg.GetAddress().Hex())
	assert.Equal(t, "Sign in with Ethereum to the Scalar Healer application", msg.GetStatement())
	assert.Equal(t, "https://"+domain, msg.GetURI())
	assert.Equal(t, "1", msg.GetVersion())
	assert.Equal(t, int64(1), msg.GetChainID())
	assert.NotEmpty(t, msg.GetNonce())
}

func TestValidate(t *testing.T) {
	domain := "example.com"
	address := common.HexToAddress("0x71C7656EC7ab88b098defB751B7401B5f6d8976F")

	msg := eip4361.NewSiweMessage(domain, address)

	// Test with invalid signature
	invalidSig := "10000000000000000000000000000000000000000009"
	_, err := eip4361.Validate(msg.String(), invalidSig)
	assert.Error(t, err)
}

func TestSignAndValidate(t *testing.T) {
	domain := "example.com"

	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	msg := eip4361.NewSiweMessage(domain, address)

	data := []byte(msg.String())
	hash := crypto.Keccak256Hash([]byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)))

	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	msg, err = eip4361.Validate(msg.String(), "0x"+hex.EncodeToString(signature))
	assert.NoError(t, err)

	t.Log(msg)
	assert.Equal(t, domain, msg.GetDomain())
}
