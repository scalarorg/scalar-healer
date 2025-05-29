package daemon

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/scalarorg/scalar-healer/config"
	"github.com/scalarorg/scalar-healer/pkg/db"
	"github.com/scalarorg/scalar-healer/pkg/db/sqlc"
	"github.com/scalarorg/scalar-healer/pkg/evm"
	"github.com/scalarorg/scalar-healer/pkg/utils/funcs"
	"golang.org/x/crypto/sha3"
)

func (s *Service) initGenesis(ctx context.Context) {
	s.initCustodianGroups(ctx)
	s.initProtocols(ctx)
	s.initGatewayAddresses(ctx)
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

	err = s.CombinedAdapter.SaveCustodianGroups(ctx, custodianGroups)
	if err != nil {
		panic(err)
	}
}

func (s *Service) initProtocols(ctx context.Context) {
	protocolCfgPath := fmt.Sprintf("%s/protocols.json", s.ConfigPath)
	protocols, err := config.ReadJsonArrayConfig[sqlc.Protocol](protocolCfgPath)
	if err != nil {
		panic(err)
	}

	for ind, protocol := range protocols {
		uid := sha3.Sum256([]byte(protocol.CustodianGroupName))
		protocol.CustodianGroupUid = uid[:]
		protocols[ind] = protocol
	}

	log.Info().Msgf("protocols: %v", protocols)
	err = s.CombinedAdapter.SaveProtocols(ctx, protocols)
	if err != nil {
		panic(err)
	}

	s.initTokens(ctx, protocols)
}

func (s *Service) initGatewayAddresses(ctx context.Context) {
	evmCfgPath := fmt.Sprintf("%s/evm.json", s.ConfigPath)
	configs, err := config.ReadJsonArrayConfig[evm.EvmNetworkConfig](evmCfgPath)
	if err != nil {
		panic(err)
	}

	chainIds := make([]uint64, 0)
	gatewayAddresses := make([][]byte, 0)

	for _, config := range configs {
		chainIds = append(chainIds, config.ChainID)
		add := funcs.Must(hex.DecodeString(strings.TrimPrefix(config.Gateway, "0x")))
		gatewayAddresses = append(gatewayAddresses, add)
	}

	err = s.CombinedAdapter.CreateGatewayAddresses(ctx, gatewayAddresses, chainIds)
	if err != nil {
		panic(err)
	}
}

func (s *Service) initTokens(ctx context.Context, protocols []sqlc.Protocol) {
	tokens := make([]sqlc.Token, 0)
	for _, protocol := range protocols {
		for _, evmClient := range s.EvmClients {
			tokenAddr := evmClient.GetTokenAddressBySymbol(protocol.Symbol)
			if tokenAddr == nil {
				panic(fmt.Sprintf("Token address not found for symbol %s", protocol.Symbol))
			}
			token := sqlc.Token{
				Symbol:  protocol.Symbol,
				Address: tokenAddr.Bytes(),
				ChainID: db.ConvertUint64ToNumeric(evmClient.EvmConfig.ChainID),
				Active:  true,
			}
			tokens = append(tokens, token)
		}
	}
	s.CombinedAdapter.SaveTokens(ctx, tokens)
}
