package daemon

import (
	"context"
	"encoding/hex"
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/scalarorg/data-models/chains"
	"github.com/scalarorg/scalar-healer/pkg/db/sqlc"
	"github.com/scalarorg/scalar-healer/pkg/evm"
	contracts "github.com/scalarorg/scalar-healer/pkg/evm/contracts/generated"
	"github.com/scalarorg/scalar-healer/pkg/utils"
)

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

	wg := sync.WaitGroup{}
	recoverSessions := NewCustodiansRecoverRedeemSessions()
	for _, client := range s.EvmClients {
		wg.Add(1)
		go func() {
			defer wg.Done()

			chainRedeemSessions, err := s.recoverRedeemSession(ctx, client.EvmConfig.GetId(), groupIds)
			if err != nil {
				log.Error().Err(err).Msgf("[Service][RecoverEvmSessions] recover session error: %s", err)
				panic(fmt.Sprintf("[Service][RecoverEvmSessions] cannot recover sessions for evm client %s", client.EvmConfig.GetId()))
			}
			log.Info().
				Str("chainId", client.EvmConfig.GetId()).
				Msg("[Service][RecoverEvmSessions] add evm session")

			recoverSessions.AddRecoverSessions(client.EvmConfig.GetId(), chainRedeemSessions)

		}()
	}
	wg.Wait()
	log.Info().Msgf("[Service][RecoverEvmSessions] finished get SwitchPhase And redeemTx from evm chains")

	recoverSessions.ConstructSessions()

	redeemSessions := make([]sqlc.RedeemSession, 0)
	for _, groupRedeemSessions := range recoverSessions {
		for chain, event := range groupRedeemSessions.SwitchPhaseEvents {
			if len(event) == 0 {
				continue
			}
			for _, event := range event {
				redeemSessions = append(redeemSessions, sqlc.RedeemSession{
					CustodianGroupUid: event.CustodianGroupId[:],
					Sequence:          int64(event.Sequence),
					Chain:             chain,
					CurrentPhase:      sqlc.PhaseFromUint8(event.To),
				})
			}
		}
	}

	log.Info().
		Interface("recoverSessions", recoverSessions).
		Msg("[Service][RecoverEvmSessions] recover sessions")

	err = s.CombinedAdapter.SaveRedeemSessions(ctx, redeemSessions)
	if err != nil {
		log.Error().Err(err).Msgf("[Service][RecoverEvmSessions] cannot save redeem sessions")
		panic(fmt.Sprintf("[Service][RecoverEvmSessions] cannot save redeem sessions"))
	}

	log.Info().Msgf("[Service][RecoverEvmSessions] finished RecoverEvmSessions")
}

func (s *Service) recoverRedeemSession(ctx context.Context, chainId string, groups []string) (*evm.ChainRedeemSessions, error) {
	redeemSessions := evm.NewChainRedeemSessions()

	batch, err := s.CombinedAdapter.GetBatchNumberOfLatestSwitchedPhaseEvents(ctx, 2, chainId, groups)
	if err != nil {
		return nil, fmt.Errorf("failed to get batch number of latest switched phase events: %+v", err)
	}

	if len(batch) == 0 {
		return nil, fmt.Errorf("batch switched phase events is empty: %+v", err)
	}

	for _, events := range batch {
		if err := s.processEventsInBatch(events, redeemSessions); err != nil {
			return nil, err
		}
	}

	return redeemSessions, nil
}

func (s *Service) processEventsInBatch(events []chains.SwitchedPhase, redeemSessions *evm.ChainRedeemSessions) error {
	switch len(events) {
	case 2:
		return s.processTwoEvents(events, redeemSessions)
	case 1:
		return s.processFirstEvent(events[0], redeemSessions)
	default:
		return fmt.Errorf("invalid number of events in batch: %d", len(events))
	}
}

func (s *Service) processTwoEvents(events []chains.SwitchedPhase, redeemSessions *evm.ChainRedeemSessions) error {
	for _, event := range events {
		s.appendSwitchPhaseEvent(event, redeemSessions)
	}
	return nil
}

func (s *Service) processFirstEvent(event chains.SwitchedPhase, redeemSessions *evm.ChainRedeemSessions) error {
	if event.To != sqlc.RedeemPhasePREPARING.Uint8() || event.SessionSequence != 1 {
		return fmt.Errorf("invalid first event: phase=%d, sequence=%d", event.To, event.SessionSequence)
	}
	s.appendSwitchPhaseEvent(event, redeemSessions)
	return nil
}

