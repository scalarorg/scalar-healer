package healer_test

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/scalarorg/scalar-healer/pkg/db"
	"github.com/scalarorg/scalar-healer/pkg/db/healer"
	"github.com/scalarorg/scalar-healer/pkg/db/sqlc"
	testutils "github.com/scalarorg/scalar-healer/pkg/test_utils"
	"github.com/stretchr/testify/assert"
)

func TestCreateBridge(t *testing.T) {
	testutils.RunWithTestDB(func(ctx context.Context, repo db.HealderAdapter) error {
		err := repo.SaveBridgeRequest(ctx, "evm|222", common.Address{}, []byte{0x1, 0x2, 0x3}, []byte{0x1, 0x2, 0x3}, 0)
		if err != nil {
			t.Errorf("failed to save bridge request: %v", err)
		}

		bridgeRequests, err := repo.ListBridgeRequests(ctx, common.Address{}, 0, 10)
		if err != nil {
			t.Errorf("failed to get bridge requests: %v", err)
		}

		if len(bridgeRequests) != 1 {
			t.Errorf("expected 1 bridge request, got %d", len(bridgeRequests))
		}

		nonce := repo.GetNonce(ctx, common.Address{})
		if nonce != 1 {
			t.Errorf("Expect %d, got %d", 1, nonce)
		}
		return nil
	})
}

func TestCreateGatewayAddress(t *testing.T) {
	testutils.RunWithTestDB(func(ctx context.Context, repo db.HealderAdapter) error {
		err := (repo).(*healer.HealerRepository).CreateGatewayAddress(ctx, sqlc.CreateGatewayAddressParams{
			Chain:   "evm|1",
			Address: common.MaxAddress.Bytes(),
		})
		if err != nil {
			t.Errorf("failed to save create gateway address: %v", err)
		}
		return nil
	})
}

func TestCreateTokens(t *testing.T) {
	testutils.RunWithTestDB(func(ctx context.Context, repo db.HealderAdapter) error {
		// TODO: add query to add protocol first
		err := repo.SaveTokens(ctx, []sqlc.Token{
			{
				Symbol:  "ETH",
				ChainID: db.ConvertUint64ToNumeric(1),
				Address: common.MaxAddress.Bytes(),
				Active:  true,
			},
		})
		if err != nil {
			t.Errorf("failed to save tokens: %v", err)
		}

		pg := (repo).(*healer.HealerRepository)

		allTokens, err := pg.ListTokens(ctx)
		if err != nil {
			t.Errorf("failed to list tokens: %v", err)
		}

		t.Logf("all tokens: %v", allTokens)

		return nil
	})
}

func TestCreateCustodianGroups(t *testing.T) {
	testutils.RunWithTestDB(func(ctx context.Context, repo db.HealderAdapter) error {
		mockUid := []byte{0x1, 0x2, 0x3}
		mockName := "test"
		mockQuorum := int64(1)
		mockBitcoinPubkey := common.MaxHash[:]
		mockCustodians := [][]byte{common.MaxAddress.Bytes(), common.MaxAddress.Bytes()}
		_ = mockCustodians

		err := repo.SaveCustodianGroups(ctx, []sqlc.CustodianGroup{
			{
				Uid:           mockUid,
				Name:          mockName,
				Quorum:        mockQuorum,
				BitcoinPubkey: mockBitcoinPubkey,
			},
		})

		if err != nil {
			t.Errorf("failed to save custodian groups: %v", err)
		}

		pg := (repo).(*healer.HealerRepository)

		grs, err := pg.GetAllCustodianGroups(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(grs))

		gr, err := pg.GetCustodianGroupByUID(ctx, []byte{0x1, 0x2, 0x3})
		assert.NoError(t, err)
		assert.Equal(t, "test", gr.Name)
		assert.Equal(t, 1, int(gr.Quorum))
		assert.Equal(t, common.MaxHash, common.BytesToHash(gr.BitcoinPubkey))
		return nil
	})
}

func TestSaveProtocols(t *testing.T) {
	testutils.RunWithTestDB(func(ctx context.Context, repo db.HealderAdapter) error {
		err := repo.SaveProtocols(ctx, []sqlc.Protocol{
			{
				Symbol:             "BTC",
				Name:               "Bitcoin",
				CustodianGroupName: "test",
				CustodianGroupUid:  []byte{0x1, 0x2, 0x3},
				Tag:                "test",
				LiquidityModel:     "test",
				Decimals:           int64(8),
				Capacity:           db.ConvertUint64ToNumeric(100000000),
				DailyMintLimit:     db.ConvertUint64ToNumeric(100000000),
				Avatar:             "",
			},
		})

		if err != nil {
			t.Errorf("failed to save protocols: %v", err)
		}

		protocol, err := repo.GetProtocol(ctx, "BTC")
		if err != nil {
			t.Errorf("failed to get protocol: %v", err)
		}
		assert.Equal(t, "Bitcoin", protocol.Name)
		assert.Equal(t, "test", protocol.CustodianGroupName)
		assert.Equal(t, []byte{0x1, 0x2, 0x3}, protocol.CustodianGroupUid)
		assert.Equal(t, "test", protocol.Tag)
		assert.Equal(t, "test", protocol.LiquidityModel)
		assert.Equal(t, "BTC", protocol.Symbol)
		assert.Equal(t, int64(8), protocol.Decimals)
		return nil
	})
}

