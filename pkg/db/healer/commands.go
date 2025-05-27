package healer

import (
	"context"
	"encoding/binary"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/scalarorg/scalar-healer/constants"
	"github.com/scalarorg/scalar-healer/pkg/db/sqlc"
	"github.com/scalarorg/scalar-healer/pkg/evm"
	"github.com/scalarorg/scalar-healer/pkg/utils/chains"
)

func (m *HealerRepository) SaveCommands(ctx context.Context, commands []sqlc.Command) error {
	var (
		commandIds [][]byte
		chains     []string
		params     [][]byte
		status     []int32
		types      []sqlc.CommandType
	)

	for _, command := range commands {
		commandIds = append(commandIds, command.CommandID)
		chains = append(chains, command.Chain)
		params = append(params, command.Params)
		status = append(status, command.Status.Int32)
		types = append(types, command.CommandType)
	}

	return m.Queries.SaveCommands(ctx, sqlc.SaveCommandsParams{
		Column1: commandIds,
		Column2: chains,
		Column3: params,
		Column4: status,
		Column5: types,
	})
}

func (m *HealerRepository) SaveCommandsAndBatchCommandsTx(ctx context.Context, commands []sqlc.Command) error {
	return m.execTx(ctx, func(q *sqlc.Queries) error {
		err := m.SaveCommands(ctx, commands)
		if err != nil {
			return err
		}

		return nil

	})
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

	return &sqlc.CommandBatch{
		CommandBatchID: commandBatchID,
		Data:      data,
		SigHash:   evm.GetSignHash(data).Bytes(),
		Status:    sqlc.COMMAND_BATCH_STATUS_PENDING.ToPgType(),
		Chain:     chain,
		ExtraData: extraData,
	}, nil
}
