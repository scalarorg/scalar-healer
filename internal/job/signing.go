package job

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
	"github.com/scalarorg/scalar-core/x/chains/types"
	"github.com/scalarorg/scalar-healer/pkg/db"
	"github.com/scalarorg/scalar-healer/pkg/db/sqlc"
	"github.com/scalarorg/scalar-healer/pkg/evm"
	"github.com/scalarorg/scalar-healer/pkg/tofnd"
)

type result struct {
	CmdID  []byte
	Musigs []tofnd.Musig
	Data   []byte
}

func HandlePendingSigning(repo db.HealderAdapter, tofndManager *tofnd.Manager) func() {
	return func() {

		log.Info().Msg("Starting signing redeem commands")

		ctx := context.Background()
		commands, err := repo.ListPendingSigningRedeemCommands(ctx)
		if err != nil {
			log.Error().Err(err).Msg("Failed to list pending signing redeem commands")
			return
		}

		log.Info().Msgf("Found %d pending signing redeem commands", len(commands))

		wg := sync.WaitGroup{}

		results := make(chan result, len(commands))

		for _, cmd := range commands {
			wg.Add(1)
			go func(cmd sqlc.RedeemCommand) {
				defer wg.Done()
				musigs, err := tofndManager.Sign(ctx, cmd.SigHash)
				if err != nil {
					log.Error().Err(err).Msg("Failed to sign redeem command")
					return
				}
				results <- result{CmdID: cmd.ID, Musigs: musigs, Data: cmd.Data}
			}(cmd)
		}

		wg.Wait()
		close(results)

		signatures := make([][]byte, 0, len(commands))
		ids := make([][]byte, 0, len(commands))
		execute_data := make([][]byte, 0, len(commands))

		for result := range results {
			json, err := json.Marshal(result.Musigs)
			if err != nil {
				log.Error().Err(err).Msg("Failed to marshal signature")
				return
			}
			signatures = append(signatures, json)
			ids = append(ids, result.CmdID)

			ex, err := getExecuteDataAndSigs(result.Data, result.Musigs)
			if err != nil {
				log.Error().Err(err).Msg("Failed to get execute data")
				return
			}

			execute_data = append(execute_data, ex)
		}

		err = repo.SubmitRedeemCommandSignatures(ctx, ids, signatures, execute_data)
		if err != nil {
			log.Error().Err(err).Msg("Failed to submit redeem command signatures")
			return
		}
		log.Info().Msg("Submitted redeem command signatures")
	}
}

func getExecuteDataAndSigs(data []byte, musigs []tofnd.Musig) ([]byte, error) {
	addresses, weights, threshold, signatures, err := getProof(musigs)
	if err != nil {
		return nil, err
	}

	return createExecuteDataMultisig(data, addresses, weights, threshold, signatures)
}

func getProof(musigs []tofnd.Musig) ([]common.Address, []*big.Int, *big.Int, [][]byte, error) {
	if len(musigs) == 0 {
		return nil, nil, nil, nil, errors.New("no musigs")
	}

	addresses := make([]common.Address, 0, len(musigs))
	weights := make([]*big.Int, 0, len(musigs))
	threshold := big.NewInt(int64(musigs[0].Threshold))
	signatures := make([][]byte, 0, len(musigs))

	for _, musig := range musigs {
		addresses = append(addresses, crypto.PubkeyToAddress(*musig.Pubkey))
		weights = append(weights, big.NewInt(int64(musig.Weight)))
		signatures = append(signatures, musig.Sig.ToHomesteadSig())
	}

	return addresses, weights, threshold, signatures, nil
}

func createExecuteDataMultisig(data []byte, addresses []common.Address, weights []*big.Int, threshold *big.Int, signatures [][]byte) ([]byte, error) {
	proof, err := getWeightedSignaturesProof(addresses, weights, threshold, signatures)
	if err != nil {
		return nil, err
	}

	executeData, err := evm.ExecuteDataArguments.Pack(data, proof)
	if err != nil {
		return nil, err
	}

	return evm.ScalarGatewayABI.Pack(types.ScalarGatewayFuncExecute, executeData)
}

func getWeightedSignaturesProof(addresses []common.Address, weights []*big.Int, threshold *big.Int, signatures [][]byte) ([]byte, error) {
	proof, err := evm.CommandBatchArguments.Pack(addresses, weights, threshold, signatures)
	if err != nil {
		return nil, err
	}

	return proof, nil
}
