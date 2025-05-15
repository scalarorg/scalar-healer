package evm

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	"github.com/scalarorg/scalar-healer/pkg/db"
	contracts "github.com/scalarorg/scalar-healer/pkg/evm/contracts/generated"
)

var (
	ALL_EVENTS = []string{
		EVENT_EVM_CONTRACT_CALL,
		EVENT_EVM_CONTRACT_CALL_WITH_TOKEN,
		EVENT_EVM_TOKEN_SENT,
		EVENT_EVM_CONTRACT_CALL_APPROVED,
		EVENT_EVM_COMMAND_EXECUTED,
		//EVENT_EVM_TOKEN_DEPLOYED,
		//EVENT_EVM_REDEEM_TOKEN,
	}
)

// // Go routine for process missing logs
// func (c *EvmClient) ProcessMissingLogs() {
// 	mapEvents := map[string]abi.Event{}
// 	for _, event := range scalarGatewayAbi.Events {
// 		mapEvents[event.ID.String()] = event
// 	}
// 	for {
// 		logs := c.MissingLogs.GetLogs(10)
// 		if len(logs) == 0 {
// 			if c.MissingLogs.IsRecovered() {
// 				log.Info().Str("Chain", c.EvmConfig.ID).Msg("[EvmClient] [ProcessMissingLogs] no logs to process, recovered flag is true, exit")
// 				break
// 			} else {
// 				log.Info().Str("Chain", c.EvmConfig.ID).Msg("[EvmClient] [ProcessMissingLogs] no logs to process, recover is in progress, sleep 1 second then continue")
// 				time.Sleep(time.Second)
// 				continue
// 			}
// 		}
// 		log.Info().Str("Chain", c.EvmConfig.ID).Int("Number of logs", len(logs)).Msg("[EvmClient] [ProcessMissingLogs] processing logs")
// 		for _, txLog := range logs {
// 			topic := txLog.Topics[0].String()
// 			event, ok := mapEvents[topic]
// 			if !ok {
// 				log.Error().Str("topic", topic).Any("txLog", txLog).Msg("[EvmClient] [ProcessMissingLogs] event not found")
// 				continue
// 			}
// 			log.Debug().
// 				Str("chainId", c.EvmConfig.GetId()).
// 				Str("eventName", event.Name).
// 				Str("txHash", txLog.TxHash.String()).
// 				Msg("[EvmClient] [ProcessMissingLogs] start processing missing event")

// 			err := c.handleEventLog(event, txLog)
// 			if err != nil {
// 				log.Error().Err(err).Msg("[EvmClient] [ProcessMissingLogs] failed to handle event log")
// 			}

// 		}
// 	}
// 	//Waiting for all redeem confirm request to be handled and redeem commands ready in the pending command queue

// 	//Recover all redeem commands
// 	mapExecutingEvents := c.MissingLogs.GetExecutingEvents()
// 	groupUids := []string{}
// 	for groupUid, executingEvent := range mapExecutingEvents {
// 		log.Info().Str("Chain", c.EvmConfig.ID).Str("GroupUid", groupUid).Msgf("[EvmClient] [RecoverAllEvents] handle switched phase event to executing")
// 		c.HandleSwitchPhase(executingEvent)
// 		groupUids = append(groupUids, groupUid)
// 	}
// 	log.Info().Str("Chain", c.EvmConfig.ID).Msg("[EvmClient] [ProcessMissingLogs] finished processing all missing evm events")
// }

// // Recover all events after recovering
// func (c *EvmClient) RecoverAllEvents(ctx context.Context, groups []string) error {
// 	currentBlockNumber, err := c.Client.BlockNumber(context.Background())
// 	if err != nil {
// 		return fmt.Errorf("failed to get current block number: %w", err)
// 	}
// 	log.Info().Str("Chain", c.EvmConfig.ID).Uint64("Current BlockNumber", currentBlockNumber).
// 		Msg("[EvmClient] [RecoverAllEvents] recovering all events")

