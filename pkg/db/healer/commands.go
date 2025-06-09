package healer

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/scalarorg/scalar-healer/constants"
	"github.com/scalarorg/scalar-healer/pkg/db/sqlc"
	"github.com/scalarorg/scalar-healer/pkg/evm"
	"github.com/scalarorg/scalar-healer/pkg/utils/chains"
	"github.com/scalarorg/scalar-healer/pkg/utils/slices"
)

func (m *HealerRepository) SaveCommands(ctx context.Context, commands []*sqlc.Command) error {
	var (
		commandIds [][]byte
		chains     []string
		params     [][]byte
		status     []int32
		types      []string
		payloads   [][]byte
	)

	for _, command := range commands {
		commandIds = append(commandIds, command.CommandID)
		chains = append(chains, command.Chain)
		params = append(params, command.Params)
		status = append(status, command.Status.Int32)
		types = append(types, command.CommandType.String())
		payloads = append(payloads, command.Payload)
	}

	return m.Queries.SaveCommands(ctx, sqlc.SaveCommandsParams{
		Column1: commandIds,
		Column2: chains,
		Column3: params,
		Column4: status,
		Column5: types,
		Column6: payloads,
	})
}

func (m *HealerRepository) SaveCommandBatches(ctx context.Context, commandBatches []*sqlc.CommandBatch) error {
	var (
		commandBatchIds [][]byte
		chains          []string
		bunchData       [][]byte
		sigHashes       [][]byte
		status          []int32
		bunchExtraData  [][]byte
	)

	for _, commandBatch := range commandBatches {
		commandBatchIds = append(commandBatchIds, commandBatch.CommandBatchID)
		chains = append(chains, commandBatch.Chain)
		bunchData = append(bunchData, commandBatch.Data)
		sigHashes = append(sigHashes, commandBatch.SigHash)
		status = append(status, commandBatch.Status.Int32)
		bunchExtraData = append(bunchExtraData, commandBatch.ExtraData)
	}

	return m.Queries.SaveCommandBatches(ctx, sqlc.SaveCommandBatchesParams{
		Column1: commandBatchIds,
		Column2: chains,
		Column3: bunchData,
		Column4: sigHashes,
		Column5: status,
		Column6: bunchExtraData,
	})
}

func (m *HealerRepository) SaveCommandsAndBatchCommandsTx(ctx context.Context, commands []sqlc.Command) error {
	return m.execTx(ctx, func(cv context.Context, q *sqlc.Queries) error {
		var err error

		cmds := slices.Map(commands, func(cmd sqlc.Command) *sqlc.Command {
			return &cmd
		})

		batches, err := NewCommandBatches(cmds)
		if err != nil {
			return err
		}

		err = m.SaveCommands(ctx, cmds)
		if err != nil {
			return err
		}

		err = m.SaveCommandBatches(ctx, batches)
		if err != nil {
			return err
		}

		return nil
	})
}

func (m *HealerRepository) GetCommandBatches(ctx context.Context) ([]sqlc.CommandBatch, error) {
	return m.Queries.GetCommandBatches(ctx)
}

func (m *HealerRepository) GetCommandBatchByID(ctx context.Context, commandBatchID []byte) (sqlc.CommandBatch, error) {
	return m.Queries.GetCommandBatchByID(ctx, commandBatchID)
}

func NewCommandBatches(cmds []*sqlc.Command) ([]*sqlc.CommandBatch, error) {
	// 1. Group by chain first
	chainCmdMap := make(map[string][]*sqlc.Command)
	for _, cmd := range cmds {
		chainCmdMap[cmd.Chain] = append(chainCmdMap[cmd.Chain], cmd)
	}

	// 2. Create command batch for each chain, collect updated commands
	var commandBatches []*sqlc.CommandBatch
	for chain, cmds := range chainCmdMap {
		chainID, err := chains.ChainName(chain).GetChainID()
		if err != nil {
			return nil, err
		}

		for i := 0; i < len(cmds); i += constants.BATCH_SIZE {
			end := i + constants.BATCH_SIZE
			if end > len(cmds) {
				end = len(cmds)
			}
			batch, err := newCommandBatch(chain, chainID, cmds[i:end])
			if err != nil {
				return nil, err
			}
			commandBatches = append(commandBatches, batch)
		}
	}

	return commandBatches, nil
}

func newCommandBatch(chain string, chainID *big.Int, cmds []*sqlc.Command) (*sqlc.CommandBatch, error) {
	var commandIDs []sqlc.CommandID
	var commands []sqlc.CommandType
	var commandParams [][]byte

	var extraData [][]byte

	for _, cmd := range cmds {
		commandIDs = append(commandIDs, sqlc.CommandID(cmd.CommandID))
		commands = append(commands, cmd.CommandType)
		commandParams = append(commandParams, cmd.Params)
		extraData = append(extraData, cmd.Payload)
	}

	data, err := evm.PackArguments(chainID, commandIDs, commands, commandParams)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, uint64(now.Unix()))

	commandBatchID := crypto.Keccak256(bz, data)

	// update command batch id
	for _, cmd := range cmds {
		cmd.CommandBatchID = commandBatchID
	}

	encodedExtraData, err := json.Marshal(extraData)
	if err != nil {
		return nil, err
	}

	return &sqlc.CommandBatch{
		CommandBatchID: commandBatchID,
		Data:           data,
		SigHash:        evm.GetSignHash(data).Bytes(),
		Status:         sqlc.COMMAND_BATCH_STATUS_PENDING.ToPgType(),
		Chain:          chain,
		ExtraData:      encodedExtraData,
	}, nil
}
