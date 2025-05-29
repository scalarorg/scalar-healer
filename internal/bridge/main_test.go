package bridge_test

import (
	"context"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/scalarorg/scalar-healer/cmd/api/server"
	"github.com/scalarorg/scalar-healer/pkg/db"
	"github.com/scalarorg/scalar-healer/pkg/db/healer"
	"github.com/scalarorg/scalar-healer/pkg/db/sqlc"
	testutils "github.com/scalarorg/scalar-healer/pkg/test_utils"
)

var (
	testServer *server.Server
	dbAdapter  db.HealderAdapter
)

func TestMain(m *testing.M) {
	var code int
	testutils.RunWithTestDB(func(ctx context.Context, repo db.HealderAdapter) error {
		// Setup test data
		gatewayAddr := common.HexToAddress("0x24a1dB57Fa3ecAFcbaD91d6Ef068439acEeAe090")
		pg := (repo).(*healer.HealerRepository)
		if pg == nil {
			panic("repo is not postgres")
		}

		err := pg.Queries.CreateGatewayAddress(ctx, sqlc.CreateGatewayAddressParams{
			Chain:   "evm|1",
			Address: gatewayAddr.Bytes(),
		})

		if err != nil {
			panic(err)
		}

		dbAdapter = repo
		testServer = server.New(repo)
		code = m.Run()
		return nil
	})
	os.Exit(code)
}

func cleanup() {
	(dbAdapter).(*healer.HealerRepository).TruncateTables(context.Background(), "bridge_requests", "nonces")
}
