package evm

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	"github.com/scalarorg/data-models/chains"
	contracts "github.com/scalarorg/scalar-healer/pkg/evm/contracts/generated"
)

func (ec *EvmClient) HandleContractCallWithToken(ctx context.Context, event *contracts.IScalarGatewayContractCallWithToken) error {
	//0. Preprocess the event
	ec.preprocessContractCallWithToken(event)
	//1. Convert into a RelayData instance then store to the db
	contractCallWithToken, err := ec.ContractCallWithToken2Model(event)
	if err != nil {
		return fmt.Errorf("failed to convert ContractCallEvent to ContractCallWithToken: %w", err)
	}
	//2. update last checkpoint
	lastCheckpoint, err := ec.dbAdapter.GetLastEventCheckPoint(ctx, ec.EvmConfig.GetId(), EVENT_EVM_CONTRACT_CALL_WITH_TOKEN, ec.EvmConfig.StartBlock)
	if err != nil {
		log.Debug().Str("chainId", ec.EvmConfig.GetId()).
			Str("eventName", EVENT_EVM_CONTRACT_CALL_WITH_TOKEN).
			Msg("[EvmClient] [handleContractCallWithToken] Get event from begining")
	}
	if event.Raw.BlockNumber > lastCheckpoint.BlockNumber ||
		(event.Raw.BlockNumber == lastCheckpoint.BlockNumber && event.Raw.TxIndex > lastCheckpoint.LogIndex) {
		lastCheckpoint.BlockNumber = event.Raw.BlockNumber
		lastCheckpoint.TxHash = event.Raw.TxHash.String()
		lastCheckpoint.LogIndex = event.Raw.Index
		lastCheckpoint.EventKey = fmt.Sprintf("%s-%d-%d", event.Raw.TxHash.String(), event.Raw.BlockNumber, event.Raw.Index)
	}
	//3. store relay data to the db, update last checkpoint
	err = ec.dbAdapter.SaveContractCallWithToken(ctx, &contractCallWithToken, lastCheckpoint)
	if err != nil {
		return fmt.Errorf("failed to create evm contract call: %w", err)
	}
	return ec.AddContractCallWithTokenToEvmBatchCommand(event)
}
func (ec *EvmClient) AddContractCallWithTokenToEvmBatchCommand(event *contracts.IScalarGatewayContractCallWithToken) error {
	//TODO: implement this
	return nil
}
func (ec *EvmClient) HandleRedeemToken(ctx context.Context, event *contracts.IScalarGatewayRedeemToken) error {
	//0. Preprocess the event
	log.Info().Str("Chain", ec.EvmConfig.ID).Any("event", event).Msg("[EvmClient] [HandleRedeemToken] Start processing evm redeem token")
	//1. Convert into a RelayData instance then store to the db
	redeemToken, err := ec.RedeemTokenEvent2Model(event)
	if err != nil {
		return fmt.Errorf("failed to convert ContractCallEvent to ContractCallWithToken: %w", err)
	}
	//2. update last checkpoint
	lastCheckpoint, err := ec.dbAdapter.GetLastEventCheckPoint(ctx, ec.EvmConfig.GetId(), EVENT_EVM_REDEEM_TOKEN, ec.EvmConfig.StartBlock)
	if err != nil {
		log.Debug().Str("chainId", ec.EvmConfig.GetId()).
			Str("eventName", EVENT_EVM_CONTRACT_CALL_WITH_TOKEN).
			Msg("[EvmClient] [handleContractCallWithToken] Get event from begining")
	}
	if event.Raw.BlockNumber > lastCheckpoint.BlockNumber ||
		(event.Raw.BlockNumber == lastCheckpoint.BlockNumber && event.Raw.TxIndex > lastCheckpoint.LogIndex) {
		lastCheckpoint.BlockNumber = event.Raw.BlockNumber
		lastCheckpoint.TxHash = event.Raw.TxHash.String()
		lastCheckpoint.LogIndex = event.Raw.Index
		lastCheckpoint.EventKey = fmt.Sprintf("%s-%d-%d", event.Raw.TxHash.String(), event.Raw.BlockNumber, event.Raw.Index)
	}
	//3. store relay data to the db, update last checkpoint
	err = ec.dbAdapter.SaveContractCallWithToken(ctx, &redeemToken, lastCheckpoint)
	if err != nil {
		return fmt.Errorf("failed to create evm contract call: %w", err)
	}
	return ec.AddRedeemTokenToRedeemBatchPsbt(event)
}
func (ec *EvmClient) AddRedeemTokenToRedeemBatchPsbt(event *contracts.IScalarGatewayRedeemToken) error {
	//TODO: implement this
	return nil
}
func (ec *EvmClient) preprocessContractCallWithToken(event *contracts.IScalarGatewayContractCallWithToken) error {
	log.Info().
		Str("sender", event.Sender.Hex()).
		Str("destinationChain", event.DestinationChain).
		Str("destinationContractAddress", event.DestinationContractAddress).
		Str("payloadHash", hex.EncodeToString(event.PayloadHash[:])).
		Str("Symbol", event.Symbol).
		Uint64("Amount", event.Amount.Uint64()).
		Str("txHash", event.Raw.TxHash.String()).
		Uint("logIndex", event.Raw.Index).
		Uint("txIndex", event.Raw.TxIndex).
		Str("logData", hex.EncodeToString(event.Raw.Data)).
		Msg("[EvmClient] [preprocessContractCallWithToken] Start handle Contract call with token")
	//Todo: validate the event
	return nil
}

