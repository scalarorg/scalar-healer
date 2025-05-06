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
		err := repo.SaveBridgeRequest(ctx, 222, common.Address{}, []byte{0x1, 0x2, 0x3}, []byte{0x1, 0x2, 0x3}, 100)
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
