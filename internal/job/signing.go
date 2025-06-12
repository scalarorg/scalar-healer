package job

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/scalarorg/scalar-healer/pkg/db"
	"github.com/scalarorg/scalar-healer/pkg/db/sqlc"
	"github.com/scalarorg/scalar-healer/pkg/tofnd"
	"github.com/scalarorg/scalar-healer/pkg/utils/slices"
)

type Result struct {
	CmdID []byte
	Musig []tofnd.SignatureResult
}

func HandlePendingSigning(repo db.HealderAdapter, tofndManager *tofnd.Manager) func() {
	return func() {

		log.Info().Msg("Starting signing redeem commands")

		ctx := context.Background()
		commands, err := repo.ListPendingSigningRedeemCommands(ctx)
		if err != nil {
			log.Error().Err(err).Msg("Failed to list pending signing redeem commands")
			return
		}

		log.Info().Msgf("Found %d pending signing redeem commands", len(commands))

		wg := sync.WaitGroup{}

		results := make(chan Result, len(commands))

		for _, cmd := range commands {
			wg.Add(1)
			go func(cmd sqlc.RedeemCommand) {
				defer wg.Done()
				musig, err := tofndManager.Sign(ctx, cmd.SigHash)
				if err != nil {
					log.Error().Err(err).Msg("Failed to sign redeem command")
					return
				}
				results <- Result{CmdID: cmd.ID, Musig: musig}
			}(cmd)
		}

		wg.Wait()
		close(results)

		signatures := make([][]byte, 0, len(commands))
		ids := make([][]byte, 0, len(commands))
		for result := range results {
			formatedSigs := slices.Map(result.Musig, func(sig tofnd.SignatureResult) map[string]interface{} {
				m := make(map[string]interface{})
				m["signature"] = sig.Sig
				m["party_id"] = sig.Client.PartyID
				return m
			})

			json, err := json.Marshal(formatedSigs)
			if err != nil {
				log.Error().Err(err).Msg("Failed to marshal signature")
				return
			}
			signatures = append(signatures, json)
			ids = append(ids, result.CmdID)
		}

		err = repo.SubmitRedeemCommandSignatures(ctx, ids, signatures)
		if err != nil {
			log.Error().Err(err).Msg("Failed to submit redeem command signatures")
			return
		}
		log.Info().Msg("Submitted redeem command signatures")
	}
}
