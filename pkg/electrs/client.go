package electrs

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/scalarorg/bitcoin-vault/go-utils/chain"
	"github.com/scalarorg/data-models/scalarnet"
	"github.com/scalarorg/go-electrum/electrum"
	"github.com/scalarorg/scalar-healer/config"
	"github.com/scalarorg/scalar-healer/pkg/db"
)

type Client struct {
	electrumConfig *ElectrsConfig
	Electrs        *electrum.Client
	dbAdapter      db.HealderAdapter
	currentHeight  int
}

func NewElectrumClient(configPath string, dbAdapter db.HealderAdapter) (*Client, error) {
	electrumCfgPath := fmt.Sprintf("%s/electrum.json", configPath)
	configs, err := config.ReadJsonArrayConfig[ElectrsConfig](electrumCfgPath)
	if err != nil || len(configs) == 0 {
		return nil, fmt.Errorf("failed to read electrum configs: %w", err)
	}
	config := configs[0]
	if config.Host == "" {
		return nil, fmt.Errorf("electrum rpc host is required")
	}
	if config.Port == 0 {
		return nil, fmt.Errorf("electrum rpc port is required")
	}
	if dbAdapter == nil {
		return nil, fmt.Errorf("dbAdapter is required")
	}
	if config.Confirmations == 0 {
		config.Confirmations = 1
	}
	rpcEndpoint := fmt.Sprintf("%s:%d", config.Host, config.Port)
	electrs, err := electrum.Connect(&electrum.Options{
		Dial: func() (net.Conn, error) {
			return net.DialTimeout("tcp", rpcEndpoint, time.Second)
		},
		MethodTimeout:   time.Second,
		PingInterval:    -1,
		SoftwareVersion: "scalar-relayer",
	})
	if err != nil {
		log.Error().Err(err).Msgf("Failed to connect to electrum server at %s", rpcEndpoint)
		return nil, err
	}
	return &Client{
		electrumConfig: &config,
		Electrs:        electrs,
		dbAdapter:      dbAdapter,
	}, nil
}

func (c *Client) Start(ctx context.Context) error {
	log.Info().Msg("[ElectrsClient] start")
	params := []interface{}{}
	//Set batch size from config or default value
	params = append(params, c.electrumConfig.BatchSize)

	lastCheckpoint := c.getLastCheckpoint(ctx)
	log.Debug().Msgf("[ElectrumClient] [Start] Last checkpoint: %v", lastCheckpoint)
	if lastCheckpoint.EventKey != "" {
		params = append(params, lastCheckpoint.EventKey)
	} else if c.electrumConfig.LastVaultTx != "" {
		params = append(params, c.electrumConfig.LastVaultTx)
	}
	log.Debug().Msg("[ElectrumClient] [Start] Subscribing to new block event for request to confirm if vault transaction is get enought confirmation")
	c.Electrs.BlockchainHeaderSubscribe(ctx, c.BlockchainHeaderHandler(ctx))
	log.Debug().Msgf("[ElectrumClient] [Start] Subscribing to vault transactions with params: %v", params)
	c.Electrs.VaultTransactionSubscribe(ctx, c.VaultTxMessageHandler(ctx), params...)

	return nil
}

func (c *Client) GetSymbol(chainInfo *chain.ChainInfo, tokenAddress string) (string, error) {
	//TODO: implement this
	return "", nil
}

// Get lastcheck point from db, return default value if not found
func (c *Client) getLastCheckpoint(ctx context.Context) *scalarnet.EventCheckPoint {
	sourceChain := c.electrumConfig.SourceChain
	eventName := config.EVENT_ELECTRS_VAULT_TRANSACTION
	lastCheckpoint, err := c.dbAdapter.GetLastEventCheckPoint(ctx, sourceChain, eventName, 0)
	if err != nil {
		log.Warn().Str("chainId", sourceChain).
			Str("eventName", eventName).
			Msg("[ElectrumClient] [getLastCheckpoint] using default value")
	}
	if lastCheckpoint == nil {
		lastCheckpoint = &scalarnet.EventCheckPoint{
			ChainName:   sourceChain,
			EventName:   eventName,
			BlockNumber: 0,
			TxHash:      "",
			LogIndex:    0,
			EventKey:    "",
		}
	}
	return lastCheckpoint
}
