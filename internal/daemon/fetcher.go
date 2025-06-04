package daemon

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/rs/zerolog/log"
	"github.com/scalarorg/scalar-healer/pkg/db"
	"github.com/scalarorg/scalar-healer/pkg/db/sqlc"
	"github.com/scalarorg/scalar-healer/pkg/utils"
)

// TODO: move to env for db
const (
	MEMPOOL_URL         = "https://mempool.space/testnet4/api"
	CONFIRMATION_HEIGHT = 2
)

var BITCOIN_PARAMS = &chaincfg.TestNet4Params

type MempoolUtxo struct {
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

func GetCurrentBlock() (uint64, error) {
	url := fmt.Sprintf("%s/blocks/tip/height", MEMPOOL_URL)
	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("failed to get current block: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response body: %w", err)
	}

	var currentBlock uint64
	if err := json.Unmarshal(body, &currentBlock); err != nil {
		return 0, fmt.Errorf("failed to unmarshal current block: %w", err)
	}

	return currentBlock, nil
}

func GetUtxoList(lockingScript []byte, custodianGrUid []byte) ([]sqlc.Utxo, []uint64, error) {
	blockHeight, err := GetCurrentBlock()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get current block: %w", err)
	}

	allowedBlockHeight := blockHeight - CONFIRMATION_HEIGHT + 1

	taprootAddress, err := utils.ScriptPubKeyToAddress(lockingScript, BITCOIN_PARAMS)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to convert script pubkey to address: %w", err)
	}

	url := fmt.Sprintf("%s/address/%s/utxo", MEMPOOL_URL, taprootAddress)

	resp, err := http.Get(url)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get UTXOs: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var utxos []MempoolUtxo
	if err := json.Unmarshal(body, &utxos); err != nil {
		return nil, nil, fmt.Errorf("failed to decode UTXOs: %w", err)
	}

	utxos = FilterUTXOsByBlockHeight(utxos, allowedBlockHeight)

	utxos = SortUTXOsByBlockHeight(utxos)
	for _, utxo := range utxos {
		log.Info().Msgf("[GetUtxoList] block height: %d, txid: %s, vout: %d, amount: %d", utxo.Status.BlockHeight, utxo.Txid, utxo.Vout, utxo.Value)
	}
	utxosList := []sqlc.Utxo{}
	blockHeights := make([]uint64, len(utxos))
	for _, utxo := range utxos {
		txID, err := hex.DecodeString(utxo.Txid)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to convert txid to hash: %w", err)
		}
		utxosList = append(utxosList, sqlc.Utxo{
			TxID:              txID,
			Vout:              int64(utxo.Vout),
			AmountInSats:      db.ConvertUint64ToNumeric(utxo.Value),
			CustodianGroupUid: custodianGrUid,
			ScriptPubkey:      lockingScript,
			BlockHeight:       int64(utxo.Status.BlockHeight),
		})
		blockHeights = append(blockHeights, utxo.Status.BlockHeight)
	}

	return utxosList, blockHeights, nil
}

func SortUTXOsByBlockHeight(utxos []MempoolUtxo) []MempoolUtxo {
	sort.Slice(utxos, func(i, j int) bool {
		return (utxos[i].Status.BlockHeight < utxos[j].Status.BlockHeight) ||
			(utxos[i].Status.BlockHeight == utxos[j].Status.BlockHeight && utxos[i].Txid < utxos[j].Txid) ||
			(utxos[i].Status.BlockHeight == utxos[j].Status.BlockHeight && utxos[i].Txid == utxos[j].Txid && utxos[i].Vout < utxos[j].Vout)
	})
	return utxos
}

func FilterUTXOsByBlockHeight(utxos []MempoolUtxo, allowedBlockHeight uint64) []MempoolUtxo {
	filteredUTXOs := []MempoolUtxo{}
	for _, utxo := range utxos {
		if utxo.Status.BlockHeight <= allowedBlockHeight {
			filteredUTXOs = append(filteredUTXOs, utxo)
		}
	}
	return filteredUTXOs
}
