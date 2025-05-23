package healer

import (
	"context"

	"github.com/scalarorg/scalar-healer/pkg/db/sqlc"
)

func (m *HealerRepository) SaveRedeemSessions(ctx context.Context, redeemSessions []sqlc.RedeemSession) error {
	var uids [][]byte
	var sequences []int64
	var chains []string
	var currentPhases []string

	for _, session := range redeemSessions {
		uids = append(uids, session.CustodianGroupUid)
		sequences = append(sequences, session.Sequence)
		chains = append(chains, session.Chain)
		currentPhases = append(currentPhases, string(session.CurrentPhase))
	}

	return m.Queries.SaveRedeemSessions(ctx, sqlc.SaveRedeemSessionsParams{
		Column1: uids,
		Column2: sequences,
		Column3: chains,
		Column4: currentPhases,
	})
}
