package redeem

import (
	"fmt"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
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

	gatewayAddress, err := db.GetGatewayAddress(ctx, body.ChainID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("not found gateway address for chain: %d", body.ChainID))

	}
	log.Info().Interface("gatewayAddress", gatewayAddress).Msg("gatewayAddress")

	address := common.HexToAddress(body.Address)

	nonce := db.GetNonce(ctx, address)
	if nonce != body.Nonce {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid nonce")
	}

	_, err = db.GetTokenAddressBySymbol(ctx, body.ChainID, body.Symbol)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "token not exists")
	}

	// TODO: validate the balance on evm network
	amountz, ok := utils.StringToBigInt(body.Amount)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid amount")
	}

	// Create and convert message to EIP-712 typed data
	message := eip712.NewRedeemRequestMessage(body.Symbol, amountz, nonce)
	typedData := message.ToTypedData(*gatewayAddress, body.ChainID)

	// Verify the signature
	err = eip712.VerifySignTypedData(typedData, address, common.FromHex(body.Signature))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("failed to verify signature: %s", err.Error()))
	}

	// Save redeem request
	err = db.SaveRedeemRequest(ctx, body.ChainID, address, common.FromHex(body.Signature), amountz, body.Symbol, nonce)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to save redeem request")
	}

	return c.NoContent(http.StatusOK)
}
