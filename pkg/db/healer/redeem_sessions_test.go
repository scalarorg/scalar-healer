package healer_test

import (
	"context"
	"testing"

	"github.com/scalarorg/scalar-healer/pkg/db"
	"github.com/scalarorg/scalar-healer/pkg/db/healer"
	"github.com/scalarorg/scalar-healer/pkg/db/sqlc"
	testutils "github.com/scalarorg/scalar-healer/pkg/test_utils"
	"github.com/stretchr/testify/assert"
)

func TestSaveRedeemSessionAndChainRedeemSessionsTx(t *testing.T) {

	mockCustodianGroup := sqlc.CustodianGroup{
		Uid:           []byte{0x1},
		Name:          "scalarv36",
		BitcoinPubkey: []byte{0x1},
		Quorum:        3,
	}

	tests := []struct {
		name    string
		input   []sqlc.ChainRedeemSession
		wantErr bool
		check   func(t *testing.T, repo *healer.HealerRepository, outdated map[string][]sqlc.ChainRedeemSessionUpdate)
	}{
		{
			name: "single group with no outdated sessions",
			input: []sqlc.ChainRedeemSession{
				{
					CustodianGroupUid: []byte{0x1},
					Sequence:          1,
					CurrentPhase:      sqlc.RedeemPhasePREPARING,
					Chain:             "ethereum",
				},
			},
			wantErr: false,
			check: func(t *testing.T, repo *healer.HealerRepository, outdated map[string][]sqlc.ChainRedeemSessionUpdate) {
				session, err := repo.GetRedeemSession(context.Background(), mockCustodianGroup.Uid)
				assert.NoError(t, err)
				assert.Empty(t, outdated)

				assert.Equal(t, mockCustodianGroup.Uid, session.CustodianGroupUid)
				assert.Equal(t, int64(1), session.Sequence)
				assert.Equal(t, sqlc.RedeemPhasePREPARING, session.CurrentPhase)

				chainSession, err := repo.GetChainRedeemSession(context.Background(), mockCustodianGroup.Uid, "ethereum")
				assert.NoError(t, err)
				assert.Equal(t, mockCustodianGroup.Uid, chainSession.CustodianGroupUid)
				assert.Equal(t, int64(1), chainSession.Sequence)
				assert.Equal(t, sqlc.RedeemPhasePREPARING, chainSession.CurrentPhase)
				assert.Equal(t, "ethereum", chainSession.Chain)
			},
		},
		{
			name: "single group with outdated sessions",
			input: []sqlc.ChainRedeemSession{
				{
					CustodianGroupUid: []byte{0x1},
					Sequence:          2,
					CurrentPhase:      sqlc.RedeemPhasePREPARING,
					Chain:             "ethereum",
				},
				{
					CustodianGroupUid: []byte{0x1},
					Sequence:          1,
					CurrentPhase:      sqlc.RedeemPhaseEXECUTING,
					Chain:             "bnb",
				},
			},
			wantErr: false,
			check: func(t *testing.T, repo *healer.HealerRepository, sessions map[string][]sqlc.ChainRedeemSessionUpdate) {
				assert.Len(t, sessions, 1)
				outdated := sessions["01"]
				assert.Equal(t, int64(1), outdated[0].Sequence)

				session, err := repo.GetRedeemSession(context.Background(), mockCustodianGroup.Uid)
				assert.NoError(t, err)
				assert.NotEmpty(t, outdated)
				assert.Equal(t, mockCustodianGroup.Uid, session.CustodianGroupUid)
				assert.Equal(t, int64(2), session.Sequence)
				assert.Equal(t, sqlc.RedeemPhasePREPARING, session.CurrentPhase)
				assert.Equal(t, true, session.IsSwitching.Bool)
				chainSession, err := repo.GetChainRedeemSession(context.Background(), mockCustodianGroup.Uid, "ethereum")
				assert.NoError(t, err)
				assert.Equal(t, mockCustodianGroup.Uid, chainSession.CustodianGroupUid)
				assert.Equal(t, int64(2), chainSession.Sequence)
				assert.Equal(t, sqlc.RedeemPhasePREPARING, chainSession.CurrentPhase)
				assert.Equal(t, "ethereum", chainSession.Chain)
				chainSession, err = repo.GetChainRedeemSession(context.Background(), mockCustodianGroup.Uid, "bnb")
				assert.NoError(t, err)
				assert.Equal(t, mockCustodianGroup.Uid, chainSession.CustodianGroupUid)
				assert.Equal(t, int64(1), chainSession.Sequence)
				assert.Equal(t, sqlc.RedeemPhaseEXECUTING, chainSession.CurrentPhase)
				assert.Equal(t, "bnb", chainSession.Chain)
			},
		},
		{
			name: "single group without outdated sessions",
			input: []sqlc.ChainRedeemSession{
				{
					CustodianGroupUid: []byte{0x1},
					Sequence:          2,
					CurrentPhase:      sqlc.RedeemPhasePREPARING,
					Chain:             "ethereum",
				},
				{
					CustodianGroupUid: []byte{0x1},
					Sequence:          2,
					CurrentPhase:      sqlc.RedeemPhasePREPARING,
					Chain:             "bnb",
				},
			},
			wantErr: false,
			check: func(t *testing.T, repo *healer.HealerRepository, outdated map[string][]sqlc.ChainRedeemSessionUpdate) {
				assert.Empty(t, outdated)
				session, err := repo.GetRedeemSession(context.Background(), mockCustodianGroup.Uid)
				assert.NoError(t, err)
				assert.Empty(t, outdated)
				assert.Equal(t, mockCustodianGroup.Uid, session.CustodianGroupUid)
				assert.Equal(t, int64(2), session.Sequence)
				assert.Equal(t, sqlc.RedeemPhasePREPARING, session.CurrentPhase)
				assert.Equal(t, false, session.IsSwitching.Bool)

				chainSession, err := repo.GetChainRedeemSession(context.Background(), mockCustodianGroup.Uid, "ethereum")
				assert.NoError(t, err)
				assert.Equal(t, mockCustodianGroup.Uid, chainSession.CustodianGroupUid)
				assert.Equal(t, int64(2), chainSession.Sequence)

				chainSession, err = repo.GetChainRedeemSession(context.Background(), mockCustodianGroup.Uid, "bnb")
				assert.NoError(t, err)
				assert.Equal(t, mockCustodianGroup.Uid, chainSession.CustodianGroupUid)
				assert.Equal(t, int64(2), chainSession.Sequence)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			testutils.RunWithTestDB(func(ctx context.Context, repo db.HealderAdapter) error {
				pg := repo.(*healer.HealerRepository)
				// Insert the mock custodian group
				err := pg.SaveCustodianGroups(ctx, []sqlc.CustodianGroup{mockCustodianGroup})
				assert.NoError(t, err)

				outdated, err := pg.SaveRedeemSessionAndChainRedeemSessionsTx(ctx, tc.input)
				if tc.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					tc.check(t, pg, outdated)
				}
				return nil
			})
		})
	}
}
