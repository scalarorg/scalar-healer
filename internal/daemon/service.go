package daemon

import (
	"context"
	"encoding/hex"

	"github.com/rs/zerolog/log"
	"github.com/scalarorg/scalar-healer/pkg/btc"
	"github.com/scalarorg/scalar-healer/pkg/db"
	"github.com/scalarorg/scalar-healer/pkg/db/sqlc"
	"github.com/scalarorg/scalar-healer/pkg/electrs"
	"github.com/scalarorg/scalar-healer/pkg/evm"
	"github.com/scalarorg/scalar-healer/pkg/utils"
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
	err = s.RecoverEvmSessions(ctx, utils.Map(groups, func(group sqlc.CustodianGroup) string {
		return hex.EncodeToString(group.Uid)
	}))
	if err != nil {
		log.Warn().Err(err).Msgf("[DaemonService] cannot recover evm sessions")
		panic(err)
	}
	//Recover all events
	for _, client := range s.EvmClients {
		go client.ProcessMissingLogs()
		go func() {
			//Todo: Handle the moment when recover just finished and listner has not started yet. It around 1 second
			err := client.RecoverAllEvents(ctx, groups)
			if err != nil {
				log.Warn().Err(err).Msgf("[Relayer] [Start] cannot recover events for evm client %s", client.EvmConfig.GetId())
			} else {
				log.Info().Msgf("[Relayer] [Start] recovered missing events for evm client %s", client.EvmConfig.GetId())
				client.Start(ctx)
			}
		}()
	}
	return nil
}

func (s *Service) Stop() {
	log.Info().Msg("Daemon service stopped")
}
