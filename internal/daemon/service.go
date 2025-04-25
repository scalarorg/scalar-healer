package daemon

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/scalarorg/scalar-healer/config"
	"github.com/scalarorg/scalar-healer/pkg/btc"
	"github.com/scalarorg/scalar-healer/pkg/db"
	"github.com/scalarorg/scalar-healer/pkg/electrs"
	"github.com/scalarorg/scalar-healer/pkg/evm"
)

type Service struct {
	DbAdapter    db.DbAdapter
	ElectrClient *electrs.Client
	BtcClient    *btc.BtcClient
	EvmClients   []*evm.EvmClient
}

func NewService(config *config.Config, dbAdapter db.DbAdapter) *Service {
	electrsClient, err := electrs.NewElectrumClient(config.ConfigPath, dbAdapter)
	if err != nil {
		panic(err)
	}
	btcClient, err := btc.NewBtcClient(config.ConfigPath)
	if err != nil {
		panic(err)
	}
	evmClients, err := evm.NewEvmClients(config.ConfigPath, dbAdapter)
	if err != nil {
		panic(err)
	}
	return &Service{
		DbAdapter:    dbAdapter,
		ElectrClient: electrsClient,
		BtcClient:    btcClient,
		EvmClients:   evmClients,
	}
}

func (s *Service) Start(ctx context.Context) error {
	log.Debug().Msg("[Service] start")
	//Start electrum clients. This client can get all vault transactions from last checkpoint of begining if no checkpoint is found
	go s.ElectrClient.Start(ctx)

	groups, err := s.DbAdapter.GetAllCustodianGroups()
	if err != nil {
		log.Error().Err(err).Msg("[DaemonService] Cannot get custodian groups")
	}
	err = s.RecoverEvmSessions(groups)
	if err != nil {
		log.Warn().Err(err).Msgf("[DaemonService] cannot recover evm sessions")
		panic(err)
	}
	return nil
}