// 	//Recover switched phase event
// 	// mapPreparingEvents, mapExecutingEvents, err := c.RecoverSwitchedPhaseEvent(ctx, currentBlockNumber, groups)
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	// log.Info().Str("Chain", c.EvmConfig.ID).Msgf("[EvmClient] [RecoverAllEvents] recovered %d preparing events and %d executing events", len(mapPreparingEvents), len(mapExecutingEvents))
// 	// //First handle all preparing events
// 	// //TODO: turn on the flag Recovering
// 	// groupUids := []string{}
// 	// for groupUid, preparingEvent := range mapPreparingEvents {
// 	// 	log.Info().Str("Chain", c.EvmConfig.ID).Str("GroupUid", groupUid).Msgf("[EvmClient] [RecoverAllEvents] handle switched phase event to preparing")
// 	// 	c.HandleSwitchPhase(preparingEvent)
// 	// 	groupUids = append(groupUids, groupUid)
// 	// }
// 	// c.MissingLogs.SetLastSwitchedEvents(mapPreparingEvents, mapExecutingEvents)
// 	// //Wait for scalar network switch to preparing phase
// 	// c.WaitForSwitchingToPhase(groupUids, covExported.Preparing)
// 	//Recover all other events
// 	err = c.RecoverEvents(ctx, ALL_EVENTS, currentBlockNumber)
// 	if err != nil {
// 		return err
// 	}
// 	log.Info().Str("Chain", c.EvmConfig.ID).Msg("[EvmClient] [RecoverAllEvents] recovered all events set recovered flag to true")
// 	c.MissingLogs.SetRecovered(true)
// 	return nil
// }

