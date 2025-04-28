package daemon

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/scalarorg/scalar-healer/config"
	"github.com/scalarorg/scalar-healer/pkg/db/models"
	"golang.org/x/crypto/sha3"
)

func (s *Service) initGenesis(ctx context.Context) {
	protocolCfgPath := fmt.Sprintf("%s/protocols.json", s.ConfigPath)
	protocols, err := config.ReadJsonArrayConfig[models.Protocol](protocolCfgPath)
	if err != nil {
		panic(err)
	}
	for _, protocol := range protocols {
		protocol.CustodianGroupUid = sha3.Sum256([]byte(protocol.CustodianGroupName))
	}
	s.DbAdapter.SaveProtocols(ctx, protocols)
	s.initTokens(ctx, protocols)
}

func (s *Service) initTokens(ctx context.Context, protocols []models.Protocol) {
	tokens := make([]models.Token, 0)
	for _, protocol := range protocols {
		for _, evmClient := range s.EvmClients {
			tokenAddr := evmClient.GetTokenAddressBySymbol(protocol.Symbol)
			if tokenAddr == nil {
				log.Error().Msgf("Token address not found for symbol %s", protocol.Symbol)
				continue
			}
			token := models.Token{
				Protocol:  protocol.Name,
				Address:   tokenAddr.Bytes(),
				ChainID:   evmClient.EvmConfig.ChainID,
				Active:    true,
				Decimal:   uint64(protocol.Decimals),
				Symbol:    protocol.Symbol,
				Name:      protocol.Asset,
				Avatar:    protocol.Avatar,
				CreatedAt: uint64(time.Now().Second()),
				UpdatedAt: uint64(time.Now().Second()),
			}
			tokens = append(tokens, token)
		}
	}
	s.DbAdapter.SaveTokenInfos(ctx, tokens)
}
