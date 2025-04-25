package daemon

import (
	"encoding/hex"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/rs/zerolog/log"
	"github.com/scalarorg/scalar-healer/pkg/db/models"
	"github.com/scalarorg/scalar-healer/pkg/evm"
	contracts "github.com/scalarorg/scalar-healer/pkg/evm/contracts/generated"
)

func (s *Service) RecoverEvmSessions(groups []string) error {
	wg := sync.WaitGroup{}
	recoverSessions := CustodiansRecoverRedeemSessions{}
	for _, client := range s.EvmClients {
		wg.Add(1)
		go func() {
			defer wg.Done()
			chainRedeemSessions, err := client.RecoverRedeemSessions(groups)
			if err != nil {
				log.Warn().Err(err).Msgf("[Relayer] [Start] cannot recover sessions for evm client %s", client.EvmConfig.GetId())
			}
			if chainRedeemSessions != nil {
				log.Info().Str("chainId", client.EvmConfig.GetId()).
					//Any("chainRedeemSessions", chainRedeemSessions).
					Msg("[Relayer] [RecoverEvmSessions] add evm session")
				recoverSessions.AddRecoverSessions(client.EvmConfig.GetId(), chainRedeemSessions)
			} else {
				panic(fmt.Sprintf("[Relayer] [RecoverEvmSessions] cannot recover sessions for evm client %s", client.EvmConfig.GetId()))
			}
		}()
	}
	wg.Wait()
	log.Info().Msgf("[Relayer] [RecoverEvmSessions] finished get SwitchPhase And redeemTx from evm chains")
	recoverSessions.ConstructSessions()
	for groupUid, groupRedeemSessions := range recoverSessions.RecoverSessions {
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Info().Str("groupUid", groupUid).
				Any("maxSession", groupRedeemSessions.MaxSession).
				Any("minSession", groupRedeemSessions.MinSession).
				//Any("switchPhaseEvents", groupRedeemSessions.SwitchPhaseEvents).
				//Any("redeemTokenEvents", groupRedeemSessions.RedeemTokenEvents).
				Msg("[Relayer] [RecoverEvmSessions] recovered redeem session for each group")
			if groupRedeemSessions.MaxSession.Phase == models.Executing {
				err := s.processRecoverExecutingPhase(groupUid, groupRedeemSessions)
				if err != nil {
					log.Warn().Err(err).Msgf("[Relayer] [RecoverEvmSessions] cannot process recover executing phase for group %s", groupUid)
				}
			} else if groupRedeemSessions.MaxSession.Phase == models.Preparing {
				err := s.processRecoverPreparingPhase(groupUid, groupRedeemSessions)
				if err != nil {
					log.Warn().Err(err).Msgf("[Relayer] [RecoverEvmSessions] cannot process recover preparing phase for group %s", groupUid)
				}
			}
		}()
	}
	wg.Wait()
	log.Info().Msgf("[Relayer] [RecoverEvmSessions] finished RecoverEvmSessions")
	return nil
}
func (s *Service) processRecoverExecutingPhase(groupUid string, groupRedeemSessions *GroupRedeemSessions) error {
	log.Info().Str("groupUid", groupUid).
		Msg("[Relayer] [RecoverEvmSessions] processRecoverExecutingPhase")
	//0. Check if the redeem session is broadcasted to bitcoin network
	isBroadcasted, err := s.isRedeemSessionBroadcasted(groupRedeemSessions.RedeemTokenEvents)
	if err != nil {
		log.Warn().Err(err).Msgf("[Relayer] [processRecoverExecutingPhase] cannot check if the redeem session is broadcasted to bitcoin network")
		return err
	}
	if !isBroadcasted {
		log.Info().Msgf("[Relayer] [processRecoverExecutingPhase] redeem session is not broadcasted to bitcoin network")

		//1. Replay all switch to preparing phase event,
		expectedPhase, evmCounter, hasDifferentPhase := s.replaySwitchPhaseEvents(groupRedeemSessions.SwitchPhaseEvents, 0)
		log.Info().Int32("evmCounter", evmCounter).
			Any("ExpectedPhase", expectedPhase).
			Bool("hasDifferentPhase", hasDifferentPhase).
			Msg("[Relayer] [processRecoverExecutingPhase] first events")
		if hasDifferentPhase {
			panic("[Relayer] [processRecoverExecutingPhase] cannot recover all evm switch phase events to the same phase")
		}
		if evmCounter != int32(len(s.EvmClients)) {
			panic(fmt.Sprintf("[Relayer] [processRecoverExecutingPhase] cannot recover all evm switch phase events, evm counter is %d", evmCounter))
		}
		if expectedPhase != int32(models.Preparing) {
			panic("[Relayer] [processRecoverExecutingPhase] by design, recover first event switch to Preparing for all evm chains")
		}

		mapTxHashes, err := s.replayRedeemTransactions(groupUid, groupRedeemSessions.RedeemTokenEvents)
		if err != nil {
			log.Warn().Err(err).Msgf("[Relayer] [processRecoverExecutionPhase] cannot replay redeem transactions")
			return err
		}
		log.Info().Any("mapTxHashes", mapTxHashes).Msg("[Relayer] [processRecoverExecutionPhase] finished replay redeem transactions")
	}
	//5. Replay all switch to executing phase events
	expectedPhase, evmCounter, hasDifferentPhase := s.replaySwitchPhaseEvents(groupRedeemSessions.SwitchPhaseEvents, 1)
	log.Info().Int32("evmCounter", evmCounter).
		Any("ExpectedPhase", expectedPhase).
		Bool("hasDifferentPhase", hasDifferentPhase).
		Msg("[Relayer] [processRecoverExecutionPhase] second events")
	if hasDifferentPhase {
		panic("[Relayer] [processRecoverExecutionPhase] cannot recover all evm switch phase events")
	}
	if expectedPhase != int32(models.Executing) {
		panic(fmt.Sprintf("[Relayer] [processRecoverExecutionPhase] cannot recover all evm switch phase events, expected phase is %d", expectedPhase))
	}

	if evmCounter == int32(len(s.EvmClients)) {
		log.Info().Int32("evmCounter", evmCounter).Msg("[Relayer] [processRecoverExecutionPhase] all evm chains a in executing phase")
		err = s.replayBtcRedeemTxs(groupUid)
		if err != nil {
			log.Warn().Err(err).Msgf("[Relayer] [processRecoverExecutionPhase] cannot replay btc redeem transactions")
			return err
		}
	} else {
		log.Warn().Int32("evmCounter", evmCounter).Msg("[Relayer] [processRecoverExecutionPhase] not all evm chains are in executing phase")
	}
	return nil
}
func (s *Service) processRecoverPreparingPhase(groupUid string, groupRedeemSessions *GroupRedeemSessions) error {
	log.Info().Str("groupUid", groupUid).
		Msg("[Relayer] [RecoverEvmSessions] processRecoverPreparingPhase")
	//1. For each evm chain, replay last switch event. It can be Preparing or executing from previous session
	expectedPhase, evmCounter, hasDifferentPhase := s.replaySwitchPhaseEvents(groupRedeemSessions.SwitchPhaseEvents, 0)
	if hasDifferentPhase {
		panic("[Relayer] [processRecoverPreparingPhase] cannot recover all evm switch phase events")
	}
	if evmCounter != int32(len(s.EvmClients)) {
		panic(fmt.Sprintf("[Relayer] [processRecoverPreparingPhase] cannot recover all evm switch phase events, evm counter is %d", evmCounter))
	}

	if expectedPhase == int32(models.Preparing) {
		//3. Replay all redeem transactions
		mapTxHashes, err := s.replayRedeemTransactions(groupUid, groupRedeemSessions.RedeemTokenEvents)
		if err != nil {
			log.Warn().Err(err).Msgf("[Relayer] [processRecoverPreparingPhase] cannot replay redeem transactions")
			return err
		}
		log.Info().Any("mapTxHashes", mapTxHashes).Msg("[Relayer] [processRecoverPreparingPhase] finished replay redeem transactions")
	} else if expectedPhase == int32(models.Executing) {
		err := s.replayBtcRedeemTxs(groupUid)
		if err != nil {
			log.Warn().Err(err).Msgf("[Relayer] [processRecoverPreparingPhase] cannot replay btc redeem transactions")
			return err
		}
	}
	return nil
}

