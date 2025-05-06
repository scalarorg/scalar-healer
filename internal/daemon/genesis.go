package daemon

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/scalarorg/scalar-healer/config"
	"github.com/scalarorg/scalar-healer/pkg/db"
	"github.com/scalarorg/scalar-healer/pkg/db/sqlc"
	"golang.org/x/crypto/sha3"
)

func (s *Service) initGenesis(ctx context.Context) {
	protocolCfgPath := fmt.Sprintf("%s/protocols.json", s.ConfigPath)
	protocols, err := config.ReadJsonArrayConfig[db.Protocol](protocolCfgPath)
	if err != nil {
		panic(err)
	}
	for _, protocol := range protocols {
		protocol.CustodianGroupUid = sha3.Sum256([]byte(protocol.CustodianGroupName))
	}
	s.DbAdapter.SaveProtocols(ctx, protocols)
	s.initTokens(ctx, protocols)
}

func (s *Service) initTokens(ctx context.Context, protocols []db.Protocol) {
	tokens := make([]sqlc.Token, 0)
	for _, protocol := range protocols {
		for _, evmClient := range s.EvmClients {
			tokenAddr := evmClient.GetTokenAddressBySymbol(protocol.Symbol)
			if tokenAddr == nil {
				log.Error().Msgf("Token address not found for symbol %s", protocol.Symbol)
				continue
			}
			token := sqlc.Token{
				Protocol: protocol.Name,
				Address:  tokenAddr.Bytes(),
				ChainID:  db.ConvertUint64ToNumeric(evmClient.EvmConfig.ChainID),
				Active:   true,
				Decimal:  db.ConvertUint64ToNumeric(uint64(protocol.Decimals)),
				Symbol:   protocol.Symbol,
				Name:     protocol.Asset,
				Avatar:   protocol.Avatar,
			}
			tokens = append(tokens, token)
		}
	}
	s.DbAdapter.SaveTokens(ctx, tokens)
}