func (s *Service) appendSwitchPhaseEvent(event chains.SwitchedPhase, redeemSessions *evm.ChainRedeemSessions) {
	redeemSessions.AppendSwitchPhaseEvent(event.CustodianGroupUid, &contracts.IScalarGatewaySwitchPhase{
		CustodianGroupId: common.HexToHash(event.CustodianGroupUid),
		Sequence:         event.SessionSequence,
		From:             event.From,
		To:               event.To,
	})
}

// func (s *Service) processRecoverExecutingPhase(ctx context.Context, groupUid string, groupRedeemSessions *GroupRedeemSessions) error {
// 	log.Info().Str("groupUid", groupUid).
// 		Msg("[Service][RecoverEvmSessions] processRecoverExecutingPhase")
// 	//0. Check if the redeem session is broadcasted to bitcoin network
// 	isBroadcasted, err := s.isRedeemSessionBroadcasted(groupRedeemSessions.RedeemTokenEvents)
// 	if err != nil {
// 		log.Warn().Err(err).Msgf("[Service][processRecoverExecutingPhase] cannot check if the redeem session is broadcasted to bitcoin network")
// 		return err
// 	}
// 	if !isBroadcasted {
// 		log.Info().Msgf("[Service][processRecoverExecutingPhase] redeem session is not broadcasted to bitcoin network")

// 		//1. Replay all switch to preparing phase event,
// 		expectedPhase, evmCounter, hasDifferentPhase := s.replaySwitchPhaseEvents(groupRedeemSessions.SwitchPhaseEvents, 0)
// 		log.Info().Int32("evmCounter", evmCounter).
// 			Any("ExpectedPhase", expectedPhase).
// 			Bool("hasDifferentPhase", hasDifferentPhase).
// 			Msg("[Service][processRecoverExecutingPhase] first events")
// 		if hasDifferentPhase {
// 			panic("[Service][processRecoverExecutingPhase] cannot recover all evm switch phase events to the same phase")
// 		}
// 		if evmCounter != int32(len(s.EvmClients)) {
// 			panic(fmt.Sprintf("[Service][processRecoverExecutingPhase] cannot recover all evm switch phase events, evm counter is %d", evmCounter))
// 		}
// 		if expectedPhase != int32(db.Preparing) {
// 			panic("[Relayer] [processRecoverExecutingPhase] by design, recover first event switch to Preparing for all evm chains")
// 		}

// 		mapTxHashes, err := s.replayRedeemTransactions(ctx, groupUid, groupRedeemSessions.RedeemTokenEvents)
// 		if err != nil {
// 			log.Warn().Err(err).Msgf("[Service][processRecoverExecutionPhase] cannot replay redeem transactions")
// 			return err
// 		}
// 		log.Info().Any("mapTxHashes", mapTxHashes).Msg("[Service][processRecoverExecutionPhase] finished replay redeem transactions")
// 	}
// 	//5. Replay all switch to executing phase events
// 	expectedPhase, evmCounter, hasDifferentPhase := s.replaySwitchPhaseEvents(groupRedeemSessions.SwitchPhaseEvents, 1)
// 	log.Info().Int32("evmCounter", evmCounter).
// 		Any("ExpectedPhase", expectedPhase).
// 		Bool("hasDifferentPhase", hasDifferentPhase).
// 		Msg("[Service][processRecoverExecutionPhase] second events")
// 	if hasDifferentPhase {
// 		panic("[Service][processRecoverExecutionPhase] cannot recover all evm switch phase events")
// 	}
// 	if expectedPhase != int32(db.Executing) {
// 		panic(fmt.Sprintf("[Relayer] [processRecoverExecutionPhase] cannot recover all evm switch phase events, expected phase is %d", expectedPhase))
// 	}

// 	if evmCounter == int32(len(s.EvmClients)) {
// 		log.Info().Int32("evmCounter", evmCounter).Msg("[Service][processRecoverExecutionPhase] all evm chains are in executing phase")
// 		err = s.replayBtcRedeemTxs(groupUid)
// 		if err != nil {
// 			log.Warn().Err(err).Msgf("[Service][processRecoverExecutionPhase] cannot replay btc redeem transactions")
// 			return err
// 		}
// 	} else {
// 		log.Warn().Int32("evmCounter", evmCounter).Msg("[Service][processRecoverExecutionPhase] not all evm chains are in executing phase")
// 	}
// 	return nil
// }
// func (s *Service) processRecoverPreparingPhase(ctx context.Context, groupUid string, groupRedeemSessions *GroupRedeemSessions) error {
// 	log.Info().Str("groupUid", groupUid).
// 		Msg("[Service][RecoverEvmSessions] processRecoverPreparingPhase")
// 	//1. For each evm chain, replay last switch event. It can be Preparing or executing from previous session

