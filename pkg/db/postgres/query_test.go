package postgres_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/scalarorg/scalar-healer/pkg/db"
	"github.com/scalarorg/scalar-healer/pkg/db/postgres"
	"github.com/scalarorg/scalar-healer/pkg/db/sqlc"
	testutils "github.com/scalarorg/scalar-healer/pkg/test_utils"
	"github.com/stretchr/testify/assert"
)

func TestCreateBridge(t *testing.T) {
	testutils.RunWithTestDB(func(ctx context.Context, repo db.DbAdapter) error {
		err := repo.SaveBridgeRequest(ctx, 222, common.Address{}, []byte{0x1, 0x2, 0x3}, []byte{0x1, 0x2, 0x3}, 0)
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
	testutils.RunWithTestDB(func(ctx context.Context, repo db.DbAdapter) error {
		err := (repo).(*postgres.PostgresRepository).CreateGatewayAddress(ctx, sqlc.CreateGatewayAddressParams{
			ChainID: pgtype.Numeric{
				Int:   big.NewInt(1),
				Exp:   0,
				NaN:   false,
				Valid: true,
			},
			Address: common.MaxAddress.Bytes(),
		})
		if err != nil {
			t.Errorf("failed to save create gateway address: %v", err)
		}
		return nil
	})
}

func TestCreateTokens(t *testing.T) {
	testutils.RunWithTestDB(func(ctx context.Context, repo db.DbAdapter) error {
		err := repo.SaveTokens(ctx, []sqlc.Token{
			{
				Symbol:   "ETH",
				ChainID:  db.ConvertUint64ToNumeric(1),
				Protocol: "SCALAR",
				Address:  common.MaxAddress.Bytes(),
				Name:     "Ethereum",
				Decimal:  db.ConvertUint64ToNumeric(8),
				Avatar:   "",
				Active:   true,
			},
		})
		if err != nil {
			t.Errorf("failed to save tokens: %v", err)
		}

		pg := (repo).(*postgres.PostgresRepository)

		allTokens, err := pg.ListTokens(ctx)
		if err != nil {
			t.Errorf("failed to list tokens: %v", err)
		}

		t.Logf("all tokens: %v", allTokens)

		return nil
	})
}

func TestCreateCustodianGroups(t *testing.T) {
	testutils.RunWithTestDB(func(ctx context.Context, repo db.DbAdapter) error {
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

		pg := (repo).(*postgres.PostgresRepository)

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
	testutils.RunWithTestDB(func(ctx context.Context, repo db.DbAdapter) error {
		err := repo.SaveProtocols(ctx, []sqlc.Protocol{
			{
				Asset:              "BTC",
				Name:               "Bitcoin",
				CustodianGroupName: "test",
				CustodianGroupUid:  []byte{0x1, 0x2, 0x3},
				Tag:                "test",
				LiquidityModel:     "test",
				Symbol:             "BTC",
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
		assert.Equal(t, "BTC", protocol.Asset)
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
