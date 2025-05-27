package healer_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/scalarorg/scalar-healer/constants"
	"github.com/scalarorg/scalar-healer/pkg/db"
	"github.com/scalarorg/scalar-healer/pkg/db/healer"
	"github.com/scalarorg/scalar-healer/pkg/db/sqlc"
	"github.com/scalarorg/scalar-healer/pkg/evm"
	testutils "github.com/scalarorg/scalar-healer/pkg/test_utils"
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
				for i := range batches {
					t.Logf("batches[%d]: %+v", i, batches[i])
					extraData, err := batches[i].GetExtraData()
					assert.NoError(t, err)
					t.Logf("extraData: %+v", extraData)
				}
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

				for i := range batches {
					t.Logf("batches[%d]: %+v", i, batches[i])
					extraData, err := batches[i].GetExtraData()
					assert.NoError(t, err)
					t.Logf("extraData: %+v", extraData)
				}
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
		})
	}
}

func TestSaveCommandBatches(t *testing.T) {

	mockExtraData := [][][]byte{
		{[]byte{0x1, 0x2}, []byte{0x3, 0x4}, []byte{0x5, 0x6}},
		{[]byte{0x7, 0x8}, []byte{0x9, 0xa}, []byte{0xb, 0xc}},
	}
	tests := []struct {
		name    string
		input   func() []*sqlc.CommandBatch
		wantErr bool
		check   func(t *testing.T, repo *healer.HealerRepository, input []*sqlc.CommandBatch)
	}{
		{
			name: "multiple batches",
			input: func() []*sqlc.CommandBatch {
				data1, err := json.Marshal(mockExtraData[0])
				assert.NoError(t, err)
				data2, err := json.Marshal(mockExtraData[1])
				assert.NoError(t, err)
				return []*sqlc.CommandBatch{
					{
						CommandBatchID: []byte{0x1},
						Chain:          "ethereum",
						Data:           []byte{0x2},
						SigHash:        []byte{0x3},
						Status:         pgtype.Int4{Int32: int32(sqlc.COMMAND_BATCH_STATUS_PENDING), Valid: true},
						ExtraData:      data1,
					},
					{
						CommandBatchID: []byte{0x5},
						Chain:          "polygon",
						Data:           []byte{0x6},
						SigHash:        []byte{0x7},
						Status:         pgtype.Int4{Int32: int32(sqlc.COMMAND_BATCH_STATUS_PENDING), Valid: true},
						ExtraData:      data2,
					},
				}
			},
			check: func(t *testing.T, repo *healer.HealerRepository, input []*sqlc.CommandBatch) {
				ctx := context.Background()
				for _, batch := range input {
					gotBatch, err := repo.GetCommandBatchByID(ctx, batch.CommandBatchID)
					t.Log("got batch: ", gotBatch)
					assert.NoError(t, err)
					assert.NotNil(t, gotBatch)
					assert.Equal(t, batch.ExtraData, gotBatch.ExtraData)
				}
			},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testutils.RunWithTestDB(func(ctx context.Context, repo db.HealderAdapter) error {
				pg := repo.(*healer.HealerRepository)

				i := tc.input()
				err := pg.SaveCommandBatches(ctx, i)
				assert.NoError(t, err)

				tc.check(t, pg, i)

				return nil
			})
		})
	}
}

func TestSaveCommands(t *testing.T) {
	tests := []struct {
		name    string
		input   []*sqlc.Command
		wantErr bool
	}{
		{
			name: "single command",
			input: []*sqlc.Command{
				{
					CommandID:   []byte{0x1},
					Chain:       "ethereum",
					Params:      []byte{0x2},
					Status:      sqlc.COMMAND_STATUS_PENDING.ToPgType(),
					CommandType: sqlc.CommandTypeSwitchPhase,
					Payload:     []byte{0x3},
				},
			},
			wantErr: false,
		},
		{
			name: "multiple commands",
			input: []*sqlc.Command{
				{
					CommandID:   []byte{0x2},
					Chain:       "ethereum",
					Params:      []byte{0x2},
					Status:      sqlc.COMMAND_STATUS_PENDING.ToPgType(),
					CommandType: sqlc.CommandTypeSwitchPhase,
					Payload:     []byte{0x3},
				},
				{
					CommandID:   []byte{0x3},
					Chain:       "polygon",
					Params:      []byte{0x4},
					Status:      sqlc.COMMAND_STATUS_PENDING.ToPgType(),
					CommandType: sqlc.CommandTypeSwitchPhase,
					Payload:     []byte{0x5},
				},
			},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testutils.RunWithTestDB(func(ctx context.Context, repo db.HealderAdapter) error {
				pg := repo.(*healer.HealerRepository)

				err := pg.SaveCommands(ctx, tc.input)
				if tc.wantErr {
					assert.Error(t, err)
					return nil
				}

				assert.NoError(t, err)

				// Verify saved data if no error expected
				if len(tc.input) > 0 {
					// Add verification logic here when query methods are available
				}
				return nil
			})
		})
	}
}

