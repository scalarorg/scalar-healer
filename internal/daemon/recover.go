package daemon

import (
	"context"
	"encoding/hex"
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/scalarorg/scalar-healer/pkg/db/sqlc"
	"github.com/scalarorg/scalar-healer/pkg/evm"
	"github.com/scalarorg/scalar-healer/pkg/utils"
)

type redeemSessionResult struct {
	chainId  string
	sessions []sqlc.ChainRedeemSession
	err      error
}

func (s *Service) RecoverEvmSessions(ctx context.Context) {
	groups, err := s.CombinedAdapter.GetAllCustodianGroups(ctx)
	if err != nil {
		log.Error().Err(err).Msg("[DaemonService] Cannot get custodian groups")
		panic(err)
	}

	log.Info().
		Int("len(groups)", len(groups)).
		Interface("groups", groups).
		Msg("[Service][RecoverEvmSessions] start recover evm sessions")

	groupIds := utils.Map(groups, func(group sqlc.CustodianGroup) string {
		return hex.EncodeToString(group.Uid)
	})

	resultChan := make(chan *redeemSessionResult, len(s.EvmClients))
	wg := sync.WaitGroup{}

	for _, client := range s.EvmClients {
		wg.Add(1)

		go func(c *evm.EvmClient) {
			defer wg.Done()
			chainId := c.EvmConfig.GetId()

			// map[custodianGroupUid] = SwitchedPhaseEvent // 1 element
			events, aErr := s.CombinedAdapter.GetBatchLastestSwitchedPhaseEvents(ctx, chainId, groupIds)
			if aErr != nil {
				resultChan <- &redeemSessionResult{
					chainId:  chainId,
					sessions: nil,
					err:      aErr,
				}
				return
			}

			sessions := make([]sqlc.ChainRedeemSession, 0)

			for _, e := range events {
				sessions = append(sessions, sqlc.ChainRedeemSession{
					Chain:             chainId,
					CustodianGroupUid: common.HexToHash(e.CustodianGroupUid).Bytes(),
					Sequence:          int64(e.SessionSequence),
					CurrentPhase:      sqlc.PhaseFromUint8(e.To),
				})
			}

			resultChan <- &redeemSessionResult{
				chainId:  chainId,
				sessions: sessions,
				err:      aErr,
			}
		}(client)
	}

	go func() {
		wg.Wait()
		close(resultChan)
		log.Info().Msgf("[Service][RecoverEvmSessions] finished get SwitchPhase And redeemTx from evm chains")
	}()

	sessions := make([]sqlc.ChainRedeemSession, 0)
	for result := range resultChan {
		if result.err != nil {
			log.Error().Err(result.err).
				Str("chainId", result.chainId).
				Msgf("[Service][RecoverEvmSessions] recover session error")
			panic(fmt.Sprintf("[Service][RecoverEvmSessions] cannot recover sessions for evm client %s", result.chainId))
		}

		sessions = append(sessions, result.sessions...)

		log.Info().
			Str("chainId", result.chainId).
			Msg("[Service][RecoverEvmSessions] add evm session")
	}

	log.Info().
		Int("len(sessions)", len(sessions)).
		Interface("sessions", sessions).
		Msg("[Service][RecoverEvmSessions] start recover evm sessions")

	outdatedSessionsByGroup, err := s.CombinedAdapter.SaveRedeemSessionAndChainRedeemSessionsTx(ctx, sessions)
	if err != nil {
		log.Error().Err(err).Msgf("[Service][RecoverEvmSessions] cannot save redeem sessions")
		panic(err)
	}

	log.Info().Msgf("[Service][RecoverEvmSessions] finished RecoverEvmSessions")

	if len(outdatedSessionsByGroup) == 0 {
		return
	}

	switchPhaseCmds := make([]sqlc.Command, 0)

	for groupUid, sessions := range outdatedSessionsByGroup {
		grUidbz := common.HexToHash(groupUid).Bytes()
		for _, session := range sessions {
			id := make([]byte, 0)
			id = append(id, grUidbz...)
			id = append(id, byte(session.Sequence))
			id = append(id, session.CurrentPhase.Bytes())

			switchPhaseCmds = append(switchPhaseCmds, sqlc.Command{
				CommandID:   sqlc.NewCommandID(id, session.Chain).Bytes(),
				CommandType: sqlc.CommandTypeSwitchPhase,
				Params:      evm.CreateSwitchPhaseParams(grUidbz, session.NewPhase),
				Chain:       session.Chain,
				Status:      sqlc.COMMAND_STATUS_PENDING.ToPgType(),
				Payload:     nil,
			})
		}
	}

	err = s.CombinedAdapter.SaveCommandsAndBatchCommandsTx(ctx, switchPhaseCmds)
	if err != nil {
		log.Error().Err(err).Msgf("[Service][RecoverEvmSessions] cannot save switch phase commands")
		panic(err)
	}

	// TODO: Simulate switch phase

	return
}
