package daemon

// import (
// 	"math"

// 	"github.com/rs/zerolog/log"
// 	"github.com/scalarorg/scalar-healer/pkg/db/sqlc"
// 	"github.com/scalarorg/scalar-healer/pkg/evm"
// 	contracts "github.com/scalarorg/scalar-healer/pkg/evm/contracts/generated"
// )

// // Eache chainId container an array of SwitchPhaseEvents,
// // with the first element is switch to Preparing phase
// type GroupRedeemSessions struct {
// 	GroupUid          string
// 	MaxSession        sqlc.RedeemSession
// 	MinSession        sqlc.RedeemSession
// 	SwitchPhaseEvents map[string][]*contracts.IScalarGatewaySwitchPhase //Map by chainId
// 	RedeemTokenEvents map[string][]*contracts.IScalarGatewayRedeemToken
// }

// /*
// * For each custodian group, maximum difference between the session of evms is a phase
//  */
// func (s *GroupRedeemSessions) Construct() {
// 	s.MinSession.Sequence = math.MaxInt64
// 	//Find the max, min session
// 	for _, switchPhaseEvent := range s.SwitchPhaseEvents {
// 		lastEvent := switchPhaseEvent[len(switchPhaseEvent)-1]
// 		if s.MaxSession.Sequence < int64(lastEvent.Sequence) {
// 			s.MaxSession.Sequence = int64(lastEvent.Sequence)
// 			s.MaxSession.CurrentPhase = sqlc.PhaseFromUint8(lastEvent.To)
// 		} else if s.MaxSession.Sequence == int64(lastEvent.Sequence) && s.MaxSession.CurrentPhase.Uint8() < lastEvent.To {
// 			s.MaxSession.CurrentPhase = sqlc.PhaseFromUint8(lastEvent.To)
// 		}
// 		if s.MinSession.Sequence > int64(lastEvent.Sequence) {
// 			s.MinSession.Sequence = int64(lastEvent.Sequence)
// 			s.MinSession.CurrentPhase = sqlc.PhaseFromUint8(lastEvent.To)
// 		} else if s.MinSession.Sequence == int64(lastEvent.Sequence) && s.MinSession.CurrentPhase.Uint8() > lastEvent.To {
// 			s.MinSession.CurrentPhase = sqlc.PhaseFromUint8(lastEvent.To)
// 		}
// 	}

// 	if s.MaxSession.CurrentPhase == sqlc.RedeemPhasePREPARING {
// 		s.ConstructPreparingPhase()
// 	} else if s.MaxSession.CurrentPhase == sqlc.RedeemPhaseEXECUTING {
// 		s.ConstructExecutingPhase()
// 	}
// }

// func (s *GroupRedeemSessions) ConstructPreparingPhase() {
// 	diff := s.MaxSession.Cmp(&s.MinSession)

// 	log.Info().Str("groupUid", s.GroupUid).Int64("diff", diff).
// 		Any("maxSession", s.MaxSession).
// 		Any("minSession", s.MinSession).
// 		Msg("[GroupRedeemSessions] [ConstructPreparingPhase]")

// 	if diff == 0 {
// 		log.Warn().Str("groupUid", s.GroupUid).Msg("[GroupRedeemSessions] [ConstructPreparingPhase] max session and min session are the same")
// 		//Each chain keep only one switch phase event to Preparing phase
// 		for chainId, switchPhaseEvent := range s.SwitchPhaseEvents {
// 			if len(switchPhaseEvent) == 0 {
// 				continue
// 			}
// 			if len(switchPhaseEvent) == 2 {
// 				s.SwitchPhaseEvents[chainId] = switchPhaseEvent[1:]
// 			}
// 		}
// 		//Keep all redeem token events of the max session's sequence
// 		for chainId, redeemTokenEvents := range s.RedeemTokenEvents {
// 			currentSessionEvents := make([]*contracts.IScalarGatewayRedeemToken, 0)
// 			for _, redeemTokenEvent := range redeemTokenEvents {
// 				if int64(redeemTokenEvent.Sequence) == s.MaxSession.Sequence {
// 					currentSessionEvents = append(currentSessionEvents, redeemTokenEvent)
// 				}
// 			}
// 			s.RedeemTokenEvents[chainId] = currentSessionEvents
// 		}
// 	} else {
// 		//These are some chains switch to the preparing phase, and some other is in execution phase from previous session
// 		//We don't need to recreate Redeem transaction to btc,
// 		//show we don't need to send RedeemEvent for confirmation
// 		s.RedeemTokenEvents = make(map[string][]*contracts.IScalarGatewayRedeemToken)
// 		//Remove old switch phase event
// 		//Find all chains with 2 events [Preparing, Executing], remove the first event
// 		for chainId, switchPhaseEvent := range s.SwitchPhaseEvents {
// 			if len(switchPhaseEvent) == 2 &&
// 				switchPhaseEvent[0].To == sqlc.RedeemPhasePREPARING.Uint8() &&
// 				switchPhaseEvent[1].To == sqlc.RedeemPhaseEXECUTING.Uint8() {
// 				s.SwitchPhaseEvents[chainId] = switchPhaseEvent[1:]
// 			}
// 		}
// 	}
// }