func TestSaveCommandsAndBatchCommandsTx(t *testing.T) {
	tests := []struct {
		name    string
		input   []sqlc.Command
		wantErr bool
		check   func(t *testing.T, repo *healer.HealerRepository)
	}{
		{
			name: "single command",
			input: []sqlc.Command{
				{
					CommandID:   evm.GetSignHash([]byte{1, 2, 3}).Bytes(),
					Chain:       "evm|11155111",
					Params:      []byte{0x2},
					Status:      sqlc.COMMAND_STATUS_PENDING.ToPgType(),
					CommandType: sqlc.CommandTypeMintToken,
					Payload:     []byte{0x3},
				},
			},
			wantErr: false,
			check: func(t *testing.T, repo *healer.HealerRepository) {
				ctx := context.Background()
				batches, err := repo.GetCommandBatches(ctx)
				assert.NoError(t, err)
				assert.Equal(t, 1, len(batches))
			},
		},
		{
			name: "multiple commands same chain",
			input: []sqlc.Command{
				{
					CommandID:   evm.GetSignHash([]byte{1, 2, 3}).Bytes(),
					Chain:       "evm|11155111",
					Params:      []byte{0x2},
					Status:      sqlc.COMMAND_STATUS_PENDING.ToPgType(),
					CommandType: sqlc.CommandTypeMintToken,
					Payload:     []byte{0x3},
				},
				{
					CommandID:   evm.GetSignHash([]byte{1, 2, 4}).Bytes(),
					Chain:       "evm|11155111",
					Params:      []byte{0x5},
					Status:      sqlc.COMMAND_STATUS_PENDING.ToPgType(),
					CommandType: sqlc.CommandTypeBurnToken,
					Payload:     []byte{0x6},
				},
			},
			wantErr: false,
			check: func(t *testing.T, repo *healer.HealerRepository) {
				ctx := context.Background()
				batches, err := repo.GetCommandBatches(ctx)
				assert.NoError(t, err)
				assert.Equal(t, 1, len(batches))

				// Verify extra data contains both payloads
				extraData, err := batches[0].GetExtraData()
				assert.NoError(t, err)
				assert.Equal(t, 2, len(extraData))
			},
		},
		{
			name: "multiple commands different chains",
			input: []sqlc.Command{
				{
					CommandID:   evm.GetSignHash([]byte{1, 2, 3}).Bytes(),
					Chain:       "evm|11155111",
					Params:      []byte{0x2},
					Status:      sqlc.COMMAND_STATUS_PENDING.ToPgType(),
					CommandType: sqlc.CommandTypeMintToken,
					Payload:     []byte{0x3},
				},
				{
					CommandID:   evm.GetSignHash([]byte{1, 2, 4}).Bytes(),
					Chain:       "evm|137",
					Params:      []byte{0x5},
					Status:      sqlc.COMMAND_STATUS_PENDING.ToPgType(),
					CommandType: sqlc.CommandTypeBurnToken,
					Payload:     []byte{0x6},
				},
			},
			wantErr: false,
			check: func(t *testing.T, repo *healer.HealerRepository) {
				ctx := context.Background()
				batches, err := repo.GetCommandBatches(ctx)
				assert.NoError(t, err)
				assert.Equal(t, 2, len(batches))

				// Each batch should have one payload
				for _, batch := range batches {
					extraData, err := batch.GetExtraData()
					assert.NoError(t, err)
					assert.Equal(t, 1, len(extraData))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testutils.RunWithTestDB(func(ctx context.Context, repo db.HealderAdapter) error {
				err := repo.SaveCommandsAndBatchCommandsTx(ctx, tt.input)
				if tt.wantErr {
					assert.Error(t, err)
					return nil
				}

				assert.NoError(t, err)
				tt.check(t, repo.(*healer.HealerRepository))
				return nil
			})
		})
	}
}
