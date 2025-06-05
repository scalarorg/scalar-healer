package evm

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/rs/zerolog/log"
	"github.com/scalarorg/bitcoin-vault/go-utils/types"
	"github.com/scalarorg/scalar-healer/config"
	contracts "github.com/scalarorg/scalar-healer/pkg/evm/contracts/generated"
)

const RETRY_INTERVAL = time.Second * 12

type Byte32 [32]uint8
type Bytes []byte
type EvmNetworkConfig struct {
	ChainID      uint64        `mapstructure:"chain_id"`
	ID           string        `mapstructure:"id"`
	Name         string        `mapstructure:"name"`
	RPCUrl       string        `mapstructure:"rpc_url"`
	AuthWeighted string        `mapstructure:"auth_weighted"`
	Gateway      string        `mapstructure:"gateway"`
	Finality     int           `mapstructure:"finality"`
	StartBlock   uint64        `mapstructure:"start_block"`
	PrivateKey   string        `mapstructure:"private_key"`
	GasLimit     uint64        `mapstructure:"gas_limit"`
	BlockTime    time.Duration `mapstructure:"blockTime"` //Timeout im ms for pending txs
	MaxRetry     int
	RecoverRange int64 `mapstructure:"recover_range"` //Max block range to recover events in single query
	RetryDelay   time.Duration
	TxTimeout    time.Duration `mapstructure:"tx_timeout"` //Timeout for send txs (~3s)
}

func (c *EvmNetworkConfig) GetChainId() uint64 {
	return c.ChainID
}
func (c *EvmNetworkConfig) GetId() string {
	return c.ID
}
func (c *EvmNetworkConfig) GetName() string {
	return c.Name
}
func (c *EvmNetworkConfig) GetFamily() string {
	return types.ChainTypeEVM.String()
}

type EvmClient struct {
	EvmConfig      *EvmNetworkConfig
	Client         *ethclient.Client
	ChainName      string
	GatewayAddress common.Address
	Gateway        *contracts.IScalarGateway
	transactOpts   *bind.TransactOpts
	retryInterval  time.Duration
}

func NewEvmClients(configPath string, evmPrivKey string) ([]*EvmClient, error) {
	evmCfgPath := fmt.Sprintf("%s/evm.json", configPath)
	configs, err := config.ReadJsonArrayConfig[EvmNetworkConfig](evmCfgPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read electrum configs: %w", err)
	}

	evmClients := make([]*EvmClient, 0, len(configs))
	for _, evmConfig := range configs {
		//Inject evm private keys
		evmConfig.PrivateKey = evmPrivKey
		//Set default value for block time if is not set
		if evmConfig.BlockTime == 0 {
			evmConfig.BlockTime = 12 * time.Second
		} else {
			evmConfig.BlockTime = evmConfig.BlockTime * time.Millisecond
		}
		//Set default gaslimit to 300000
		if evmConfig.GasLimit == 0 {
			evmConfig.GasLimit = 3000000
		}
		if evmConfig.RecoverRange == 0 {
			evmConfig.RecoverRange = 1000000
		}
		client, err := NewEvmClient(&evmConfig)
		if err != nil {
			log.Warn().Msgf("Failed to create evm client for %s: %v", evmConfig.GetName(), err)
			continue
		}
		evmClients = append(evmClients, client)
	}

	return evmClients, nil
}

func NewEvmClient(evmConfig *EvmNetworkConfig) (*EvmClient, error) {
	// Setup
	ctx := context.Background()
	log.Info().Any("evmConfig", evmConfig).Msgf("[EvmClient] [NewEvmClient] connecting to EVM network")
	// Connect to a test network
	rpc, err := rpc.DialContext(ctx, evmConfig.RPCUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to EVM network %s: %w", evmConfig.Name, err)
	}
	client := ethclient.NewClient(rpc)
	gateway, gatewayAddress, err := CreateGateway(evmConfig.Name, evmConfig.Gateway, client)
	if err != nil {
		return nil, fmt.Errorf("failed to create gateway for network %s: %w", evmConfig.Name, err)
	}
	auth, err := CreateTransactOpts(evmConfig)
	if err != nil {
		//Not fatal, we can still use the gateway without auth
		//auth is only used for sending transaction
		log.Warn().Msgf("[EvmClient] [NewEvmClient] failed to create auth for network %s: %v", evmConfig.Name, err)
		panic(err)
	}
	evmClient := &EvmClient{
		EvmConfig:      evmConfig,
		Client:         client,
		GatewayAddress: *gatewayAddress,
		Gateway:        gateway,
		transactOpts:   auth,
		retryInterval: RETRY_INTERVAL,
	}

	return evmClient, nil
}

func CreateGateway(networName string, gwAddr string, client *ethclient.Client) (*contracts.IScalarGateway, *common.Address, error) {
	if gwAddr == "" {
		return nil, nil, fmt.Errorf("gateway address is not set for network %s", networName)
	}
	gatewayAddress := common.HexToAddress(gwAddr)
	gateway, err := contracts.NewIScalarGateway(gatewayAddress, client)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize gateway contract for network %s: %w", networName, err)
	}
	return gateway, &gatewayAddress, nil
}

func CreateTransactOpts(evmConfig *EvmNetworkConfig) (*bind.TransactOpts, error) {
	if evmConfig.PrivateKey == "" {
		return nil, fmt.Errorf("private key is not set for network %s", evmConfig.Name)
	}
	privateKey, err := crypto.HexToECDSA(evmConfig.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key for network %s: %w", evmConfig.Name, err)
	}
	chainID := big.NewInt(int64(evmConfig.ChainID))
	transactOpts, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create auth for network %s: %w", evmConfig.Name, err)
	}
	transactOpts.GasLimit = evmConfig.GasLimit
	return transactOpts, nil
}

func (ec *EvmClient) CreateCallOpts() (*bind.CallOpts, error) {
	callOpt := &bind.CallOpts{
		From:    ec.transactOpts.From,
		Context: context.Background(),
	}
	return callOpt, nil
}

func (c *EvmClient) SetAuth(auth *bind.TransactOpts) {
	c.transactOpts = auth
}