// 	// TODO: Update the current session into Healer DB

// 	if expectedPhase == int32(db.Preparing) {
// 		//3. Replay all redeem transactions
// 		mapTxHashes, err := s.replayRedeemTransactions(ctx, groupUid, groupRedeemSessions.RedeemTokenEvents)
// 		if err != nil {
// 			log.Warn().Err(err).Msgf("[Service][processRecoverPreparingPhase] cannot replay redeem transactions")
// 			return err
// 		}
// 		log.Info().Any("mapTxHashes", mapTxHashes).Msg("[Relayer] [processRecoverPreparingPhase] finished replay redeem transactions")
// 	} else if expectedPhase == int32(db.Executing) {
// 		err := s.replayBtcRedeemTxs(groupUid)
// 		if err != nil {
// 			log.Warn().Err(err).Msgf("[Service][processRecoverPreparingPhase] cannot replay btc redeem transactions")
// 			return err
// 		}
// 	}

// 	return nil
// }

/*
 * check if the current redeem session is broadcasted to bitcoin network by checking if the first input utxo is present in the bitcoin network
 */
func (s *Service) isRedeemSessionBroadcasted(mapRedeemTokenEvents map[string][]*contracts.IScalarGatewayRedeemToken) (bool, error) {
	log.Info().Msgf("[Service][isRedeemSessionBroadcasted] checking if the current redeem session is broadcasted to bitcoin network")
	if s.BtcClient == nil {
		return false, fmt.Errorf("[Service][isRedeemSessionBroadcasted] btc client is undefined")
	}
	var firstRedeemTokenEvent *contracts.IScalarGatewayRedeemToken
	for _, redeemTxs := range mapRedeemTokenEvents {
		if len(redeemTxs) > 0 {
			firstRedeemTokenEvent = redeemTxs[0]
			break
		}
	}
	params := evm.RedeemTokenPayloadWithType{}
	err := params.AbiUnpack(firstRedeemTokenEvent.Payload)
	if err != nil {
		return false, fmt.Errorf("[Service][isRedeemSessionBroadcasted] cannot unpack redeem token payload: %s", err)
	}
	if len(params.Utxos) == 0 {
		return false, fmt.Errorf("[Service][isRedeemSessionBroadcasted] no utxos found in redeem token payload")
	}

	if s.BtcClient.Config().ID == firstRedeemTokenEvent.DestinationChain {
		txId := hex.EncodeToString(params.Utxos[0].TxID[:])
		outResult, err := s.BtcClient.GetTxOut(txId, params.Utxos[0].Vout)
		if err != nil {
			return false, fmt.Errorf("[Service][isRedeemSessionBroadcasted] cannot get utxo for redeem token event: %s", err)
		}
		if outResult == nil {
			return true, nil
		}
	} else {
		return false, nil
	}

	return false, nil
}

// //
// func (s *Service) replaySwitchPhaseEvents(mapSwitchPhaseEvents map[string][]*contracts.IScalarGatewaySwitchPhase, index int) (int32, int32, bool) {
// 	wg := sync.WaitGroup{}
// 	var hasDifferentPhase atomic.Bool
// 	var expectedPhase atomic.Int32
// 	var evmCounter atomic.Int32
// 	expectedPhase.Store(-1)

