package daemon

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/scalarorg/scalar-healer/pkg/btc"
	"github.com/scalarorg/scalar-healer/pkg/db"
	"github.com/scalarorg/scalar-healer/pkg/electrs"
	"github.com/scalarorg/scalar-healer/pkg/evm"
)

type Service struct {
	ConfigPath   string
	DbAdapter    db.DbAdapter
	ElectrClient *electrs.Client
	BtcClient    *btc.BtcClient
	EvmClients   []*evm.EvmClient
}

func NewService(configPath string, evmPrivKey string, dbAdapter db.DbAdapter) *Service {
	electrsClient, err := electrs.NewElectrumClient(configPath, dbAdapter)
	if err != nil {
		panic(err)
	}
	btcClient, err := btc.NewBtcClient(configPath)
	if err != nil {
		panic(err)
	}
	evmClients, err := evm.NewEvmClients(configPath, evmPrivKey, dbAdapter)
	if err != nil {
		panic(err)
	}
	return &Service{
		ConfigPath:   configPath,
		DbAdapter:    dbAdapter,
		ElectrClient: electrsClient,
		BtcClient:    btcClient,
		EvmClients:   evmClients,
	}
}

func (s *Service) Start(ctx context.Context) error {
	s.initGenesis(ctx)
	log.Debug().Msg("[Service] starting")
	//Start electrum clients. This client can get all vault transactions from last checkpoint of begining if no checkpoint is found
	go s.ElectrClient.Start(ctx)

	groups, err := s.DbAdapter.GetAllCustodianGroups(ctx)
	if err != nil {
		log.Error().Err(err).Msg("[DaemonService] Cannot get custodian groups")
	}
	err = s.RecoverEvmSessions(ctx, groups)
	if err != nil {
		log.Warn().Err(err).Msgf("[DaemonService] cannot recover evm sessions")
		panic(err)
	}
	return nil
}
