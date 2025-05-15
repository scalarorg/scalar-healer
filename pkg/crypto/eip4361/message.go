package eip4361

import (
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spruceid/siwe-go"
)

func NewSiweMessage(domain string, address common.Address) *siwe.Message {
	options := make(map[string]interface{})
	options["statement"] = "Sign in with Ethereum to the Scalar Healer application"
	options["chainId"] = 1 // Default to Ethereum mainnet
	options["issuedAt"] = time.Now()
	options["expirationTime"] = time.Now().Add(time.Hour)

	msg, err := siwe.InitMessage(domain, address.Hex(), fmt.Sprintf("https://%s", domain), siwe.GenerateNonce(), options)
	if err != nil {
		panic(err)
	}
	return msg

}

func Validate(message, signature string) (*siwe.Message, error) {
	msg, err := siwe.ParseMessage(message)
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(signature, "0x") {
		signature = "0x" + signature
	}

	if len(signature) != 132 {
		return nil, fmt.Errorf("invalid signature length")
	}

	_, err = msg.Verify(signature, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	return msg, nil
}
