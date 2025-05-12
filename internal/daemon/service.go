package daemon

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/scalarorg/scalar-healer/pkg/btc"
	"github.com/scalarorg/scalar-healer/pkg/db"
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
	defer utils.Recover()

	s.initGenesis(ctx)
	log.Debug().Msg("[Service] starting")
	// go s.ElectrClient.Start(ctx)

	go s.RecoverEvmSessions(ctx)

	// s.ProcessMissingLogs(ctx)

	select {
	case <-ctx.Done():
		return nil

	}
}

func (s *Service) Stop() {
	log.Info().Msg("Daemon service stopped")
}