// /*
// For each evm chain, we need to recover from the last Event which switch to PrepringPhase
// So, we need to recover one event Preparing if it is last
// or 2 last events, Preparing and Executing, beetween 2 this events, relayer push all redeem transactions of current session
// */
// // func (c *EvmClient) RecoverSwitchedPhaseEvent(ctx context.Context, blockNumber uint64, groups []*covExported.CustodianGroup) (
// // 	map[string]*contracts.IScalarGatewaySwitchPhase, map[string]*contracts.IScalarGatewaySwitchPhase, error) {
// // 	expectingGroups := map[string]string{}
// // 	for _, group := range groups {
// // 		groupUid := strings.TrimPrefix(group.UID.Hex(), "0x")
// // 		expectingGroups[groupUid] = group.Name
// // 	}
// // 	mapPreparingEvents := map[string]*contracts.IScalarGatewaySwitchPhase{}
// // 	mapExecutingEvents := map[string]*contracts.IScalarGatewaySwitchPhase{}
// // 	event, ok := scalarGatewayAbi.Events[EVENT_EVM_SWITCHED_PHASE]
// // 	if !ok {
// // 		return nil, nil, fmt.Errorf("switched phase event not found")
// // 	}
// // 	recoverRange := uint64(100000)
// // 	if c.EvmConfig.RecoverRange > 0 && c.EvmConfig.RecoverRange < 100000 {
// // 		recoverRange = c.EvmConfig.RecoverRange
// // 	}
// // 	var fromBlock uint64
// // 	if blockNumber < recoverRange {
// // 		fromBlock = 0
// // 	} else {
// // 		fromBlock = blockNumber - recoverRange
// // 	}
// // 	toBlock := blockNumber
// // 	for len(expectingGroups) > 0 {
// // 		query := ethereum.FilterQuery{
// // 			FromBlock: big.NewInt(int64(fromBlock)),
// // 			ToBlock:   big.NewInt(int64(toBlock)),
// // 			Addresses: []common.Address{c.GatewayAddress},
// // 			Topics:    [][]common.Hash{{event.ID}},
// // 		}
// // 		logs, err := c.Client.FilterLogs(context.Background(), query)
// // 		if err != nil {
// // 			return nil, nil, fmt.Errorf("failed to filter logs: %w", err)
// // 		}
// // 		for i := len(logs) - 1; i >= 0; i-- {
// // 			switchedPhase := &contracts.IScalarGatewaySwitchPhase{
// // 				Raw: logs[i],
// // 			}
// // 			err := parser.ParseEventData(&logs[i], event.Name, switchedPhase)
// // 			if err != nil {
// // 				return nil, nil, fmt.Errorf("failed to parse event %s: %w", event.Name, err)
// // 			}
// // 			groupUid := hex.EncodeToString(switchedPhase.CustodianGroupId[:])
// // 			groupUid = strings.TrimPrefix(groupUid, "0x")
// // 			switch switchedPhase.To {
// // 			case uint8(covExported.Preparing):
// // 				log.Info().Str("groupUid", groupUid).Msg("[EvmClient] [RecoverSwitchedPhaseEvent] found preparing event")
// // 				_, ok := mapPreparingEvents[groupUid]
// // 				if !ok {
// // 					mapPreparingEvents[groupUid] = switchedPhase
// // 				}
// // 				delete(expectingGroups, groupUid)
// // 			case uint8(covExported.Executing):
// // 				log.Info().Str("groupUid", groupUid).Msg("[EvmClient] [RecoverSwitchedPhaseEvent] found executing event")
// // 				_, ok := mapExecutingEvents[groupUid]
// // 				if !ok {
// // 					mapExecutingEvents[groupUid] = switchedPhase
// // 				}
// // 			}
// // 		}
// // 		if fromBlock <= c.EvmConfig.StartBlock {
// // 			break
// // 		}
// // 		toBlock = fromBlock - 1
// // 		if fromBlock < recoverRange+c.EvmConfig.StartBlock {
// // 			fromBlock = c.EvmConfig.StartBlock
// // 		} else {
// // 			fromBlock = fromBlock - recoverRange
// // 		}
// // 	}
// // 	if len(expectingGroups) > 0 {
// // 		return nil, nil, fmt.Errorf("some groups are not found: %v", expectingGroups)
// // 	}
// // 	return mapPreparingEvents, mapExecutingEvents, nil
// // }
// func (c *EvmClient) RecoverEvents(ctx context.Context, eventNames []string, currentBlockNumber uint64) error {
// 	topics := []common.Hash{}
// 	mapEvents := map[string]abi.Event{}
// 	for _, eventName := range eventNames {
// 		//We recover switched phase event in separate function
// 		if eventName == EVENT_EVM_SWITCHED_PHASE {
// 			continue
// 		}
// 		event, ok := scalarGatewayAbi.Events[eventName]
// 		if ok {
// 			topics = append(topics, event.ID)
// 			mapEvents[event.ID.String()] = event
// 		}
// 	}
// 	var lastCheckpoint *scalarnet.EventCheckPoint
// 	var err error
// 	if c.dbAdapter != nil {
// 		lastCheckpoint, err = c.dbAdapter.GetLastCheckPoint(c.EvmConfig.GetId(), c.EvmConfig.StartBlock)
// 		if err != nil {
// 			log.Warn().Err(err).Msgf("[EvmClient] [RecoverEvents] failed to get last checkpoint use default value")
// 		}
// 	} else {
// 		log.Warn().Msgf("[EvmClient] [RecoverEvents] dbAdapter is nil, use default value")
// 		lastCheckpoint = &scalarnet.EventCheckPoint{
// 			ChainName:   c.EvmConfig.ID,
// 			EventName:   "",
// 			BlockNumber: c.EvmConfig.StartBlock,
// 		}
// 	}
// 	log.Info().Str("Chain", c.EvmConfig.ID).
// 		Str("GatewayAddress", c.GatewayAddress.String()).
// 		Str("EventNames", strings.Join(eventNames, ",")).
// 		Any("LastCheckpoint", lastCheckpoint).Msg("[EvmClient] [RecoverEvents] start recovering events")
// 	recoverRange := uint64(100000)
// 	if c.EvmConfig.RecoverRange > 0 && c.EvmConfig.RecoverRange < 100000 {
// 		recoverRange = c.EvmConfig.RecoverRange
// 	}
// 	fromBlock := lastCheckpoint.BlockNumber
// 	logCounter := 0
// 	for fromBlock < currentBlockNumber {
// 		query := ethereum.FilterQuery{
// 			FromBlock: big.NewInt(int64(fromBlock)),
// 			Addresses: []common.Address{c.GatewayAddress},
// 			Topics:    [][]common.Hash{topics},
// 		}
// 		if fromBlock+recoverRange < currentBlockNumber {
// 			query.ToBlock = big.NewInt(int64(fromBlock + recoverRange))
// 		} else {
// 			query.ToBlock = big.NewInt(int64(currentBlockNumber))
// 		}
// 		log.Info().Str("Chain", c.EvmConfig.ID).Msgf("[EvmClient] [RecoverEvents] querying logs fromBlock: %d, toBlock: %d", fromBlock, query.ToBlock.Uint64())
// 		logs, err := c.Client.FilterLogs(context.Background(), query)
// 		if err != nil {
// 			return fmt.Errorf("failed to filter logs: %w", err)
// 		}
// 		if len(logs) > 0 {
// 			log.Info().Str("Chain", c.EvmConfig.ID).Msgf("[EvmClient] [RecoverEvents] found %d logs, fromBlock: %d, toBlock: %d", len(logs), fromBlock, query.ToBlock)
// 			c.MissingLogs.AppendLogs(logs)
// 			logCounter += len(logs)
// 			if c.dbAdapter != nil {
// 				c.UpdateLastCheckPoint(mapEvents, logs, query.ToBlock.Uint64())
// 			}
// 		} else {
// 			log.Info().Str("Chain", c.EvmConfig.ID).Msgf("[EvmClient] [RecoverEvents] no logs found, fromBlock: %d, toBlock: %d", fromBlock, query.ToBlock)
// 		}
// 		//Set fromBlock to the next block number for next iteration
// 		fromBlock = query.ToBlock.Uint64() + 1
// 	}
// 	log.Info().
// 		Str("Chain", c.EvmConfig.ID).
// 		Uint64("CurrentBlockNumber", currentBlockNumber).
// 		Int("TotalLogs", logCounter).
// 		Msg("[EvmClient] [FinishRecover] recovered all events")
// 	return nil
// }

