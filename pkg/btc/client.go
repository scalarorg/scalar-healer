package btc

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"github.com/rs/zerolog/log"
	"github.com/scalarorg/scalar-healer/config"
)

func NewBtcClient(configPath string) (*BtcClient, error) {
	btcConfigPath := fmt.Sprintf("%s/btc.json", configPath)
	btcConfigs, err := config.ReadJsonArrayConfig[BtcNetworkConfig](btcConfigPath)
	if err != nil {
		panic(err)
	}
	if len(btcConfigs) == 0 {
		panic("no btc config found")
	}
	btcConfig := btcConfigs[0]
	if btcConfig.MempoolUrl == "" {
		panic(fmt.Sprintf("mempool url is not set for %s", btcConfig.Name))
	}
	// Set max fee rate to 0.10 if not set
	if btcConfig.MaxFeeRate == 0 {
		btcConfig.MaxFeeRate = 0.10
	}
	connCfg := &rpcclient.ConnConfig{
		Host:         fmt.Sprintf("%s:%d", btcConfig.Host, btcConfig.Port),
		User:         btcConfig.User,
		Pass:         btcConfig.Password,
		HTTPPostMode: true,
		DisableTLS:   btcConfig.SSL == nil || !*btcConfig.SSL,
	}

	// Create new client
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create BTC client for network %s: %w", btcConfig.Network, err)
	}
	btcClient := &BtcClient{
		btcConfig: &btcConfig,
		client:    client,
	}
	return btcClient, nil
}

// GetTransaction retrieves detailed information about a transaction given its ID
// func (c *BtcClient) GetTransaction(txID string) (*btcjson.GetTransactionResult, error) {
// 	result, err := c.client.RawRequest("gettransaction", []json.RawMessage{[]byte(`"` + txID + `"`)})
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get transaction info: %w", err)
// 	}

// 	var tx btcjson.GetTransactionResult
// 	if err := json.Unmarshal(result, &tx); err != nil {
// 		return nil, fmt.Errorf("failed to unmarshal transaction: %w", err)
// 	}
// 	return &tx, nil
// }

func (c *BtcClient) Config() *BtcNetworkConfig {
	return c.btcConfig
}

func (c *BtcClient) BroadcastTx(tx *wire.MsgTx, maxFeeRate *float64) (*chainhash.Hash, error) {
	// If testnet4, create Command then call c.RpcClient.SendCmd(cmd)
	if c.btcConfig.Network == "testnet4" {
		rawTx, err := CreateRawTx(tx)
		if err != nil {
			return nil, err
		}
		log.Debug().Msgf("Send rawTx: %s\n", rawTx)
		cmd := c.creatSendRawTransactionCmd(rawTx, maxFeeRate)
		res := c.client.SendCmd(cmd)
		// Cast the response to FutureTestMempoolAcceptResult and call Receive
		future := rpcclient.FutureSendRawTransactionResult(res)
		return future.Receive()
	} else {
		// Otherwise, call c.RpcClient.SendRawTransaction(tx, true)x
		return c.client.SendRawTransaction(tx, true)
	}
}
func (c *BtcClient) BroadcastRawTx(signedPsbtHex string) (*chainhash.Hash, error) {
	cmd := c.creatSendRawTransactionCmd(signedPsbtHex, &c.btcConfig.MaxFeeRate)
	res := c.client.SendCmd(cmd)
	// Cast the response to FutureTestMempoolAcceptResult and call Receive
	future := rpcclient.FutureSendRawTransactionResult(res)
	return future.Receive()
}

func (c *BtcClient) FindBroadcastedTx(signedPsbtHex string) (*chainhash.Hash, bool, error) {
	// First decode the hex string to MsgTx
	msgTx := wire.NewMsgTx(wire.TxVersion)
	txHash := msgTx.TxHash()
	txBytes, err := hex.DecodeString(signedPsbtHex)
	if err != nil {
		return nil, false, err
	}
	if err := msgTx.Deserialize(bytes.NewReader(txBytes)); err != nil {
		return &txHash, false, err
	}

	tx, err := c.client.GetTransaction(&txHash)
	if err != nil {
		return &txHash, false, err
	}
	log.Debug().Msgf("[BtcClient] FindBroadcastedTx: %v", tx)
	return &txHash, true, nil
}

