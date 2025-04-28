package mongo

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
)

func (m *MongoRepository) GetNonce(ctx context.Context, address common.Address) uint64 {
	log.Info().Msgf("GetRedeemNonce: %s", address.Hex())
	return 0
}
