package daemon

import (
	"bytes"
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/scalarorg/scalar-healer/config"
	"github.com/scalarorg/scalar-healer/pkg/db"
	"github.com/scalarorg/scalar-healer/pkg/db/sqlc"
	"golang.org/x/crypto/sha3"
)

func (s *Service) initGenesis(ctx context.Context) {
	s.initCustodianGroups(ctx)
	s.initProtocols(ctx)
}

func (s *Service) initCustodianGroups(ctx context.Context) {
	custodianGroupCfgPath := fmt.Sprintf("%s/custodian_groups.json", s.ConfigPath)
	custodianGroups, err := config.ReadJsonArrayConfig[sqlc.CustodianGroup](custodianGroupCfgPath)
	if err != nil {
		panic(err)
	}

	for _, custodianGroup := range custodianGroups {
		uid := sha3.Sum256([]byte(custodianGroup.Name))
		log.Info().Msgf("Custodian Group's UID: %x", custodianGroup.Uid)
		log.Info().Msgf("Custodian Group's BitcoinPubKey: %x", custodianGroup.BitcoinPubkey)
		if !bytes.Equal(uid[:], custodianGroup.Uid) {
			panic("custodian group uid is not correct")
		}
	}

	s.DbAdapter.SaveCustodianGroups(ctx, custodianGroups)
}

func (s *Service) initProtocols(ctx context.Context) {
	protocolCfgPath := fmt.Sprintf("%s/protocols.json", s.ConfigPath)
	protocols, err := config.ReadJsonArrayConfig[sqlc.Protocol](protocolCfgPath)
	if err != nil {
		panic(err)
	}

	newProtocols := make([]sqlc.Protocol, len(protocols))
	for ind, protocol := range protocols {
		uid := sha3.Sum256([]byte(protocol.CustodianGroupName))
		protocol.CustodianGroupUid = uid[:]
		protocols[ind] = protocol
	}
	s.DbAdapter.SaveProtocols(ctx, newProtocols)

	s.initTokens(ctx, protocols)
}

func (s *Service) initTokens(ctx context.Context, protocols []sqlc.Protocol) {
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
