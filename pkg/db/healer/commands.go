package healer

import (
	"context"

	"github.com/scalarorg/data-models/chains"
	"github.com/scalarorg/scalar-healer/pkg/db/sqlc"
)

func (m *HealerRepository) UpdateEvmCommandExecuted(ctx context.Context, cmdExecuted *chains.CommandExecuted) error {
	return nil
}

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
