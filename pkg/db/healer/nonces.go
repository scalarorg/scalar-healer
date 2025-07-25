package healer

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/scalarorg/scalar-healer/pkg/db/sqlc"
)

func (m *HealerRepository) GetNonce(ctx context.Context, address common.Address) uint64 {
	num, err := m.Queries.GetNonce(ctx, address.Bytes())
	if err != nil {
		return 0
	}
	return sqlc.ConvertNumericToUint64(num) + 1
}
