package job

import (
	"context"

	"github.com/scalarorg/scalar-healer/pkg/db"
	"github.com/scalarorg/scalar-healer/pkg/db/healer"
	"github.com/scalarorg/scalar-healer/pkg/db/sqlc"
	"github.com/scalarorg/scalar-healer/pkg/tofnd"
)

func HandlePendingSigning(repo db.HealderAdapter, tofndManager *tofnd.Manager) func() {
	return func() {
		// ctx := context.Background()
		// commands, err := repo.ListPendingSigningRedeemCommands(ctx)
		// if err != nil {
		// 	// log error
		// 	return
		// }
		// for _, req := range requests {
		// 	go handleSigning(ctx, repo, tofndClients, req)
		// }
	}
}

func handlePendingSigning(repo *healer.HealerRepository, tofndClients []*tofnd.Client) {
	// ctx := context.Background()
	// requests, err := repo.ListPendingSigningRequests(ctx)
	// if err != nil {
	// 	// log error
	// 	return
	// }
	// for _, req := range requests {
	// 	go handleSigning(ctx, repo, tofndClients, req)
	// }
}

func handleSigning(ctx context.Context, repo *healer.HealerRepository, tofndClients []*tofnd.Client, req sqlc.RedeemCommand) {
	// msg := prepareSignMessage(req) // implement as needed
	// var sigs [][]byte
	// for _, client := range tofndClients {
	// 	sig, err := client.Sign(ctx, msg)
	// 	if err != nil {
	// 		// log error, optionally retry
	// 		continue
	// 	}
	// 	sigs = append(sigs, sig)
	// }
	// if len(sigs) >= req.Quorum {
	// 	repo.UpdateRedeemRequestStatus(req.ID, "signed", sigs)
	// } else {
	// 	repo.UpdateRedeemRequestStatus(req.ID, "signing_failed", nil)
	// }
}
