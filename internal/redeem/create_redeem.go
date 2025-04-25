package redeem

import (
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/labstack/echo/v4"
	"github.com/scalarorg/scalar-healer/pkg/db/mongo"
	"github.com/scalarorg/scalar-healer/pkg/eip712"
	"github.com/scalarorg/scalar-healer/pkg/utils"
)

type CreateRedeemRequest struct {
	Address   string `json:"address" validate:"eth_addr"`
	Signature string `json:"signature" validate:"hexadecimal"` // not need to validte length
	ChainID   uint64 `json:"chain_id" validate:"required,gte=0"`
	Symbol    string `json:"symbol" validate:"required"`
	Amount    string `json:"amount" validate:"required"` // bigint format
	Nonce     uint64 `json:"nonce" validate:"required,gte=0"`
}

// RedeemRequestTypes defines the EIP-712 type structure for order signing
var RedeemRequestTypes = apitypes.Types{
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
}

const primaryType = "RedeemRequest"

func CreateRedeem(c echo.Context) error {
	var body CreateRedeemRequest
	if err := utils.BindAndValidate(c, &body); err != nil {
		return err
	}

	db := mongo.GetRepositoryFromContext(c)

	// TODO: check nonce match with nonce in db
	// TODO: parse amount to bigint
	db.GetRedeemNonce(common.Address{})

	// Create message data for EIP-712 signing
	message := map[string]interface{}{
		"symbol": body.Symbol,
		"amount": body.Amount,
		"nonce":  body.Nonce,
	}

	// Create EIP-712 typed data
	typedData := eip712.CreateTypedData(RedeemRequestTypes, primaryType, getDomain(common.Address{}), message)

	// Generate hash for the typed data
	hash, err := eip712.HashTypedData(typedData)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, common.Bytes2Hex(hash))
}

func getDomain(gatewayAddress common.Address) *eip712.TypedDataDomain {
	return &eip712.TypedDataDomain{
		Name:              "ScalarGateway",
		Version:           "1",
		ChainId:           1,
		VerifyingContract: gatewayAddress,
	}
}
