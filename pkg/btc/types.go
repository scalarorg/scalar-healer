package btc

import (
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"github.com/scalarorg/bitcoin-vault/go-utils/types"
)

const COMPONENT_NAME = "BtcClient"

type SigningType int

const (
	CUSTODIAL_ONLY SigningType = iota
	USER_PROTOCOL
	USER_CUSTODIAL
	PROTOCOL_CUSTODIAL
)

type BtcClient struct {
	btcConfig *BtcNetworkConfig
	client    *rpcclient.Client
}

type BtcClientInterface interface {
	SendTx(tx *wire.MsgTx, maxFeeRate *float64) (*chainhash.Hash, error)
	TestMempoolAccept(txs []*wire.MsgTx, maxFeeRatePerKb float64) ([]*btcjson.TestMempoolAcceptResult, error)
}

type BtcNetworkConfig struct {
	Network    string  `mapstructure:"network"`
	ID         string  `mapstructure:"id"`
	ChainID    uint64  `mapstructure:"chain_id"`
	Name       string  `mapstructure:"name"`
	Type       string  `mapstructure:"type"`
	Host       string  `mapstructure:"host"`
	Port       int     `mapstructure:"port"`
	User       string  `mapstructure:"user"`
	Password   string  `mapstructure:"password"`
	SSL        *bool   `mapstructure:"ssl,omitempty"`
	MempoolUrl string  `mapstructure:"mempool_url,omitempty"`
	PrivateKey string  `mapstructure:"private_key,omitempty"`
	Address    *string `mapstructure:"address,omitempty"` //Taproot address. Todo: set it as parameter from scalar note
	MaxFeeRate float64 `mapstructure:"max_fee_rate,omitempty"`
}

func (c *BtcNetworkConfig) GetChainId() uint64 {
	return c.ChainID
}
func (c *BtcNetworkConfig) GetId() string {
	return c.ID
}
func (c *BtcNetworkConfig) GetName() string {
	return c.Name
}
func (c *BtcNetworkConfig) GetFamily() string {
	return types.ChainTypeBitcoin.String()
}

type Utxo struct {
	Txid   string `json:"txid"`
	Vout   uint32 `json:"vout"`
	Status struct {
		Confirmed   bool   `json:"confirmed"`
		BlockHeight uint64 `json:"block_height"`
		BlockHash   string `json:"block_hash"`
		BlockTime   uint64 `json:"block_time"`
	} `json:"status"`
	Value uint64 `json:"value"`
}