/*
 * check if the current redeem session is broadcasted to bitcoin network by checking if the first input utxo is present in the bitcoin network
 */
func (s *Service) isRedeemSessionBroadcasted(mapRedeemTokenEvents map[string][]*contracts.IScalarGatewayRedeemToken) (bool, error) {
	log.Info().Msgf("[Relayer] [isRedeemSessionBroadcasted] checking if the current redeem session is broadcasted to bitcoin network")
	if s.BtcClient == nil {
		return false, fmt.Errorf("[Relayer] [isRedeemSessionBroadcasted] btc client is undefined")
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
		return false, fmt.Errorf("[Relayer] [isRedeemSessionBroadcasted] cannot unpack redeem token payload: %s", err)
	}
	if len(params.Utxos) == 0 {
		return false, fmt.Errorf("[Relayer] [isRedeemSessionBroadcasted] no utxos found in redeem token payload")
	}

	if s.BtcClient.Config().ID == firstRedeemTokenEvent.DestinationChain {
		txId := hex.EncodeToString(params.Utxos[0].TxID[:])
		outResult, err := s.BtcClient.GetTxOut(txId, params.Utxos[0].Vout)
		if err != nil {
			return false, fmt.Errorf("[Relayer] [isRedeemSessionBroadcasted] cannot get utxo for redeem token event: %s", err)
		}
		if outResult == nil {
			return true, nil
		}
	} else {
		return false, nil
	}

	return false, nil
}

