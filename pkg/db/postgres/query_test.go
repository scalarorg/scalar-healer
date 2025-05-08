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
