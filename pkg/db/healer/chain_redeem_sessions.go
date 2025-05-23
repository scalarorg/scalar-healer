package healer

import (
	"context"

	"github.com/scalarorg/scalar-healer/pkg/db/sqlc"
)

func (m *HealerRepository) SaveChainRedeemSessions(ctx context.Context, redeemSessions []sqlc.ChainRedeemSession) error {
	var chains []string
	var uids [][]byte
	var sequences []int64
	var currentPhases []string

	for _, session := range redeemSessions {
		chains = append(chains, session.Chain)
		uids = append(uids, session.CustodianGroupUid)
		sequences = append(sequences, session.Sequence)
		currentPhases = append(currentPhases, string(session.CurrentPhase))
	}

	return m.Queries.SaveChainRedeemSessions(ctx, sqlc.SaveChainRedeemSessionsParams{
		Column1: chains,
		Column2: uids,
		Column3: sequences,
		Column4: currentPhases,
	})
}
