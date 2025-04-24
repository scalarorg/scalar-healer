package electrs

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/scalarorg/data-models/scalarnet"
	"github.com/scalarorg/go-electrum/electrum"
	"github.com/scalarorg/scalar-healer/config"
	"github.com/scalarorg/scalar-healer/pkg/db"
)

type Client struct {
	electrumConfig *ElectrsConfig
	Electrs        *electrum.Client
	dbAdapter      db.DbAdapter
	currentHeight  int
}

func NewElectrumClients(config *ElectrsConfig, dbAdapter db.DbAdapter) (*Client, error) {
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
		electrumConfig: config,
		Electrs:        electrs,
		dbAdapter:      dbAdapter,
	}, nil
}

func (c *Client) Start(ctx context.Context) error {
	params := []interface{}{}
	//Set batch size from config or default value
	params = append(params, c.electrumConfig.BatchSize)

	lastCheckpoint := c.getLastCheckpoint()
	log.Debug().Msgf("[ElectrumClient] [Start] Last checkpoint: %v", lastCheckpoint)
	if lastCheckpoint.EventKey != "" {
		params = append(params, lastCheckpoint.EventKey)
	} else if c.electrumConfig.LastVaultTx != "" {
		params = append(params, c.electrumConfig.LastVaultTx)
	}
	log.Debug().Msg("[ElectrumClient] [Start] Subscribing to new block event for request to confirm if vault transaction is get enought confirmation")
	c.Electrs.BlockchainHeaderSubscribe(ctx, c.BlockchainHeaderHandler)
	log.Debug().Msgf("[ElectrumClient] [Start] Subscribing to vault transactions with params: %v", params)
	c.Electrs.VaultTransactionSubscribe(ctx, c.VaultTxMessageHandler, params...)

	return nil
}

func (c *Client) GetSymbol(chainId string, tokenAddress string) (string, error) {
	//TODO: implement this
	return "", nil
}

// Get lastcheck point from db, return default value if not found
func (c *Client) getLastCheckpoint() *scalarnet.EventCheckPoint {
	sourceChain := c.electrumConfig.SourceChain
	lastCheckpoint, err := c.dbAdapter.GetLastEventCheckPoint(sourceChain, config.EVENT_ELECTRS_VAULT_TRANSACTION, 0)
	if err != nil {
		log.Warn().Str("chainId", sourceChain).
			Str("eventName", config.EVENT_ELECTRS_VAULT_TRANSACTION).
			Msg("[ElectrumClient] [getLastCheckpoint] using default value")
	}
	return lastCheckpoint
}
