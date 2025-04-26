package eip712_test

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/joho/godotenv"
	"github.com/scalarorg/scalar-healer/pkg/crypto/eip712"
)

var mockAddress = common.HexToAddress("0x24a1dB57Fa3ecAFcbaD91d6Ef068439acEeAe090")
var mockMessage = map[string]interface{}{
	"symbol": "ETH",
	"amount": big.NewInt(123456),
	"nonce":  big.NewInt(0),
}
var mockTypedData = eip712.CreateTypedData(apitypes.Types{
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
}, "RedeemRequest", &eip712.TypedDataDomain{
	Name:              "ScalarGateway",
	Version:           "1",
	ChainId:           1,
	VerifyingContract: mockAddress,
}, mockMessage)

func TestHashTypedData(t *testing.T) {
	// Generate hash for the typed data
	hash, err := eip712.HashTypedData(mockTypedData)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(common.Bytes2Hex(hash))

	if common.Bytes2Hex(hash) != "3cbf3ba715881b1fea2442af84a30e5074f835c2fd0039b5690fce2d20b9b8f1" {
		t.Fatal("hash is not correct")
	}
}

func TestSignTypedData1(t *testing.T) {
	err := godotenv.Load(os.ExpandEnv("../../../.env.test"))
	if err != nil {
		log.Fatalf("Error loading .env.test file: %s", err)
	}

	privKey := os.Getenv("PRIVATE_KEY")
	fmt.Println("private key: ", privKey)
	privateKey, err := crypto.HexToECDSA(privKey)
	if err != nil {
		t.Fatal(err)
	}

	// Generate signature for the typed data
	signature, err := eip712.SignTypedData(mockTypedData, privateKey)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(common.Bytes2Hex(signature))

	// verify signature
	hash, err := eip712.HashTypedData(mockTypedData)
	if err != nil {
		t.Fatal(err)
	}

	// recover public key
	publicKey, err := crypto.SigToPub(hash, signature)
	if err != nil {
		t.Fatal(err)
	}
	// get address
	address := crypto.PubkeyToAddress(*publicKey)
	t.Log(address.Hex())

	if address.Hex() != mockAddress.Hex() {
		t.Fatal("address is not correct")
	}
}