func (ec *EvmClient) HandleTokenSent(ctx context.Context, event *contracts.IScalarGatewayTokenSent) error {
	//0. Preprocess the event
	ec.preprocessTokenSent(event)
	//1. Convert into a RelayData instance then store to the db
	tokenSent, err := ec.TokenSentEvent2Model(event)
	if err != nil {
		log.Error().Err(err).Msg("[EvmClient] [HandleTokenSent] failed to convert TokenSentEvent to model data")
		return err
	}
	//For evm, the token sent is verified immediately by the scalarnet
	tokenSent.Status = chains.TokenSentStatusVerifying
	//2. update last checkpoint
	lastCheckpoint, err := ec.dbAdapter.GetLastEventCheckPoint(ctx, ec.EvmConfig.GetId(), EVENT_EVM_TOKEN_SENT, ec.EvmConfig.StartBlock)
	if err != nil {
		log.Debug().Str("chainId", ec.EvmConfig.GetId()).
			Str("eventName", EVENT_EVM_TOKEN_SENT).
			Msg("[EvmClient] [handleTokenSent] Get event from begining")
	}
	if event.Raw.BlockNumber > lastCheckpoint.BlockNumber ||
		(event.Raw.BlockNumber == lastCheckpoint.BlockNumber && event.Raw.TxIndex > lastCheckpoint.LogIndex) {
		lastCheckpoint.BlockNumber = event.Raw.BlockNumber
		lastCheckpoint.TxHash = event.Raw.TxHash.String()
		lastCheckpoint.LogIndex = event.Raw.Index
		lastCheckpoint.EventKey = fmt.Sprintf("%s-%d-%d", event.Raw.TxHash.String(), event.Raw.BlockNumber, event.Raw.Index)
	}
	//3. store relay data to the db, update last checkpoint
	err = ec.dbAdapter.SaveTokenSent(ctx, &tokenSent, lastCheckpoint)
	if err != nil {
		return fmt.Errorf("failed to create evm token send: %w", err)
	}
	return ec.AddTokenSentToEvmBatchCommand(event)
}

func (ec *EvmClient) AddTokenSentToEvmBatchCommand(event *contracts.IScalarGatewayTokenSent) error {
	//TODO: implement this
	return nil
}
func (ec *EvmClient) preprocessTokenSent(event *contracts.IScalarGatewayTokenSent) error {
	log.Info().
		Str("sender", event.Sender.Hex()).
		Str("destinationChain", event.DestinationChain).
		Str("destinationAddress", event.DestinationAddress).
		Str("txHash", event.Raw.TxHash.String()).
		Str("symbol", event.Symbol).
		Uint64("amount", event.Amount.Uint64()).
		Uint("logIndex", event.Raw.Index).
		Uint("txIndex", event.Raw.TxIndex).
		Str("logData", hex.EncodeToString(event.Raw.Data)).
		Msg("[EvmClient] [preprocessTokenSent] Start handle TokenSent")
	//Todo: validate the event
	return nil
}