// func (c *EvmClient) UpdateLastCheckPoint(events map[string]abi.Event, logs []types.Log, lastBlock uint64) {
// 	eventCheckPoints := map[string]scalarnet.EventCheckPoint{}
// 	for _, txLog := range logs {
// 		topic := txLog.Topics[0].String()
// 		event, ok := events[topic]
// 		if !ok {
// 			log.Error().Str("topic", topic).Any("txLog", txLog).Msg("[EvmClient] [UpdateLastCheckPoint] event not found")
// 			continue
// 		}
// 		checkpoint, ok := eventCheckPoints[event.Name]
// 		if !ok {
// 			checkpoint = scalarnet.EventCheckPoint{
// 				ChainName:   c.EvmConfig.ID,
// 				EventName:   event.Name,
// 				BlockNumber: lastBlock,
// 				LogIndex:    txLog.Index,
// 				TxHash:      txLog.TxHash.String(),
// 			}
// 			eventCheckPoints[event.Name] = checkpoint
// 		} else {
// 			checkpoint.BlockNumber = lastBlock
// 			checkpoint.LogIndex = txLog.Index
// 			checkpoint.TxHash = txLog.TxHash.String()
// 		}
// 	}
// 	c.dbAdapter.UpdateLastEventCheckPoints(eventCheckPoints)
// }

/*
For each evm chain, we need to get 2 last switch phase events
so we have to case, [Preparing, Executing], [Executing, Preparing] or [Preparing]
*/

func (c *EvmClient) RecoverRedeemSessions(ctx context.Context, groups []common.Hash) (*ChainRedeemSessions, error) {

	redeemSessions := NewChainRedeemSessions()

	for _, gr := range groups {
		query := ethereum.FilterQuery{
			Addresses: []common.Address{c.GatewayAddress},
			Topics:    [][]common.Hash{{switchPhaseEvent.ID}, {gr}},
			FromBlock: big.NewInt(int64(c.EvmConfig.StartBlock)),
		}

		events, err := c.Client.FilterLogs(ctx, query)
		if err != nil {
			return nil, fmt.Errorf("failed to filter logs: %w", err)
		}

		if len(events) == 0 {
			log.Warn().
				Str("Chain", c.EvmConfig.ID).
				Msg("[EvmClient] [RecoverRedeemSessions] no events found")
			return nil, fmt.Errorf("[EvmClient] [RecoverRedeemSessions]: no events found")
		}

		for i := len(events) - 1; i >= 0; i-- {
			event := events[i]
			switchedPhase, err := parseSwitchPhaseEvent(event)
			if err != nil {
				log.Error().Err(err).Msgf("[EvmClient] [RecoverRedeemSessions] failed to parse event %s", switchPhaseEvent.Name)
				return nil, fmt.Errorf("[EvmClient] [RecoverRedeemSessions]: failed to parse event %s", switchPhaseEvent.Name)
			}

			counter := redeemSessions.AppendSwitchPhaseEvent(gr.Hex(), switchedPhase)
			if counter == 2 || switchedPhase.To == uint8(db.Preparing) && switchedPhase.Sequence == 1 {
				break
			}
		}

	}

	return redeemSessions, nil
}

func parseSwitchPhaseEvent(log types.Log) (*contracts.IScalarGatewaySwitchPhase, error) {
	eventName := EVENT_EVM_SWITCHED_PHASE
	switchedPhase := &contracts.IScalarGatewaySwitchPhase{
		Raw: log,
	}
	err := ParseEventData(&log, eventName, switchedPhase)
	if err != nil {
		return nil, fmt.Errorf("failed to parse event %s: %w", eventName, err)
	}
	return switchedPhase, nil
}

func parseRedeemTokenEvent(log types.Log) (*contracts.IScalarGatewayRedeemToken, error) {
	eventName := EVENT_EVM_REDEEM_TOKEN
	redeemToken := &contracts.IScalarGatewayRedeemToken{
		Raw: log,
	}
	err := ParseEventData(&log, eventName, redeemToken)
	if err != nil {
		return nil, fmt.Errorf("failed to parse event %s: %w", eventName, err)
	}
	return redeemToken, nil
}
