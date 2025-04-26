package redeem

import (
	"fmt"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/labstack/echo/v4"
	"github.com/scalarorg/scalar-healer/pkg/crypto/eip712"
	"github.com/scalarorg/scalar-healer/pkg/db/mongo"
	"github.com/scalarorg/scalar-healer/pkg/utils"
)

type CreateRedeemRequest struct {
	Address   string `json:"address" validate:"eth_addr"`
	Signature string `json:"signature" validate:"hexadecimal"` // not need to validte length
	ChainID   uint64 `json:"chain_id" validate:"required,gte=0"`
	Symbol    string `json:"symbol" validate:"required"`
	Amount    string `json:"amount" validate:"required"` // bigint format
	Nonce     uint64 `json:"nonce" validate:"gte=0"`
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

	ctx := c.Request().Context()

	db := mongo.GetRepositoryFromContext(c)

	gatewayAddress := db.GetGatewayAddress(ctx, body.ChainID)
	if gatewayAddress == nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("not found gateway address for chain id: %d", body.ChainID))
	}

	address := common.HexToAddress(body.Address)

	nonce := db.GetRedeemNonce(ctx, address)
	if nonce != body.Nonce {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid nonce")
	}

	if !db.CheckTokenExists(ctx, body.Symbol) {
		return echo.NewHTTPError(http.StatusBadRequest, "token not exists")
	}

	amountz, ok := utils.StringToBigInt(body.Amount)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid amount")
	}

	// Create message data for EIP-712 signing
	message := map[string]interface{}{
		"symbol": body.Symbol,
		"amount": amountz,
		"nonce":  big.NewInt(int64(nonce)),
	}

	// Create EIP-712 typed data
	typedData := eip712.CreateTypedData(RedeemRequestTypes, primaryType, getDomain(*gatewayAddress), message)

	// Verify the signature
	err := eip712.VerifySignTypedData(typedData, address, common.FromHex(body.Signature))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("invalid signature: %s", err.Error()))
	}

	// Save redeem request
	err = db.SaveRedeemRequest(ctx, body.ChainID, address, common.FromHex(body.Signature), amountz, body.Symbol, nonce)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to save redeem request")
	}

	return c.NoContent(http.StatusOK)
}

func getDomain(gatewayAddress common.Address) *eip712.TypedDataDomain {
	return &eip712.TypedDataDomain{
		Name:              "ScalarGateway",
		Version:           "1",
		ChainId:           1,
		VerifyingContract: gatewayAddress,
	}
}