// func (ec *EvmClient) handleContractCall(event *contracts.IScalarGatewayContractCall) error {
// 	//0. Preprocess the event
// 	ec.preprocessContractCall(event)
// 	//1. Convert into a RelayData instance then store to the db
// 	contractCall, err := ec.ContractCallEvent2Model(event)
// 	if err != nil {
// 		return fmt.Errorf("failed to convert ContractCallEvent to RelayData: %w", err)
// 	}
// 	//2. update last checkpoint
// 	lastCheckpoint, err := ec.dbAdapter.GetLastEventCheckPoint(ec.EvmConfig.GetId(), config.EVENT_EVM_CONTRACT_CALL)
// 	if err != nil {
// 		log.Debug().Str("chainId", ec.EvmConfig.GetId()).
// 			Str("eventName", config.EVENT_EVM_CONTRACT_CALL).
// 			Msg("[EvmClient] [handleContractCall] Get event from begining")
// 	}
// 	if event.Raw.BlockNumber > lastCheckpoint.BlockNumber ||
// 		(event.Raw.BlockNumber == lastCheckpoint.BlockNumber && event.Raw.TxIndex > lastCheckpoint.LogIndex) {
// 		lastCheckpoint.BlockNumber = event.Raw.BlockNumber
// 		lastCheckpoint.TxHash = event.Raw.TxHash.String()
// 		lastCheckpoint.LogIndex = event.Raw.Index
// 		lastCheckpoint.EventKey = fmt.Sprintf("%s-%d-%d", event.Raw.TxHash.String(), event.Raw.BlockNumber, event.Raw.Index)
// 	}
// 	//3. store relay data to the db, update last checkpoint
// 	err = ec.dbAdapter.CreateContractCall(contractCall, lastCheckpoint)
// 	if err != nil {
// 		return fmt.Errorf("failed to create evm contract call: %w", err)
// 	}
// 	//2. Send to the bus
// 	confirmTxs := events.ConfirmTxsRequest{
// 		ChainName: ec.EvmConfig.GetId(),
// 		TxHashs:   map[string]string{contractCall.TxHash: contractCall.DestinationChain},
// 	}
// 	if ec.eventBus != nil {
// 		ec.eventBus.BroadcastEvent(&events.EventEnvelope{
// 			EventType:        events.EVENT_EVM_CONTRACT_CALL,
// 			DestinationChain: events.SCALAR_NETWORK_NAME,
// 			Data:             confirmTxs,
// 		})
// 	} else {
// 		log.Warn().Msg("[EvmClient] [handleContractCall] event bus is undefined")
// 	}
// 	return nil
// }
// func (ec *EvmClient) preprocessContractCall(event *contracts.IScalarGatewayContractCall) error {
// 	log.Info().
// 		Str("sender", event.Sender.Hex()).
// 		Str("destinationChain", event.DestinationChain).
// 		Str("destinationContractAddress", event.DestinationContractAddress).
// 		Str("payloadHash", hex.EncodeToString(event.PayloadHash[:])).
// 		Str("txHash", event.Raw.TxHash.String()).
// 		Uint("logIndex", event.Raw.Index).
// 		Uint("txIndex", event.Raw.TxIndex).
// 		Str("logData", hex.EncodeToString(event.Raw.Data)).
// 		Msg("[EvmClient] [preprocessContractCall] Start handle Contract call")
// 	//Todo: validate the event
// 	return nil
// }

