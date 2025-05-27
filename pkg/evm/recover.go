package evm

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	"github.com/scalarorg/scalar-healer/pkg/db/sqlc"
	contracts "github.com/scalarorg/scalar-healer/pkg/evm/contracts/generated"
)

var (
	ALL_EVENTS = []string{
		EVENT_EVM_CONTRACT_CALL,
		EVENT_EVM_CONTRACT_CALL_WITH_TOKEN,
		EVENT_EVM_TOKEN_SENT,
		EVENT_EVM_CONTRACT_CALL_APPROVED,
		EVENT_EVM_COMMAND_EXECUTED,
	}
)

var abiEventsMap = map[string]abi.Event{}
var switchPhaseEvent abi.Event
var redeemTokenEvent abi.Event

func init() {
	for _, event := range scalarGatewayAbi.Events {
		abiEventsMap[event.Name] = event
	}

	switchPhaseEvent = abiEventsMap[EVENT_EVM_SWITCHED_PHASE]
	redeemTokenEvent = abiEventsMap[EVENT_EVM_REDEEM_TOKEN]
}

func (c *EvmClient) RecoverRedeemSessions(ctx context.Context, groups []common.Hash) (*ChainRedeemSessions, error) {

	redeemSessions := NewChainRedeemSessions()

	// just get events in range 500 blocks to prevent rate limit
	currentBlockNumber, err := c.Client.BlockNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get current block number: %w", err)
	}

	baseQuery := ethereum.FilterQuery{
		Addresses: []common.Address{c.GatewayAddress},
		ToBlock:   big.NewInt(int64(currentBlockNumber)),
		FromBlock: big.NewInt(int64(currentBlockNumber - uint64(c.EvmConfig.RecoverRange) + 1)),
	}

	log.Info().
		Str("Number of groups", fmt.Sprintf("%d", len(groups))).
		Msg("[EvmClient] [RecoverRedeemSessions] start recovering redeem sessions")

	for _, gr := range groups {
		isEnoughEvents := false
		query := baseQuery
		query.Topics = [][]common.Hash{{switchPhaseEvent.ID}, {gr}}
		for !isEnoughEvents && query.FromBlock.Uint64() > c.EvmConfig.StartBlock {
			log.Info().
				Str("Chain", c.EvmConfig.ID).
				Str("Group", gr.String()).
				Msgf("[EvmClient] [RecoverRedeemSessions] querying logs fromBlock: %d, toBlock: %d", query.FromBlock.Uint64(), query.ToBlock.Uint64())
			events, err := c.Client.FilterLogs(ctx, query)
			if err != nil {
				return nil, fmt.Errorf("failed to filter logs: %w", err)
			}

			if len(events) > 0 {
				for i := len(events) - 1; i >= 0; i-- {
					event := events[i]
					switchedPhaseEvent, err := parseSwitchPhaseEvent(event)
					if err != nil {
						log.Error().Err(err).Msgf("[EvmClient] [RecoverRedeemSessions] failed to parse event %s", switchPhaseEvent.Name)
						return nil, fmt.Errorf("[EvmClient] [RecoverRedeemSessions]: failed to parse event %s", switchPhaseEvent.Name)
					}

					counter := redeemSessions.AppendSwitchPhaseEvent(gr.Hex(), switchedPhaseEvent)
					if counter == 2 || switchedPhaseEvent.To == sqlc.RedeemPhasePREPARING.Uint8() && switchedPhaseEvent.Sequence == 1 {
						log.Debug().
							Str("Chain", c.EvmConfig.ID).
							Str("Group", gr.String()).
							Msgf("[EvmClient] [RecoverRedeemSessions] found %d events", counter)
						isEnoughEvents = true
						break
					}
				}
			}

			query.ToBlock = big.NewInt(int64(query.ToBlock.Uint64() - uint64(c.EvmConfig.RecoverRange)))
			query.FromBlock = big.NewInt(int64(query.FromBlock.Uint64() - uint64(c.EvmConfig.RecoverRange)))
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