func TestSaveUtxoSnapshot(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(t *testing.T, repo *healer.HealerRepository)
		input   []sqlc.Utxo
		wantErr bool
	}{
		{
			name:    "empty snapshot",
			setup:   func(t *testing.T, repo *healer.HealerRepository) {},
			input:   []sqlc.Utxo{},
			wantErr: false,
		},
		{
			name:  "valid snapshot with same custodian group",
			setup: func(t *testing.T, repo *healer.HealerRepository) {},
			input: []sqlc.Utxo{
				{
					TxID:              []byte{0x1},
					Vout:              1,
					ScriptPubkey:      []byte{0x2},
					AmountInSats:      pgtype.Numeric{Int: common.Big1, Valid: true},
					CustodianGroupUid: []byte{0x3},
					BlockHeight:       100,
				},
				{
					TxID:              []byte{0x2},
					Vout:              1,
					ScriptPubkey:      []byte{0x02},
					AmountInSats:      pgtype.Numeric{Int: common.Big2, Valid: true},
					CustodianGroupUid: []byte{0x3},
					BlockHeight:       100,
				},
			},
			wantErr: false,
		},
		{
			name:  "invalid snapshot with mixed groups",
			setup: func(t *testing.T, repo *healer.HealerRepository) {},
			input: []sqlc.Utxo{
				{
					TxID:              []byte{0x1},
					Vout:              1,
					ScriptPubkey:      []byte{0x2},
					AmountInSats:      pgtype.Numeric{Int: common.Big1, Valid: true},
					CustodianGroupUid: []byte{0x3},
					BlockHeight:       100,
				},
				{
					TxID:              []byte{0x1},
					Vout:              1,
					ScriptPubkey:      []byte{0x2},
					AmountInSats:      pgtype.Numeric{Int: common.Big2, Valid: true},
					CustodianGroupUid: []byte{0x4}, // Different group
					BlockHeight:       100,
				},
			},
			wantErr: true,
		},

		{
			name:  "invalid snapshot with mixed script pubkeys",
			setup: func(t *testing.T, repo *healer.HealerRepository) {},
			input: []sqlc.Utxo{
				{
					TxID:              []byte{0x1},
					Vout:              1,
					ScriptPubkey:      []byte{0x1},
					AmountInSats:      pgtype.Numeric{Int: common.Big1, Valid: true},
					CustodianGroupUid: []byte{0x4},
					BlockHeight:       100,
				},
				{
					TxID:              []byte{0x1},
					Vout:              1,
					ScriptPubkey:      []byte{0x2},
					AmountInSats:      pgtype.Numeric{Int: common.Big2, Valid: true},
					CustodianGroupUid: []byte{0x4}, // Different group
					BlockHeight:       100,
				},
			},
			wantErr: true,
		},

		{
			name: "block height validation",
			setup: func(t *testing.T, repo *healer.HealerRepository) {
				// Insert existing UTXO with higher block height
				err := repo.SaveUTXOs(context.Background(), sqlc.SaveUTXOsParams{
					Column1: [][]byte{{0x1}},
					Column2: []int64{1},
					Column3: [][]byte{{0x2}},
					Column4: []pgtype.Numeric{{Int: common.Big1, Valid: true}},
					Column5: [][]byte{{0x3}},
					Column6: []int64{200}, // Higher block height
				})
				assert.NoError(t, err)
			},
			input: []sqlc.Utxo{
				{
					TxID:              []byte{0x1},
					Vout:              1,
					ScriptPubkey:      []byte{0x2},
					AmountInSats:      pgtype.Numeric{Int: common.Big1, Valid: true},
					CustodianGroupUid: []byte{0x3},
					BlockHeight:       100, // Lower block height
				},
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testutils.RunWithTestDB(func(ctx context.Context, repo db.HealderAdapter) error {
				pg := repo.(*healer.HealerRepository)
				tc.setup(t, pg)

				err := pg.SaveUtxoSnapshot(ctx, tc.input)
				if tc.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)

					// Verify saved data if no error expected
					if len(tc.input) > 0 {
						utxos, err := pg.GetUTXOsByCustodianGroupUID(ctx, tc.input[0].CustodianGroupUid)
						assert.NoError(t, err)
						assert.Equal(t, len(tc.input), len(utxos))
					}
				}
				return nil
			})
		})
	}
}