//	func (ec *EvmClient) HandleContractCallApproved(event *contracts.IScalarGatewayContractCallApproved) error {
//		//0. Preprocess the event
//		err := ec.preprocessContractCallApproved(event)
//		if err != nil {
//			return fmt.Errorf("failed to preprocess contract call approved: %w", err)
//		}
//		//1. Convert into a RelayData instance then store to the db
//		contractCallApproved, err := ec.ContractCallApprovedEvent2Model(event)
//		if err != nil {
//			return fmt.Errorf("failed to convert ContractCallApprovedEvent to RelayData: %w", err)
//		}
//		err = ec.dbAdapter.SaveSingleValue(&contractCallApproved)
//		if err != nil {
//			return fmt.Errorf("failed to create contract call approved: %w", err)
//		}
//		// Find relayData from the db by combination (contractAddress, sourceAddress, payloadHash)
//		// This contract call (initiated by the user call to the source chain) is approved by EVM network
//		// So anyone can execute it on the EVM by broadcast the corresponding payload to protocol's smart contract on the destination chain
//		destContractAddress := strings.TrimLeft(event.ContractAddress.Hex(), "0x")
//		sourceAddress := strings.TrimLeft(event.SourceAddress, "0x")
//		payloadHash := strings.TrimLeft(hex.EncodeToString(event.PayloadHash[:]), "0x")
//		relayDatas, err := ec.dbAdapter.FindContractCallByParams(sourceAddress, destContractAddress, payloadHash)
//		if err != nil {
//			log.Error().Err(err).Msg("[EvmClient] [handleContractCallApproved] find relay data")
//			return err
//		}
//		log.Debug().Str("contractAddress", event.ContractAddress.String()).
//			Str("sourceAddress", event.SourceAddress).
//			Str("payloadHash", hex.EncodeToString(event.PayloadHash[:])).
//			Any("relayDatas count", len(relayDatas)).
//			Msg("[EvmClient] [handleContractCallApproved] query relaydata by ContractCall")
//		//3. Execute payload in the found relaydatas
//		executeResults, err := ec.executeDestinationCall(event, relayDatas)
//		if err != nil {
//			log.Warn().Err(err).Any("executeResults", executeResults).Msg("[EvmClient] [handleContractCallApproved] execute destination call")
//		}
//		// Done; Don't need to send to the bus
//		// TODO: Do we need to update relay data atomically?
//		err = ec.dbAdapter.UpdateBatchContractCallStatus(executeResults, len(executeResults))
//		if err != nil {
//			return fmt.Errorf("failed to update relay data status to executed: %w", err)
//		}
//		return nil
//	}
// func (ec *EvmClient) executeDestinationCall(event *contracts.IScalarGatewayContractCallApproved, contractCalls []chains.ContractCall) ([]db.ContractCallExecuteResult, error) {
// 	executeResults := []db.ContractCallExecuteResult{}
// 	executed, err := ec.isContractCallExecuted(event)
// 	if err != nil {
// 		return executeResults, fmt.Errorf("[EvmClient] [executeDestinationCall] failed to check if contract call is approved: %w", err)
// 	}
// 	if executed {
// 		//Update the relay data status to executed
// 		for _, contractCall := range contractCalls {
// 			executeResults = append(executeResults, db.ContractCallExecuteResult{
// 				Status:  chains.ContractCallStatusSuccess,
// 				EventId: contractCall.EventID,
// 			})
// 		}
// 		return executeResults, fmt.Errorf("destination contract call is already executed")
// 	}
// 	if len(contractCalls) > 0 {
// 		for _, contractCall := range contractCalls {
// 			if len(contractCall.Payload) == 0 {
// 				continue
// 			}
// 			log.Info().Str("payload", hex.EncodeToString(contractCall.Payload)).
// 				Msg("[EvmClient] [executeDestinationCall]")
// 			receipt, err := ec.ExecuteDestinationCall(event.ContractAddress, event.CommandId, event.SourceChain, event.SourceAddress, contractCall.Payload)
// 			if err != nil {
// 				return executeResults, fmt.Errorf("execute destination call with error: %w", err)
// 			}

// 			log.Info().Any("txReceipt", receipt).Msg("[EvmClient] [executeDestinationCall]")

// 			if receipt.Hash() != (common.Hash{}) {
// 				executeResults = append(executeResults, db.ContractCallExecuteResult{
// 					Status:  chains.ContractCallStatusSuccess,
// 					EventId: contractCall.EventID,
// 				})
// 			} else {
// 				executeResults = append(executeResults, db.ContractCallExecuteResult{
// 					Status:  chains.ContractCallStatusFailed,
// 					EventId: contractCall.EventID,
// 				})
// 			}
// 		}
// 	}
// 	return executeResults, nil
// }

