package eip712_test

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/joho/godotenv"
	"github.com/scalarorg/scalar-healer/pkg/crypto/eip712"
)

var mockAddress = common.HexToAddress("0x24a1dB57Fa3ecAFcbaD91d6Ef068439acEeAe090")

type suite struct {
	data eip712.EIP712Message
	want string
}

func TestHashTypedData(t *testing.T) {
	tests := []suite{
		{
			data: eip712.NewRedeemRequestMessage(&eip712.BaseRequest{
				Nonce: uint64(0),
			}, "ETH", big.NewInt(123456)),
			want: "f5043f1952bbc2803a9bc1ff8cff68dbfbcc3f229d2d8f780e21c6890b390dd4",
		},
	}

	for _, tt := range tests {
		hash, err := eip712.HashTypedData(tt.data.ToTypedData(mockAddress, 1))
		if err != nil {
			t.Fatal(err)
		}
		if common.Bytes2Hex(hash) != tt.want {
			t.Fatal("hash is not correct")
		}
	}

}

func TestSignTypedData(t *testing.T) {
	err := godotenv.Load(os.ExpandEnv("../../../.env.test"))
	if err != nil {
		log.Fatalf("Error loading .env.test file: %s", err)
	}

	privKey := os.Getenv("PRIVATE_KEY")
	privateKey, err := crypto.HexToECDSA(privKey)
	if err != nil {
		t.Fatal(err)
	}

	suites := []suite{
		{
			data: eip712.NewRedeemRequestMessage(&eip712.BaseRequest{
				Nonce: 0,
			}, "ETH", big.NewInt(123456)),
			want: mockAddress.Hex(),
		},
		{
			data: eip712.NewBridgeRequestMessage(&eip712.BaseRequest{
				Nonce: 0,
			}, big.NewInt(1), common.MaxHash),
			want: mockAddress.Hex(),
		},
	}

	for _, tt := range suites {
		typedData := tt.data.ToTypedData(mockAddress, 1)
		signature, err := eip712.SignTypedData(typedData, privateKey)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println("Signature: ", common.Bytes2Hex(signature))
		// verify signature
		hash, err := eip712.HashTypedData(typedData)
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
		if address.Hex() != tt.want {
			t.Fatal("address is not correct")
		}
	}

}