// 	for _, evmClient := range s.EvmClients {
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			chainId := evmClient.EvmConfig.GetId()
// 			switchPhaseEvents, ok := mapSwitchPhaseEvents[chainId]
// 			if !ok || len(switchPhaseEvents) == 0 {
// 				log.Warn().Msgf("[Service][processRecoverPreparingPhase] cannot find redeem session for evm client %s", chainId)
// 				return
// 			}
// 			if index >= len(switchPhaseEvents) {
// 				log.Warn().Str("chainId", chainId).
// 					Int("index", index).
// 					Msgf("[Service][processRecoverPreparingPhase] Switchphase event not found")
// 				return
// 			}
// 			switchPhaseEvent := switchPhaseEvents[index]
// 			expectedPhaseValue := expectedPhase.Load()
// 			if expectedPhaseValue == -1 {
// 				expectedPhase.Store(int32(switchPhaseEvent.To))
// 			} else if expectedPhaseValue != int32(switchPhaseEvent.To) {
// 				log.Warn().Msgf("[Service][processRecoverPreparingPhase] found switch phase event with different phase")
// 				hasDifferentPhase.Store(true)
// 				return
// 			}
// 			err := evmClient.HandleSwitchPhase(ctx, switchPhaseEvent)
// 			if err != nil {
// 				log.Warn().Err(err).Msgf("[Service][processRecoverPreparingPhase] cannot handle switch phase event for evm client %s", chainId)
// 			} else {
// 				evmCounter.Add(1)
// 			}
// 		}()
// 	}
// 	wg.Wait()
// 	return expectedPhase.Load(), evmCounter.Load(), hasDifferentPhase.Load()
// }

func (s *Service) replayRedeemTransactions(ctx context.Context, groupUid string, mapRedeemTokenEvents map[string][]*contracts.IScalarGatewayRedeemToken) (map[string][]string, error) {
	mapTxHashes := sync.Map{}
	wg := sync.WaitGroup{}
	for _, evmClient := range s.EvmClients {
		wg.Add(1)
		go func() {
			defer wg.Done()
			chainId := evmClient.EvmConfig.GetId()
			redeemTokenEvents, ok := mapRedeemTokenEvents[chainId]
			if !ok {
				log.Warn().Str("ChainId", chainId).Msgf("[Service][processRecoverExecutionPhase] no redeemToken event for repaylaying")
				return
			}
			// Scalar network will utxoSnapshot request on each confirm RedeemToken event
			for _, redeemTokenEvent := range redeemTokenEvents {
				err := evmClient.HandleRedeemToken(ctx, redeemTokenEvent)
				if err != nil {
					log.Warn().
						Str("chainId", chainId).
						Any("redeemTokenEvent", redeemTokenEvent).
						Err(err).Msgf("[Service][processRecoverExecutionPhase] cannot handle redeem token event")
				} else {
					value, loaded := mapTxHashes.LoadOrStore(redeemTokenEvent.DestinationChain, []string{redeemTokenEvent.Raw.TxHash.Hex()})
					if loaded {
						mapTxHashes.Store(redeemTokenEvent.DestinationChain, append(value.([]string), redeemTokenEvent.Raw.TxHash.Hex()))
					}
				}
			}
			log.Info().Str("ChainId", chainId).Int("RedeemTx count", len(redeemTokenEvents)).
				Msgf("[Service][processRecoverExecutionPhase] finished handle redeem token events")
		}()
	}
	wg.Wait()
	result := map[string][]string{}
	mapTxHashes.Range(func(key, value interface{}) bool {
		result[key.(string)] = value.([]string)
		return true
	})
	return result, nil
}

func (s *Service) replayBtcRedeemTxs(groupUid string) error {
	log.Info().Str("groupUid", groupUid).Msgf("[Service][processRecoverPreparingPhase] replay btc redeem transactions")
	groupBytes, err := hex.DecodeString(groupUid)
	if err != nil {
		return fmt.Errorf("[Service][processRecoverPreparingPhase] cannot decode group uid: %s", err)
	}

	_ = groupBytes

	// TODO: Get redeem session from DB

	// redeemSession, err := s.ScalarClient.GetRedeemSession(groupBytes)
	// if err != nil {
	// 	return fmt.Errorf("[Service][processRecoverPreparingPhase] cannot get redeem session for group %s", groupUid)
	// }
	// redeemTxs := s.ScalarClient.PickCacheRedeemTx(groupUid, redeemSession.Session.Sequence)
	// log.Info().Any("redeemTxs", redeemTxs).Msgf("[Service][replayBtcRedeemTxs] redeem txs in cache")
	// for chainId, redeemTxs := range redeemTxs {
	// 	err := s.ScalarClient.BroadcastRedeemTxsConfirmRequest(chainId, groupUid, redeemTxs)
	// 	if err != nil {
	// 		return fmt.Errorf("[Service][processRecoverPreparingPhase] cannot broadcast redeem txs confirm request for group %s", groupUid)
	// 	}
	// }
	return nil
}
