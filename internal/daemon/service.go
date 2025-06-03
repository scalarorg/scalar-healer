package daemon

import (
	"context"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/scalarorg/scalar-healer/pkg/btc"
	"github.com/scalarorg/scalar-healer/pkg/db"
	"github.com/scalarorg/scalar-healer/pkg/electrs"
	"github.com/scalarorg/scalar-healer/pkg/evm"
	"github.com/scalarorg/scalar-healer/pkg/utils"
)

type Service struct {
	ConfigPath      string
	CombinedAdapter db.CombinedAdapter
	ElectrClient    *electrs.Client
	BtcClient       *btc.BtcClient
	EvmClients      []*evm.EvmClient
}

func NewService(configPath string, evmPrivKey string, dbAdapter db.CombinedAdapter) *Service {
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
		ConfigPath:      configPath,
		CombinedAdapter: dbAdapter,
		ElectrClient:    electrsClient,
		BtcClient:       btcClient,
		EvmClients:      evmClients,
	}
}

func (s *Service) Start(ctx context.Context) error {
	defer utils.Recover()

	s.initGenesis(ctx)
	log.Debug().Msg("[Service] starting")

	wg := sync.WaitGroup{}
	wg.Add(2)

	// go func() {
	// 	defer wg.Done()
	// 	s.ElectrClient.Start(ctx)
	// }()

	go func() {
		defer wg.Done()
		s.RecoverEvmSessions(ctx)
	}()

	wg.Wait()

	s.DoJob(ctx)

	select {
	case <-ctx.Done():
		return nil
	}
}

func (s *Service) Stop() {
	log.Info().Msg("Daemon service stopped")
}
