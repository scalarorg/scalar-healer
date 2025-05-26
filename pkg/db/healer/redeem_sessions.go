package healer

import (
	"context"
	"encoding/hex"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/scalarorg/scalar-healer/pkg/db/sqlc"
)

func (m *HealerRepository) SaveRedeemSessionAndChainRedeemSessionsTx(ctx context.Context, chainRedeemSessions []sqlc.ChainRedeemSession) ([]sqlc.ChainRedeemSession, error) {
	var outdatedSession []sqlc.ChainRedeemSession

	err := m.execTx(ctx, func(q *sqlc.Queries) error {
		// 1. Group redeem sessions by group id
		sessionsByGroup := make(map[string][]sqlc.ChainRedeemSession)
		for _, session := range chainRedeemSessions {
			cusGrUID := hex.EncodeToString(session.CustodianGroupUid)
			if _, ok := sessionsByGroup[cusGrUID]; !ok {
				sessionsByGroup[cusGrUID] = make([]sqlc.ChainRedeemSession, 0)
			}
			sessionsByGroup[cusGrUID] = append(sessionsByGroup[cusGrUID], session)
		}

		redeemSessions := make([]sqlc.RedeemSession, 0)

		// 2. Compare the latest redeem session with the chain redeem session
		for _, sessions := range sessionsByGroup {
			lastestRedeemSessionMap := make(map[int]bool)
			latestRedeemSession := sessions[0]
			lastestRedeemSessionMap[0] = true
			for index, session := range sessions {
				if session.Cmp(&latestRedeemSession) > 0 {
					latestRedeemSession = session
					lastestRedeemSessionMap = make(map[int]bool)
					lastestRedeemSessionMap[index] = true
				} else if session.Cmp(&latestRedeemSession) == 0 {
					lastestRedeemSessionMap[index] = true
				}
			}
			// 3. Save the outdated redeem session
			for index, session := range sessions {
				if _, ok := lastestRedeemSessionMap[index]; !ok {
					outdatedSession = append(outdatedSession, session)
				}
			}

			// 4. Collect the redeem session for group
			redeemSessions = append(redeemSessions, sqlc.RedeemSession{
				CustodianGroupUid: latestRedeemSession.CustodianGroupUid,
				Sequence:          latestRedeemSession.Sequence,
				CurrentPhase:      latestRedeemSession.CurrentPhase,
				IsSwitching: pgtype.Bool{
					Bool:  len(outdatedSession) > 0,
					Valid: true,
				},
			})
		}

		// 5. Save the chain redeem session
		err := m.saveChainRedeemSessions(ctx, chainRedeemSessions)
		if err != nil {
			return err
		}

		// 6. Save the redeem session
		err = m.saveRedeemSessions(ctx, redeemSessions)
		if err != nil {
			return err
		}

		return nil
	})

	return outdatedSession, err
}

func (m *HealerRepository) saveRedeemSessions(ctx context.Context, redeemSessions []sqlc.RedeemSession) error {
	var uids [][]byte
	var sequences []int64
	var currentPhases []string
	var isSwitchings []bool

	for _, session := range redeemSessions {
		uids = append(uids, session.CustodianGroupUid)
		sequences = append(sequences, session.Sequence)
		currentPhases = append(currentPhases, string(session.CurrentPhase))
		isSwitchings = append(isSwitchings, session.IsSwitching.Bool)
	}

	return m.Queries.SaveRedeemSessions(ctx, sqlc.SaveRedeemSessionsParams{
		Column1: uids,
		Column2: sequences,
		Column3: currentPhases,
		Column4: isSwitchings,
	})
}

func (m *HealerRepository) saveChainRedeemSessions(ctx context.Context, redeemSessions []sqlc.ChainRedeemSession) error {
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

func (m *HealerRepository) GetRedeemSession(ctx context.Context, groupUid []byte) (*sqlc.RedeemSession, error) {
	result, err := m.Queries.GetRedeemSession(ctx, groupUid)
	return &result, err
}

func (m *HealerRepository) GetChainRedeemSession(ctx context.Context, grUID []byte, chain string) (*sqlc.ChainRedeemSession, error) {
	result, err := m.Queries.GetChainRedeemSession(ctx, sqlc.GetChainRedeemSessionParams{
		Chain:             chain,
		CustodianGroupUid: grUID,
	})
	return &result, err
}
