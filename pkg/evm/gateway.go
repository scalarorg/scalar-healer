package evm

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
)

func (c *EvmClient) GetTokenAddressBySymbol(symbol string) *common.Address {
	callOpts, err := c.CreateCallOpts()
	if err != nil {
		log.Error().Err(err).Msg("[EvmClient] CreateCallOpts")
		return nil
	}
	tokenAddr, err := c.Gateway.TokenAddresses(callOpts, symbol)
	if err != nil {
		log.Error().Err(err).Msg("[EvmClient] TokenAddresses")
		return nil
	}
	return &tokenAddr
}