// Check if ContractCall is already executed
func (ec *EvmClient) isContractCallExecuted(event *contracts.IScalarGatewayContractCallApproved) (bool, error) {
	if ec.transactOpts == nil {
		log.Error().
			Str("commandId", hex.EncodeToString(event.CommandId[:])).
			Str("sourceChain", event.SourceChain).
			Str("sourceAddress", event.SourceAddress).
			Str("contractAddress", event.ContractAddress.String()).
			Msg("[EvmClient] [isContractCallExecuted] auth is nil")
		return false, fmt.Errorf("auth is nil")
	}
	callOpt := &bind.CallOpts{
		From:    ec.transactOpts.From,
		Context: context.Background(),
	}
	approved, err := ec.Gateway.IsContractCallApproved(callOpt, event.CommandId, event.SourceChain, event.SourceAddress, event.ContractAddress, event.PayloadHash)
	if err != nil {
		return false, fmt.Errorf("failed to check if contract call is approved: %w", err)
	}
	return !approved, nil
}

func (ec *EvmClient) preprocessContractCallApproved(event *contracts.IScalarGatewayContractCallApproved) error {
	log.Info().Any("event", event).Msgf("[EvmClient] [handleContractCallApproved]")
	//Todo: validate the event
	return nil
}

func (ec *EvmClient) HandleCommandExecuted(ctx context.Context, event *contracts.IScalarGatewayExecuted) error {
	//0. Preprocess the event
	// ec.preprocessCommandExecuted(event)
	// //1. Convert into a RelayData instance then store to the db
	// cmdExecuted := ec.CommandExecutedEvent2Model(event)
	// err := ec.dbAdapter.UpdateEvmCommandExecuted(ctx, &cmdExecuted)
	// if err != nil {
	// 	log.Error().Err(err).Msg("[EvmClient] HandleCommandExecuted failed")
	// }
	return nil
}

func (ec *EvmClient) preprocessCommandExecuted(event *contracts.IScalarGatewayExecuted) error {
	log.Info().Any("event", event).Msg("[EvmClient] [ExecutedHandler] Start processing evm command executed")
	//Todo: validate the event
	return nil
}

func (ec *EvmClient) SubmitTx(signedTx *ethtypes.Transaction, retryAttempt int) (*ethtypes.Receipt, error) {
	if retryAttempt >= ec.EvmConfig.MaxRetry {
		return nil, fmt.Errorf("max retry exceeded")
	}

	// Create a new context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), ec.EvmConfig.TxTimeout)
	defer cancel()

	// Log transaction details
	log.Debug().
		Interface("tx", signedTx).
		Msg("Submitting transaction")

	// Send the transaction using the new context
	err := ec.Client.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Error().
			Err(err).
			Str("rpcUrl", ec.EvmConfig.RPCUrl).
			Str("walletAddress", ec.transactOpts.From.String()).
			Str("to", signedTx.To().String()).
			Str("data", hex.EncodeToString(signedTx.Data())).
			Msg("[EvmClient.SubmitTx] Failed to submit transaction")

		// Sleep before retry
		time.Sleep(ec.EvmConfig.RetryDelay)

		log.Debug().
			Int("attempt", retryAttempt+1).
			Msg("Retrying transaction")

		return ec.SubmitTx(signedTx, retryAttempt+1)
	}

	// Wait for transaction receipt using the new context
	receipt, err := bind.WaitMined(ctx, ec.Client, signedTx)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for transaction receipt: %w", err)
	}

	log.Debug().
		Interface("receipt", receipt).
		Msg("Transaction receipt received")

	return receipt, nil
}

func (ec *EvmClient) WaitForTransaction(hash string) (*ethtypes.Receipt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ec.EvmConfig.TxTimeout)
	defer cancel()

	txHash := common.HexToHash(hash)
	tx, _, err := ec.Client.TransactionByHash(ctx, txHash)
	if err != nil {
		return nil, err
	}
	return bind.WaitMined(ctx, ec.Client, tx)
}
