package healer

import (
	"context"
	"encoding/hex"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/scalarorg/scalar-healer/pkg/db/sqlc"
)

func (m *HealerRepository) SaveRedeemSessionAndChainRedeemSessionsTx(ctx context.Context, chainRedeemSessions []sqlc.ChainRedeemSession) (map[string][]sqlc.ChainRedeemSessionUpdate, error) {
	var outdatedSessionsByGroup map[string][]sqlc.ChainRedeemSessionUpdate

	expiredAt := time.Now().Add(time.Hour * 2).Unix()

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
			outdatedSessions := make([]sqlc.ChainRedeemSessionUpdate, 0)

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
					outdatedSessions = append(outdatedSessions, sqlc.ChainRedeemSessionUpdate{
						Chain:             session.Chain,
						CustodianGroupUid: session.CustodianGroupUid,
						Sequence:          session.Sequence,
						CurrentPhase:      session.CurrentPhase,
						NewPhase:          latestRedeemSession.CurrentPhase,
					})
				}
			}

			if len(outdatedSessions) > 0 {
				if outdatedSessionsByGroup == nil {
					outdatedSessionsByGroup = make(map[string][]sqlc.ChainRedeemSessionUpdate)
				}
				outdatedSessionsByGroup[hex.EncodeToString(latestRedeemSession.CustodianGroupUid)] = outdatedSessions
			}

			// 4. Collect the redeem session for group
			rs := sqlc.RedeemSession{
				CustodianGroupUid: latestRedeemSession.CustodianGroupUid,
				Sequence:          latestRedeemSession.Sequence,
				CurrentPhase:      latestRedeemSession.CurrentPhase,
				IsSwitching: pgtype.Bool{
					Bool:  len(outdatedSessions) > 0,
					Valid: true,
				},
			}

			if len(outdatedSessions) == 0 {
				rs.PhaseExpiredAt = expiredAt
			}

			redeemSessions = append(redeemSessions, rs)
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

	return outdatedSessionsByGroup, err
}

func (m *HealerRepository) saveRedeemSessions(ctx context.Context, redeemSessions []sqlc.RedeemSession) error {
	var uids [][]byte
	var sequences []int64
	var currentPhases []string
	var isSwitchings []bool
	var expiredAts []int64

	for _, session := range redeemSessions {
		uids = append(uids, session.CustodianGroupUid)
		sequences = append(sequences, session.Sequence)
		currentPhases = append(currentPhases, string(session.CurrentPhase))
		isSwitchings = append(isSwitchings, session.IsSwitching.Bool)
		expiredAts = append(expiredAts, session.PhaseExpiredAt)
	}

	return m.Queries.SaveRedeemSessions(ctx, sqlc.SaveRedeemSessionsParams{
		Column1: uids,
		Column2: sequences,
		Column3: currentPhases,
		Column4: isSwitchings,
		Column5: expiredAts,
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