// func (s *GroupRedeemSessions) ConstructExecutingPhase() {
// 	//For both case diff == 0 and diff = 1, we need to resend the redeem transaction to the scalar network
// 	//Expecting all chains are switching to the executing phase
// 	for chainId, switchPhaseEvent := range s.SwitchPhaseEvents {
// 		if int64(switchPhaseEvent[0].Sequence) < s.MaxSession.Sequence {
// 			log.Warn().Str("chainId", chainId).Any("First preparing event", switchPhaseEvent[0]).
// 				Msg("[Service][RecoverRedeemSessions] Session is too low. Some thing wrong")
// 		}
// 	}
// 	//We resend to the scalar network onlye the redeem transaction of the last session
// 	for chainId, redeemTokenEvent := range s.RedeemTokenEvents {
// 		for _, event := range redeemTokenEvent {
// 			if int64(event.Sequence) < s.MaxSession.Sequence {
// 				log.Warn().Str("chainId", chainId).Any("Redeem transaction", event).
// 					Msg("[Service][RecoverRedeemSessions] Redeem transaction is too low. Some thing wrong")
// 			}
// 		}
// 	}
// }

// // Store all evm recovering redeem sessions
// type custodiansRecoverRedeemSessions map[string]*GroupRedeemSessions

// func NewCustodiansRecoverRedeemSessions() custodiansRecoverRedeemSessions {
// 	return make(custodiansRecoverRedeemSessions)
// }

// func (s custodiansRecoverRedeemSessions) AddRecoverSessions(chainId string, chainRedeemSessions *evm.ChainRedeemSessions) {
// 	for groupUid, switchPhaseEvent := range chainRedeemSessions.SwitchPhaseEvents {
// 		if len(switchPhaseEvent) == 0 {
// 			continue
// 		}
// 		groupSession, ok := s[groupUid]
// 		if !ok {
// 			groupSession = &GroupRedeemSessions{
// 				GroupUid:          groupUid,
// 				SwitchPhaseEvents: make(map[string][]*contracts.IScalarGatewaySwitchPhase),
// 				RedeemTokenEvents: make(map[string][]*contracts.IScalarGatewayRedeemToken),
// 			}
// 		}
// 		groupSession.SwitchPhaseEvents[chainId] = switchPhaseEvent
// 		s[groupUid] = groupSession
// 	}
// 	for groupUid, redeemTokenEvent := range chainRedeemSessions.RedeemTokenEvents {
// 		if len(redeemTokenEvent) == 0 {
// 			continue
// 		}
// 		groupSession, ok := s[groupUid]
// 		if !ok {
// 			log.Warn().Msgf("[Service][AddRecoverSessions] no recover session found for group %s", groupUid)
// 			groupSession = &GroupRedeemSessions{
// 				GroupUid:          groupUid,
// 				SwitchPhaseEvents: make(map[string][]*contracts.IScalarGatewaySwitchPhase),
// 				RedeemTokenEvents: make(map[string][]*contracts.IScalarGatewayRedeemToken),
// 			}
// 		}
// 		groupSession.RedeemTokenEvents[chainId] = redeemTokenEvent
// 		s[groupUid] = groupSession
// 	}
// }

// func (s custodiansRecoverRedeemSessions) ConstructSessions() {
// 	log.Info().Msg("[Service][ConstructSessions] start construct sessions")
// 	for _, groupSession := range s {
// 		groupSession.Construct()
// 	}
// }

// func (s custodiansRecoverRedeemSessions) GroupByChain() map[string]*evm.ChainRedeemSessions {
// 	mapChainRedeemSessions := make(map[string]*evm.ChainRedeemSessions)
// 	for groupUid, groupSession := range s {
// 		for chainId, switchPhaseEvent := range groupSession.SwitchPhaseEvents {
// 			mapChainRedeemSessions[chainId] = &evm.ChainRedeemSessions{
// 				SwitchPhaseEvents: map[string][]*contracts.IScalarGatewaySwitchPhase{groupUid: switchPhaseEvent},
// 				RedeemTokenEvents: map[string][]*contracts.IScalarGatewayRedeemToken{},
// 			}
// 		}
// 		for chainId, redeemTokenEvent := range groupSession.RedeemTokenEvents {
// 			chainRedeemSessions, ok := mapChainRedeemSessions[chainId]
// 			if !ok {
// 				mapChainRedeemSessions[chainId] = &evm.ChainRedeemSessions{
// 					SwitchPhaseEvents: map[string][]*contracts.IScalarGatewaySwitchPhase{},
// 					RedeemTokenEvents: map[string][]*contracts.IScalarGatewayRedeemToken{groupUid: redeemTokenEvent},
// 				}
// 			} else {
// 				chainRedeemSessions.RedeemTokenEvents[groupUid] = redeemTokenEvent
// 			}
// 		}
// 	}
// 	return mapChainRedeemSessions
// }
