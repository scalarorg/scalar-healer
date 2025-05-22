package healer

import (
	"context"

	"github.com/scalarorg/scalar-healer/pkg/db/sqlc"
)

func (m *HealerRepository) GetAllCustodianGroups(ctx context.Context) ([]sqlc.CustodianGroup, error) {
	return m.Queries.GetAllCustodianGroups(ctx)
}

func (m *HealerRepository) GetCustodianGroup(ctx context.Context, uid []byte) (sqlc.CustodianGroup, error) {
	return m.Queries.GetCustodianGroupByUID(ctx, uid)
}

func (m *HealerRepository) SaveCustodianGroups(ctx context.Context, grs []sqlc.CustodianGroup) error {
	var uids [][]byte
	var names []string
	var bitcoinPubkeys [][]byte
	var quorums []int64
	for _, gr := range grs {
		uids = append(uids, gr.Uid)
		names = append(names, gr.Name)
		bitcoinPubkeys = append(bitcoinPubkeys, gr.BitcoinPubkey)
		quorums = append(quorums, gr.Quorum)
	}

	return m.Queries.SaveCustodianGroups(ctx, sqlc.SaveCustodianGroupsParams{
		Column1: uids,
		Column2: names,
		Column3: bitcoinPubkeys,
		Column4: quorums,
	})
}