func (s *Service) replaySwitchPhaseEvents(mapSwitchPhaseEvents map[string][]*contracts.IScalarGatewaySwitchPhase, index int) (int32, int32, bool) {
	wg := sync.WaitGroup{}
	var hasDifferentPhase atomic.Bool
	var expectedPhase atomic.Int32
	var evmCounter atomic.Int32
	expectedPhase.Store(-1)

	for _, evmClient := range s.EvmClients {
		wg.Add(1)
		go func() {
			defer wg.Done()
			chainId := evmClient.EvmConfig.GetId()
			switchPhaseEvents, ok := mapSwitchPhaseEvents[chainId]
			if !ok || len(switchPhaseEvents) == 0 {
				log.Warn().Msgf("[Relayer] [processRecoverPreparingPhase] cannot find redeem session for evm client %s", chainId)
				return
			}
			if index >= len(switchPhaseEvents) {
				log.Warn().Str("chainId", chainId).
					Int("index", index).
					Msgf("[Relayer] [processRecoverPreparingPhase] Switchphase event not found")
				return
			}
			switchPhaseEvent := switchPhaseEvents[index]
			expectedPhaseValue := expectedPhase.Load()
			if expectedPhaseValue == -1 {
				expectedPhase.Store(int32(switchPhaseEvent.To))
			} else if expectedPhaseValue != int32(switchPhaseEvent.To) {
				log.Warn().Msgf("[Relayer] [processRecoverPreparingPhase] found switch phase event with different phase")
				hasDifferentPhase.Store(true)
				return
			}
			err := evmClient.HandleSwitchPhase(switchPhaseEvent)
			if err != nil {
				log.Warn().Err(err).Msgf("[Relayer] [processRecoverPreparingPhase] cannot handle switch phase event for evm client %s", chainId)
			} else {
				evmCounter.Add(1)
			}
		}()
	}
	wg.Wait()
	return expectedPhase.Load(), evmCounter.Load(), hasDifferentPhase.Load()
}

func (s *Service) replayRedeemTransactions(groupUid string, mapRedeemTokenEvents map[string][]*contracts.IScalarGatewayRedeemToken) (map[string][]string, error) {
	mapTxHashes := sync.Map{}
	wg := sync.WaitGroup{}
	for _, evmClient := range s.EvmClients {
		wg.Add(1)
		go func() {
			defer wg.Done()
			chainId := evmClient.EvmConfig.GetId()
			redeemTokenEvents, ok := mapRedeemTokenEvents[chainId]
			if !ok {
				log.Warn().Str("ChainId", chainId).Msgf("[Relayer] [processRecoverExecutionPhase] no redeemToken event for repaylaying")
				return
			}
			// Scalar network will utxoSnapshot request on each confirm RedeemToken event
			for _, redeemTokenEvent := range redeemTokenEvents {
				err := evmClient.HandleRedeemToken(redeemTokenEvent)
				if err != nil {
					log.Warn().
						Str("chainId", chainId).
						Any("redeemTokenEvent", redeemTokenEvent).
						Err(err).Msgf("[Relayer] [processRecoverExecutionPhase] cannot handle redeem token event")
				} else {
					value, loaded := mapTxHashes.LoadOrStore(redeemTokenEvent.DestinationChain, []string{redeemTokenEvent.Raw.TxHash.Hex()})
					if loaded {
						mapTxHashes.Store(redeemTokenEvent.DestinationChain, append(value.([]string), redeemTokenEvent.Raw.TxHash.Hex()))
					}
				}
			}
			log.Info().Str("ChainId", chainId).Int("RedeemTx count", len(redeemTokenEvents)).
				Msgf("[Relayer] [processRecoverExecutionPhase] finished handle redeem token events")
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
	log.Info().Str("groupUid", groupUid).Msgf("[Relayer] [processRecoverPreparingPhase] replay btc redeem transactions")
	// groupBytes, err := hex.DecodeString(groupUid)
	// if err != nil {
	// 	return fmt.Errorf("[Relayer] [processRecoverPreparingPhase] cannot decode group uid: %s", err)
	// }
	// redeemSession, err := s.ScalarClient.GetRedeemSession(groupBytes)
	// if err != nil {
	// 	return fmt.Errorf("[Relayer] [processRecoverPreparingPhase] cannot get redeem session for group %s", groupUid)
	// }
	// redeemTxs := s.ScalarClient.PickCacheRedeemTx(groupUid, redeemSession.Session.Sequence)
	// log.Info().Any("redeemTxs", redeemTxs).Msgf("[Relayer] [replayBtcRedeemTxs] redeem txs in cache")
	// for chainId, redeemTxs := range redeemTxs {
	// 	err := s.ScalarClient.BroadcastRedeemTxsConfirmRequest(chainId, groupUid, redeemTxs)
	// 	if err != nil {
	// 		return fmt.Errorf("[Relayer] [processRecoverPreparingPhase] cannot broadcast redeem txs confirm request for group %s", groupUid)
	// 	}
	// }
	return nil
}
