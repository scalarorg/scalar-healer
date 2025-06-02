package daemon

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
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

type CustodianGroupConfig struct {
	ID            int64            `json:"id"`
	Uid           []byte           `json:"uid"`
	Name          string           `json:"name"`
	BitcoinPubkey []byte           `json:"bitcoin_pubkey"`
	Quorum        int64            `json:"quorum"`
	Custodians    []sqlc.Custodian `json:"custodians"`
}

func (s *Service) initCustodianGroups(ctx context.Context) {
	custodianGroupCfgPath := fmt.Sprintf("%s/custodian_groups.json", s.ConfigPath)
	custodianGroups, err := config.ReadJsonArrayConfig[CustodianGroupConfig](custodianGroupCfgPath)
	if err != nil {
		panic(err)
	}

	custodianGrs := make([]sqlc.CustodianGroup, 0)

	for _, custodianGroup := range custodianGroups {
		uid := sha3.Sum256([]byte(custodianGroup.Name))
		log.Info().Msgf("Custodian Group's UID: %x", custodianGroup.Uid)
		log.Info().Msgf("Custodian Group's BitcoinPubKey: %x", custodianGroup.BitcoinPubkey)
		if !bytes.Equal(uid[:], custodianGroup.Uid) {
			panic("custodian group uid is not correct")
		}

		custodians, err := json.Marshal(custodianGroup.Custodians)
		if err != nil {
			panic(err)
		}
		custodianGrs = append(custodianGrs, sqlc.CustodianGroup{
			ID:            custodianGroup.ID,
			Uid:           uid[:],
			Name:          custodianGroup.Name,
			BitcoinPubkey: custodianGroup.BitcoinPubkey,
			Quorum:        custodianGroup.Quorum,
			Custodians:    custodians,
		})

	}

	err = s.CombinedAdapter.SaveCustodianGroups(ctx, custodianGrs)
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

	chains := make([]string, 0)
	gatewayAddresses := make([][]byte, 0)

	for _, config := range configs {
		chains = append(chains, config.ID)
		add := funcs.Must(hex.DecodeString(strings.TrimPrefix(config.Gateway, "0x")))
		gatewayAddresses = append(gatewayAddresses, add)
	}

	err = s.CombinedAdapter.CreateGatewayAddresses(ctx, gatewayAddresses, chains)
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
				Chain:   evmClient.EvmConfig.ID,
				Symbol:  protocol.Symbol,
				Address: tokenAddr.Bytes(),
				ChainID: db.ConvertUint64ToNumeric(evmClient.EvmConfig.ChainID),
				Active:  true,
			}
			tokens = append(tokens, token)
		}
	}
	err := s.CombinedAdapter.SaveTokens(ctx, tokens)
	if err != nil {
		panic(err)
	}
}