func (c *BtcClient) TestMempoolRawTx(signedPsbtHex string) (*chainhash.Hash, error) {
	// First decode the hex string to MsgTx
	msgTx := wire.NewMsgTx(wire.TxVersion)
	txBytes, err := hex.DecodeString(signedPsbtHex)
	if err != nil {
		return nil, err
	}
	if err := msgTx.Deserialize(bytes.NewReader(txBytes)); err != nil {
		return nil, err
	}

	// Test mempool accept
	results, err := c.TestMempoolAccept([]*wire.MsgTx{msgTx}, c.btcConfig.MaxFeeRate)
	if err != nil {
		return nil, err
	}

	// If accepted, return the transaction hash
	if len(results) > 0 && results[0].Allowed {
		hash := msgTx.TxHash()
		return &hash, nil
	}

	return nil, fmt.Errorf("transaction rejected: %v", results[0].RejectReason)
}

// if maxFeeRate is not nil, set the feeSetting parameter
// otherwise, don't set the feeSetting parameter use default value which is set by bitcoind 0.10
func (c *BtcClient) creatSendRawTransactionCmd(rawTxHex string, maxFeeRate *float64) *btcjson.SendRawTransactionCmd {
	if maxFeeRate != nil {
		return btcjson.NewBitcoindSendRawTransactionCmd(rawTxHex, *maxFeeRate)
	}
	return &btcjson.SendRawTransactionCmd{
		HexTx:      rawTxHex,
		FeeSetting: nil,
	}
}
func (c *BtcClient) TestMempoolAccept(txs []*wire.MsgTx, maxFeeRatePerKb float64) ([]*btcjson.TestMempoolAcceptResult, error) {
	if c.btcConfig.Network == "testnet4" {
		// Add some checks to make sure the txs are valid
		rawTxns, err := CreateRawTxs(txs)
		if err != nil {
			return nil, err
		}
		res := c.client.SendCmd(btcjson.NewTestMempoolAcceptCmd(rawTxns, maxFeeRatePerKb))
		// Cast the response to FutureTestMempoolAcceptResult and call Receive
		future := rpcclient.FutureTestMempoolAcceptResult(res)
		return future.Receive()
	} else {
		return c.client.TestMempoolAccept(txs, maxFeeRatePerKb)
	}
}

func (c *BtcClient) GetTxOut(txHex string, vout uint32) (*btcjson.GetTxOutResult, error) {
	log.Debug().Msgf("GetTxOut: %s, vout: %d", txHex, vout)
	var txHash chainhash.Hash
	txBytes, err := hex.DecodeString(strings.TrimPrefix(txHex, "0x"))
	if err != nil {
		return nil, err
	}
	if len(txBytes) != chainhash.HashSize {
		return nil, fmt.Errorf("txBytes length is not %d", chainhash.HashSize)
	}
	for i := 0; i < chainhash.HashSize/2; i++ {
		txHash[i], txHash[chainhash.HashSize-1-i] = txBytes[chainhash.HashSize-1-i], txBytes[i]
	}
	txOut, err := c.client.GetTxOut(&txHash, vout, true)
	return txOut, err
}
func CreateRawTx(tx *wire.MsgTx) (string, error) {
	// Serialize the transaction and convert to hex string.
	buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
	// TODO(yy): add similar checks found in `BtcDecode` to
	// `BtcEncode` - atm it just serializes bytes without any
	// bitcoin-specific checks.
	if err := tx.Serialize(buf); err != nil {
		return "", err
	}
	// Sanity check the provided tx is valid, which can be removed
	// once we have similar checks added in `BtcEncode`.
	//
	// NOTE: must be performed after buf.Bytes is copied above.
	//
	// TODO(yy): remove it once the above TODO is addressed.
	// if err := tx.Deserialize(buf); err != nil {
	// 	err = fmt.Errorf("%w: %v", rpcclient.ErrInvalidParam, err)
	// 	return "", err
	// }
	return hex.EncodeToString(buf.Bytes()), nil
}

func CreateRawTxs(txns []*wire.MsgTx) ([]string, error) {
	// Iterate all the transactions and turn them into hex strings.
	rawTxns := make([]string, 0, len(txns))
	for _, tx := range txns {
		rawTx, err := CreateRawTx(tx)
		if err != nil {
			return nil, err
		}
		rawTxns = append(rawTxns, rawTx)

	}

	return rawTxns, nil
}
