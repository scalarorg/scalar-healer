package healer_test

import (
	"fmt"
	"testing"

	"github.com/scalarorg/scalar-healer/constants"
	"github.com/scalarorg/scalar-healer/pkg/db/healer"
	"github.com/scalarorg/scalar-healer/pkg/db/sqlc"
	"github.com/scalarorg/scalar-healer/pkg/evm"
	"github.com/stretchr/testify/assert"
)

func TestNewCommandBatches(t *testing.T) {
	tests := []struct {
		name    string
		input   []*sqlc.Command
		wantErr bool
		check   func(t *testing.T, batches []*sqlc.CommandBatch, cmds []*sqlc.Command)
	}{
		{
			name: "single chain, single batch",
			input: []*sqlc.Command{
				{
					Chain:       "evm|11155111",
					CommandID:   evm.GetSignHash([]byte{1, 2, 3}).Bytes(),
					Params:      []byte("params"),
					CommandType: sqlc.CommandTypeSwitchPhase,
					Status:      sqlc.COMMAND_STATUS_PENDING.ToPgType(),
					Payload:     []byte("payload"),
				},
			},
			wantErr: false,
			check: func(t *testing.T, batches []*sqlc.CommandBatch, cmds []*sqlc.Command) {
				assert.Equal(t, 1, len(batches))
				assert.Equal(t, batches[0].CommandBatchID, cmds[0].CommandBatchID)
			},
		},
		{
			name: "single chain, multiple batches",
			input: func() []*sqlc.Command {
				cmds := make([]*sqlc.Command, constants.BATCH_SIZE+1)
				for i := range cmds {
					cmds[i] = &sqlc.Command{
						Chain:       "evm|97",
						CommandID:   evm.GetSignHash([]byte{byte(i)}).Bytes(),
						Params:      []byte(fmt.Sprintf("params_%d", i)),
						CommandType: sqlc.CommandTypeSwitchPhase,
						Status:      sqlc.COMMAND_STATUS_PENDING.ToPgType(),
						Payload:     []byte(fmt.Sprintf("payload_%d", i)),
					}
				}
				return cmds
			}(),
			wantErr: false,
			check: func(t *testing.T, batches []*sqlc.CommandBatch, cmds []*sqlc.Command) {
				assert.Equal(t, 2, len(batches))
				assert.Equal(t, batches[0].CommandBatchID, cmds[0].CommandBatchID)
				assert.Equal(t, batches[1].CommandBatchID, cmds[constants.BATCH_SIZE].CommandBatchID)
			},
		},
		{
			name: "multiple chains",
			input: []*sqlc.Command{
				{
					Chain:     "evm|80",
					ID:        1,
					CommandID: evm.GetSignHash([]byte{1, 2, 3}).Bytes(),
				},
				{
					Chain:     "evm|81",
					ID:        2,
					CommandID: evm.GetSignHash([]byte{4, 5, 6}).Bytes(),
				},
			},
			wantErr: false,
			check: func(t *testing.T, batches []*sqlc.CommandBatch, cmds []*sqlc.Command) {
				assert.Equal(t, 2, len(batches))
				assert.Equal(t, batches[0].Chain, cmds[0].Chain)
				assert.Equal(t, batches[1].Chain, cmds[1].Chain)
			},
		},
		{
			name: "invalid chain",
			input: []*sqlc.Command{
				{
					Chain: "invalid_chain",
					ID:    1,
				},
			},
			wantErr: true,
			check: func(t *testing.T, batches []*sqlc.CommandBatch, cmds []*sqlc.Command) {
				assert.Equal(t, 0, len(batches))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			batches, err := healer.NewCommandBatches(tc.input)
			if tc.wantErr {
				assert.Error(t, err)
			}
			tc.check(t, batches, tc.input)

			// for _, batch := range batches {
			// 	for _, cmd := range batch.Commands {
			// 		assert.Equal(t, batch.Chain, cmd.Chain)
			// 	}
			// }
		})
	}
}
