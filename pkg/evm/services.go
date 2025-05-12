package evm

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	"github.com/scalarorg/relayers/pkg/clients/evm/parser"
	contracts "github.com/scalarorg/scalar-healer/pkg/evm/contracts/generated"
)

var abiEventsMap = map[string]abi.Event{}

func init() {
	for _, event := range scalarGatewayAbi.Events {
		abiEventsMap[event.Name] = event
	}
}

// Go routine for process missing logs
func (c *EvmClient) ProcessMissingLogs(ctx context.Context) {
	for {
		logs := c.MissingLogs.GetLogs(10)
		if len(logs) == 0 {
			if c.MissingLogs.IsRecovered() {
				log.Info().Str("Chain", c.EvmConfig.ID).Msg("[EvmClient] [ProcessMissingLogs] no logs to process, recovered flag is true, exit")
				break
			} else {
				log.Info().Str("Chain", c.EvmConfig.ID).Msg("[EvmClient] [ProcessMissingLogs] no logs to process, recover is in progress, sleep 1 second then continue")
				time.Sleep(time.Second)
				continue
			}
		}
		log.Info().
			Str("Chain", c.EvmConfig.ID).
			Int("Number of logs", len(logs)).
			Msg("[EvmClient] [ProcessMissingLogs] processing logs")

		for _, txLog := range logs {
			topic := txLog.Topics[0].String()
			event, ok := mapEvents[topic]
			if !ok {
				log.Error().
					Str("topic", topic).
					Any("txLog", txLog).
					Msg("[EvmClient] [ProcessMissingLogs] event not found")
				continue
			}
			log.Debug().
				Str("chainId", c.EvmConfig.GetId()).
				Str("eventName", event.Name).
				Str("txHash", txLog.TxHash.String()).
				Msg("[EvmClient] [ProcessMissingLogs] start processing missing event")

			err := c.handleEventLog(ctx, event, txLog)
			if err != nil {
				log.Error().Err(err).Msg("[EvmClient] [ProcessMissingLogs] failed to handle event log")
			}
		}
	}
}

func (c *EvmClient) handleEventLog(ctx context.Context, event *abi.Event, txLog types.Log) error {
	switch event.Name {
	case EVENT_EVM_TOKEN_SENT:
		tokenSent := &contracts.IScalarGatewayTokenSent{
			Raw: txLog,
		}
		err := parser.ParseEventData(&txLog, event.Name, tokenSent)
		if err != nil {
			return fmt.Errorf("failed to parse event %s: %w", event.Name, err)
		}
		return c.HandleTokenSent(ctx, tokenSent)
	case EVENT_EVM_REDEEM_TOKEN:
		redeemToken := &contracts.IScalarGatewayRedeemToken{
			Raw: txLog,
		}
		err := parser.ParseEventData(&txLog, event.Name, redeemToken)
		if err != nil {
			return fmt.Errorf("failed to parse event %s: %w", event.Name, err)
		}
		return c.HandleRedeemToken(ctx, redeemToken)
	default:
		return fmt.Errorf("invalid event type for %s: %T", event.Name, txLog)
	}
}
