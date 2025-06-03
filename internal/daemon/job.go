package daemon

import (
	"context"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/scalarorg/scalar-healer/pkg/db/sqlc"
)

func (s *Service) DoJob(ctx context.Context) {
	// 1. Get all redeem sessions from db
	redeemSessions, err := s.CombinedAdapter.GetCompletedRedeemSessions(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to get all redeem sessions")
		return
	}

	// spawn a goroutine for each redeem session
	var wg sync.WaitGroup
	for _, session := range redeemSessions {
		wg.Add(1)
		go func(session sqlc.RedeemSession) {
			defer wg.Done()
			s.processRedeemSession(ctx, session)
		}(session)
	}
	wg.Wait()
}

func (s *Service) processRedeemSession(ctx context.Context, session sqlc.RedeemSession) {
	// 1. phase is PREPARING
	if session.CurrentPhase == sqlc.RedeemPhasePREPARING {
		log.Info().Msgf("redeem session already completed")
		// Collect all of reserved requests
		// Aggregate, create redeem commands
		// Timeout in <M> time then broadcast redeem commands -> switch phase to Executing
		// Allow user make request to reserve utxo
		return
	}

	// 2. phase is EXECUTING
	log.Info().Msgf("processing redeem session")

}
