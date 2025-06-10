package evm

// import (
// 	"context"
// 	"encoding/hex"
// 	"fmt"
// 	"strings"

// 	chainTypes "github.com/scalarorg/bitcoin-vault/go-utils/chain"
// 	chains "github.com/scalarorg/data-models/chains"
// 	"github.com/scalarorg/data-models/scalarnet"
// 	contracts "github.com/scalarorg/scalar-healer/pkg/evm/contracts/generated"
// 	"github.com/scalarorg/scalar-healer/pkg/utils"
// )

// func (c *EvmClient) ContractCallEvent2Model(event *contracts.IScalarGatewayContractCall) (chains.ContractCall, error) {
// 	//id := strings.ToLower(fmt.Sprintf("%s-%d", event.Raw.TxHash.String(), event.Raw.Index))
// 	//Calculate eventId by Txhash-logIndex among logs in txreceipt (AxelarEvmModule)
// 	//https://github.com/scalarorg/scalar-core/blob/main/vald/evm/gateway_tx_confirmation.go#L73
// 	//Dec 30, use logIndex directly to avoid redundant request. This must aggrees with the scalar-core vald module
// 	// receipt, err := c.Client.TransactionReceipt(context.Background(), event.Raw.TxHash)
// 	// if err != nil {
// 	// 	return models.RelayData{}, fmt.Errorf("failed to get transaction receipt: %w", err)
// 	// }
// 	// var id string
// 	// for ind, log := range receipt.Logs {
// 	// 	if log.Index == event.Raw.Index {
// 	// 		id = fmt.Sprintf("%s-%d", event.Raw.TxHash.String(), ind)
// 	// 		break
// 	// 	}
// 	// }
// 	eventId := fmt.Sprintf("%s-%d", utils.NormalizeHash(event.Raw.TxHash.String()), event.Raw.Index)
// 	senderAddress := event.Sender.String()

// 	chainInfoBytes := chainTypes.ChainInfoBytes{}
// 	err := chainInfoBytes.FromString(event.DestinationChain)
// 	if err != nil {
// 		return chains.ContractCall{}, fmt.Errorf("failed to convert destination chain: %w", err)
// 	}
// 	contractCall := chains.ContractCall{
// 		EventID:     eventId,
// 		TxHash:      event.Raw.TxHash.String(),
// 		BlockNumber: event.Raw.BlockNumber,
// 		LogIndex:    event.Raw.Index,
// 		SourceChain: c.EvmConfig.GetId(),
// 		//3 follows field are used for query to get back payload, so need to convert to lower case
// 		DestinationChain: event.DestinationChain,
// 		SourceAddress:    utils.NormalizeAddress(senderAddress, chainInfoBytes.ChainType()),
// 		PayloadHash:      utils.NormalizeHash(hex.EncodeToString(event.PayloadHash[:])),
// 		Payload:          event.Payload,
// 	}
// 	return contractCall, nil
// }

// func (c *EvmClient) ContractCallWithToken2Model(event *contracts.IScalarGatewayContractCallWithToken) (chains.ContractCallWithToken, error) {
// 	eventId := fmt.Sprintf("%s-%d", utils.NormalizeHash(event.Raw.TxHash.String()), event.Raw.Index)
// 	senderAddress := event.Sender.String()

// 	chainInfoBytes := chainTypes.ChainInfoBytes{}
// 	err := chainInfoBytes.FromString(event.DestinationChain)
// 	if err != nil {
// 		return chains.ContractCallWithToken{}, fmt.Errorf("failed to convert destination chain: %w", err)
// 	}

// 	destinationAddress, err := utils.CalculateDestinationAddress(event.Payload, &chainInfoBytes)
// 	if err != nil {
// 		return chains.ContractCallWithToken{}, fmt.Errorf("failed to calculate destination address: %w", err)
// 	}

// 	callContract := chains.ContractCall{
// 		EventID:     eventId,
// 		TxHash:      utils.NormalizeHash(event.Raw.TxHash.String()),
// 		BlockNumber: event.Raw.BlockNumber,
// 		LogIndex:    event.Raw.Index,
// 		SourceChain: c.EvmConfig.GetId(),
// 		//3 follows field are used for query to get back payload, so need to convert to lower case
// 		DestinationChain:   event.DestinationChain,
// 		DestinationAddress: utils.NormalizeAddress(destinationAddress, chainInfoBytes.ChainType()),
// 		SourceAddress:      utils.NormalizeAddress(senderAddress, chainInfoBytes.ChainType()),
// 		PayloadHash:        utils.NormalizeHash(hex.EncodeToString(event.PayloadHash[:])),
// 		Payload:            event.Payload,
// 	}
// 	contractCallWithToken := chains.ContractCallWithToken{
// 		ContractCall:         callContract,
// 		TokenContractAddress: utils.NormalizeAddress(event.DestinationContractAddress, chainInfoBytes.ChainType()),
// 		Symbol:               event.Symbol,
// 		Amount:               event.Amount.Uint64(),
// 	}
// 	return contractCallWithToken, nil
// }

// // Use contractCallToken to store the redeem token event
// func (c *EvmClient) RedeemTokenEvent2Model(event *contracts.IScalarGatewayRedeemToken) (chains.ContractCallWithToken, error) {
// 	eventId := fmt.Sprintf("%s-%d", utils.NormalizeHash(event.Raw.TxHash.String()), event.Raw.Index)
// 	senderAddress := event.Sender.String()

// 	chainInfoBytes := chainTypes.ChainInfoBytes{}
// 	err := chainInfoBytes.FromString(event.DestinationChain)
// 	if err != nil {
// 		return chains.ContractCallWithToken{}, fmt.Errorf("failed to convert destination chain: %w", err)
// 	}

// 	destinationAddress, err := utils.CalculateDestinationAddress(event.Payload, &chainInfoBytes)
// 	if err != nil {
// 		return chains.ContractCallWithToken{}, fmt.Errorf("failed to calculate destination address: %w", err)
// 	}

// 	callContract := chains.ContractCall{
// 		EventID:     eventId,
// 		TxHash:      utils.NormalizeHash(event.Raw.TxHash.String()),
// 		BlockNumber: event.Raw.BlockNumber,
// 		LogIndex:    event.Raw.Index,
// 		SourceChain: c.EvmConfig.GetId(),
// 		//3 follows field are used for query to get back payload, so need to convert to lower case
// 		DestinationChain:   event.DestinationChain,
// 		DestinationAddress: utils.NormalizeAddress(destinationAddress, chainInfoBytes.ChainType()),
// 		SourceAddress:      utils.NormalizeAddress(senderAddress, chainInfoBytes.ChainType()),
// 		PayloadHash:        utils.NormalizeHash(hex.EncodeToString(event.PayloadHash[:])),
// 		Payload:            event.Payload,
// 	}
// 	contractCallWithToken := chains.ContractCallWithToken{
// 		ContractCall:         callContract,
// 		TokenContractAddress: utils.NormalizeAddress(event.DestinationContractAddress, chainInfoBytes.ChainType()),
// 		Symbol:               event.Symbol,
// 		Amount:               event.Amount.Uint64(),
// 	}
// 	return contractCallWithToken, nil
// }

// func (c *EvmClient) TokenSentEvent2Model(event *contracts.IScalarGatewayTokenSent) (chains.TokenSent, error) {
// 	normalizedTxHash := utils.NormalizeHash(event.Raw.TxHash.String())
// 	eventId := fmt.Sprintf("%s-%d", normalizedTxHash, event.Raw.Index)
// 	senderAddress := event.Sender.String()
// 	tokenContractAddress, err := c.dbAdapter.GetTokenAddressBySymbol(context.Background(), c.EvmConfig.GetId(), event.Symbol)
// 	if err != nil {
// 		return chains.TokenSent{}, fmt.Errorf("failed to get token address by symbol: %w", err)
// 	}
// 	tokenSent := chains.TokenSent{
// 		EventID:     eventId,
// 		SourceChain: c.EvmConfig.GetId(),
// 		TxHash:      normalizedTxHash,
// 		BlockNumber: event.Raw.BlockNumber,
// 		LogIndex:    event.Raw.Index,
// 		//3 follows field are used for query to get back payload, so need to convert to lower case
// 		SourceAddress:        strings.ToLower(senderAddress),
// 		DestinationChain:     event.DestinationChain,
// 		DestinationAddress:   strings.ToLower(event.DestinationAddress),
// 		Symbol:               event.Symbol,
// 		TokenContractAddress: tokenContractAddress.Hex(),
// 		Amount:               event.Amount.Uint64(),
// 		Status:               chains.TokenSentStatusPending,
// 	}
// 	return tokenSent, nil
// }

// func (c *EvmClient) ContractCallApprovedEvent2Model(event *contracts.IScalarGatewayContractCallApproved) (scalarnet.ContractCallApproved, error) {
// 	txHash := event.Raw.TxHash.String()
// 	eventId := strings.ToLower(fmt.Sprintf("%s-%d-%d", txHash, event.SourceEventIndex, event.Raw.Index))
// 	sourceEventIndex := uint64(0)
// 	if event.SourceEventIndex != nil && event.SourceEventIndex.IsUint64() {
// 		sourceEventIndex = event.SourceEventIndex.Uint64()
// 	}
// 	record := scalarnet.ContractCallApproved{
// 		EventID:          eventId,
// 		SourceChain:      event.SourceChain,
// 		DestinationChain: c.EvmConfig.GetId(),
// 		TxHash:           strings.ToLower(txHash),
// 		CommandID:        hex.EncodeToString(event.CommandId[:]),
// 		Sender:           strings.ToLower(event.SourceAddress),
// 		ContractAddress:  strings.ToLower(event.ContractAddress.String()),
// 		PayloadHash:      strings.ToLower(hex.EncodeToString(event.PayloadHash[:])),
// 		SourceTxHash:     strings.ToLower(hex.EncodeToString(event.SourceTxHash[:])),
// 		SourceEventIndex: sourceEventIndex,
// 	}
// 	return record, nil
// }

// func (c *EvmClient) CommandExecutedEvent2Model(event *contracts.IScalarGatewayExecuted) chains.CommandExecuted {
// 	cmdExecuted := chains.CommandExecuted{
// 		SourceChain: c.EvmConfig.GetId(),
// 		Address:     event.Raw.Address.String(),
// 		TxHash:      strings.ToLower(event.Raw.TxHash.String()),
// 		BlockNumber: uint64(event.Raw.BlockNumber),
// 		LogIndex:    uint(event.Raw.Index),
// 		CommandID:   hex.EncodeToString(event.CommandId[:]),
// 	}
// 	return cmdExecuted
// }

// func (c *EvmClient) TokenDeployedEvent2Model(event *contracts.IScalarGatewayTokenDeployed) chains.TokenDeployed {
// 	tokenDeployed := chains.TokenDeployed{
// 		Chain:        c.EvmConfig.GetId(),
// 		BlockNumber:  uint64(event.Raw.BlockNumber),
// 		TxHash:       event.Raw.TxHash.String(),
// 		Symbol:       event.Symbol,
// 		TokenAddress: event.TokenAddresses.String(),
// 	}
// 	return tokenDeployed
// }

// func (c *EvmClient) SwitchPhaseEvent2Model(event *contracts.IScalarGatewaySwitchPhase) chains.SwitchedPhase {
// 	switchPhase := chains.SwitchedPhase{
// 		Chain:             c.EvmConfig.GetId(),
// 		BlockNumber:       uint64(event.Raw.BlockNumber),
// 		TxHash:            event.Raw.TxHash.String(),
// 		CustodianGroupUid: hex.EncodeToString(event.CustodianGroupId[:]),
// 		SessionSequence:   event.Sequence,
// 		From:              event.From,
// 		To:                event.To,
// 	}
// 	return switchPhase
// }

// import (
// 	"github.com/scalarorg/data-models/chains"
// )

// const (
// 	//CUSTODIAL_NETWORK_NAME        = "Custodial.Network"
// 	//SCALAR_NETWORK_NAME           = "Scalar.Network"
// 	//EVENT_BTC_SIGNATURE_REQUESTED = "Btc.SignatureRequested"
// 	//EVENT_BTC_PSBT_SIGN_REQUEST          = "Btc.PsbtSignRequest"
// 	// EVENT_CUSTODIAL_SIGNATURES_CONFIRMED = "Custodial.SignaturesConfirmed"
// 	// EVENT_ELECTRS_VAULT_TRANSACTION      = "Electrs.VaultTransaction"
// 	// EVENT_ELECTRS_REDEEM_TRANSACTION     = "Electrs.RedeemTransaction"
// 	// EVENT_ELECTRS_NEW_BLOCK              = "Electrs.NewBlock"
// 	// EVENT_SCALAR_TOKEN_SENT              = "Scalar.TokenSent"
// 	// EVENT_SCALAR_DEST_CALL_APPROVED      = "Scalar.ContractCallApproved"
// 	// EVENT_SCALAR_BATCHCOMMAND_SIGNED     = "Scalar.BatchCommandSigned"
// 	// EVENT_SCALAR_COMMAND_EXECUTED        = "Scalar.CommandExecuted"
// 	// EVENT_SCALAR_CREATE_PSBT_REQUEST     = "Scalar.CreatePsbtRequest"
// 	// EVENT_SCALAR_SWITCH_PHASE_STARTED    = "Scalar.StartedSwitchPhase"
// 	// EVENT_SCALAR_REDEEM_TOKEN_APPROVED   = "Scalar.RedeemTokenApproved"
// 	EVENT_EVM_CONTRACT_CALL_APPROVED   = "ContractCallApproved"
// 	EVENT_EVM_CONTRACT_CALL            = "ContractCall"
// 	EVENT_EVM_CONTRACT_CALL_WITH_TOKEN = "ContractCallWithToken"
// 	EVENT_EVM_REDEEM_TOKEN             = "RedeemToken"
// 	EVENT_EVM_TOKEN_SENT               = "TokenSent"
// 	EVENT_EVM_COMMAND_EXECUTED         = "Executed"
// 	EVENT_EVM_TOKEN_DEPLOYED           = "TokenDeployed"
// 	EVENT_EVM_SWITCHED_PHASE           = "SwitchPhase"
// )

// type EventEnvelope struct {
// 	DestinationChain string      // The source chain of the event
// 	EventType        string      // The name of the event in format "ComponentName.EventName"
// 	MessageID        string      // The message id of the event used add RelayData'id
// 	CommandIDs       []string    // The command ids
// 	Data             interface{} // The actual event data
// }
// type SignatureRequest struct {
// 	ExecuteParams *ExecuteParams
// 	Base64Psbt    string
// }
// type ConfirmTxsRequest struct {
// 	ChainName string
// 	TxHashs   map[string]string //Map txHash to DestinationChain, user for validate destination chain
// }

// type ConfirmRedeemTxRequest struct {
// 	Chain    string
// 	GroupUid string
// 	TxHashs  []string
// }

// type RedeemTxEvents struct {
// 	Chain     string
// 	GroupUid  string
// 	Sequence  uint64
// 	RedeemTxs []*chains.RedeemTx
// }
// type ChainBlockHeight struct {
// 	Chain  string
// 	Height uint64
// }

// import (
// 	"context"
// 	"encoding/hex"
// 	"fmt"
// 	"time"

// 	"github.com/ethereum/go-ethereum/accounts/abi/bind"
// 	"github.com/ethereum/go-ethereum/common"
// 	ethtypes "github.com/ethereum/go-ethereum/core/types"
// 	"github.com/rs/zerolog/log"
// 	"github.com/scalarorg/data-models/chains"
// 	contracts "github.com/scalarorg/scalar-healer/pkg/evm/contracts/generated"
// )

// func (ec *EvmClient) HandleContractCallWithToken(ctx context.Context, event *contracts.IScalarGatewayContractCallWithToken) error {
// 	//0. Preprocess the event
// 	ec.preprocessContractCallWithToken(event)
// 	//1. Convert into a RelayData instance then store to the db
// 	contractCallWithToken, err := ec.ContractCallWithToken2Model(event)
// 	if err != nil {
// 		return fmt.Errorf("failed to convert ContractCallEvent to ContractCallWithToken: %w", err)
// 	}
// 	//2. update last checkpoint
// 	lastCheckpoint, err := ec.dbAdapter.GetLastEventCheckPoint(ctx, ec.EvmConfig.GetId(), EVENT_EVM_CONTRACT_CALL_WITH_TOKEN, ec.EvmConfig.StartBlock)
// 	if err != nil {
// 		log.Debug().Str("chainId", ec.EvmConfig.GetId()).
// 			Str("eventName", EVENT_EVM_CONTRACT_CALL_WITH_TOKEN).
// 			Msg("[EvmClient] [handleContractCallWithToken] Get event from begining")
// 	}
// 	if event.Raw.BlockNumber > lastCheckpoint.BlockNumber ||
// 		(event.Raw.BlockNumber == lastCheckpoint.BlockNumber && event.Raw.TxIndex > lastCheckpoint.LogIndex) {
// 		lastCheckpoint.BlockNumber = event.Raw.BlockNumber
// 		lastCheckpoint.TxHash = event.Raw.TxHash.String()
// 		lastCheckpoint.LogIndex = event.Raw.Index
// 		lastCheckpoint.EventKey = fmt.Sprintf("%s-%d-%d", event.Raw.TxHash.String(), event.Raw.BlockNumber, event.Raw.Index)
// 	}
// 	//3. store relay data to the db, update last checkpoint
// 	err = ec.dbAdapter.SaveContractCallWithToken(ctx, &contractCallWithToken, lastCheckpoint)
// 	if err != nil {
// 		return fmt.Errorf("failed to create evm contract call: %w", err)
// 	}
// 	return ec.AddContractCallWithTokenToEvmBatchCommand(event)
// }
// func (ec *EvmClient) AddContractCallWithTokenToEvmBatchCommand(event *contracts.IScalarGatewayContractCallWithToken) error {
// 	//TODO: implement this
// 	return nil
// }
// func (ec *EvmClient) HandleRedeemToken(ctx context.Context, event *contracts.IScalarGatewayRedeemToken) error {
// 	//0. Preprocess the event
// 	log.Info().Str("Chain", ec.EvmConfig.ID).Any("event", event).Msg("[EvmClient] [HandleRedeemToken] Start processing evm redeem token")
// 	//1. Convert into a RelayData instance then store to the db
// 	redeemToken, err := ec.RedeemTokenEvent2Model(event)
// 	if err != nil {
// 		return fmt.Errorf("failed to convert ContractCallEvent to ContractCallWithToken: %w", err)
// 	}
// 	//2. update last checkpoint
// 	lastCheckpoint, err := ec.dbAdapter.GetLastEventCheckPoint(ctx, ec.EvmConfig.GetId(), EVENT_EVM_REDEEM_TOKEN, ec.EvmConfig.StartBlock)
// 	if err != nil {
// 		log.Debug().Str("chainId", ec.EvmConfig.GetId()).
// 			Str("eventName", EVENT_EVM_CONTRACT_CALL_WITH_TOKEN).
// 			Msg("[EvmClient] [handleContractCallWithToken] Get event from begining")
// 	}
// 	if event.Raw.BlockNumber > lastCheckpoint.BlockNumber ||
// 		(event.Raw.BlockNumber == lastCheckpoint.BlockNumber && event.Raw.TxIndex > lastCheckpoint.LogIndex) {
// 		lastCheckpoint.BlockNumber = event.Raw.BlockNumber
// 		lastCheckpoint.TxHash = event.Raw.TxHash.String()
// 		lastCheckpoint.LogIndex = event.Raw.Index
// 		lastCheckpoint.EventKey = fmt.Sprintf("%s-%d-%d", event.Raw.TxHash.String(), event.Raw.BlockNumber, event.Raw.Index)
// 	}
// 	//3. store relay data to the db, update last checkpoint
// 	err = ec.dbAdapter.SaveContractCallWithToken(ctx, &redeemToken, lastCheckpoint)
// 	if err != nil {
// 		return fmt.Errorf("failed to create evm contract call: %w", err)
// 	}
// 	return ec.AddRedeemTokenToRedeemBatchPsbt(event)
// }
// func (ec *EvmClient) AddRedeemTokenToRedeemBatchPsbt(event *contracts.IScalarGatewayRedeemToken) error {
// 	//TODO: implement this
// 	return nil
// }
// func (ec *EvmClient) preprocessContractCallWithToken(event *contracts.IScalarGatewayContractCallWithToken) error {
// 	log.Info().
// 		Str("sender", event.Sender.Hex()).
// 		Str("destinationChain", event.DestinationChain).
// 		Str("destinationContractAddress", event.DestinationContractAddress).
// 		Str("payloadHash", hex.EncodeToString(event.PayloadHash[:])).
// 		Str("Symbol", event.Symbol).
// 		Uint64("Amount", event.Amount.Uint64()).
// 		Str("txHash", event.Raw.TxHash.String()).
// 		Uint("logIndex", event.Raw.Index).
// 		Uint("txIndex", event.Raw.TxIndex).
// 		Str("logData", hex.EncodeToString(event.Raw.Data)).
// 		Msg("[EvmClient] [preprocessContractCallWithToken] Start handle Contract call with token")
// 	//Todo: validate the event
// 	return nil
// }

// func (ec *EvmClient) HandleTokenSent(ctx context.Context, event *contracts.IScalarGatewayTokenSent) error {
// 	//0. Preprocess the event
// 	ec.preprocessTokenSent(event)
// 	//1. Convert into a RelayData instance then store to the db
// 	tokenSent, err := ec.TokenSentEvent2Model(event)
// 	if err != nil {
// 		log.Error().Err(err).Msg("[EvmClient] [HandleTokenSent] failed to convert TokenSentEvent to model data")
// 		return err
// 	}
// 	//For evm, the token sent is verified immediately by the scalarnet
// 	tokenSent.Status = chains.TokenSentStatusVerifying
// 	//2. update last checkpoint
// 	lastCheckpoint, err := ec.dbAdapter.GetLastEventCheckPoint(ctx, ec.EvmConfig.GetId(), EVENT_EVM_TOKEN_SENT, ec.EvmConfig.StartBlock)
// 	if err != nil {
// 		log.Debug().Str("chainId", ec.EvmConfig.GetId()).
// 			Str("eventName", EVENT_EVM_TOKEN_SENT).
// 			Msg("[EvmClient] [handleTokenSent] Get event from begining")
// 	}
// 	if event.Raw.BlockNumber > lastCheckpoint.BlockNumber ||
// 		(event.Raw.BlockNumber == lastCheckpoint.BlockNumber && event.Raw.TxIndex > lastCheckpoint.LogIndex) {
// 		lastCheckpoint.BlockNumber = event.Raw.BlockNumber
// 		lastCheckpoint.TxHash = event.Raw.TxHash.String()
// 		lastCheckpoint.LogIndex = event.Raw.Index
// 		lastCheckpoint.EventKey = fmt.Sprintf("%s-%d-%d", event.Raw.TxHash.String(), event.Raw.BlockNumber, event.Raw.Index)
// 	}
// 	//3. store relay data to the db, update last checkpoint
// 	err = ec.dbAdapter.SaveTokenSent(ctx, &tokenSent, lastCheckpoint)
// 	if err != nil {
// 		return fmt.Errorf("failed to create evm token send: %w", err)
// 	}
// 	return ec.AddTokenSentToEvmBatchCommand(event)
// }

// func (ec *EvmClient) AddTokenSentToEvmBatchCommand(event *contracts.IScalarGatewayTokenSent) error {
// 	//TODO: implement this
// 	return nil
// }
// func (ec *EvmClient) preprocessTokenSent(event *contracts.IScalarGatewayTokenSent) error {
// 	log.Info().
// 		Str("sender", event.Sender.Hex()).
// 		Str("destinationChain", event.DestinationChain).
// 		Str("destinationAddress", event.DestinationAddress).
// 		Str("txHash", event.Raw.TxHash.String()).
// 		Str("symbol", event.Symbol).
// 		Uint64("amount", event.Amount.Uint64()).
// 		Uint("logIndex", event.Raw.Index).
// 		Uint("txIndex", event.Raw.TxIndex).
// 		Str("logData", hex.EncodeToString(event.Raw.Data)).
// 		Msg("[EvmClient] [preprocessTokenSent] Start handle TokenSent")
// 	//Todo: validate the event
// 	return nil
// }

// // func (ec *EvmClient) handleContractCall(event *contracts.IScalarGatewayContractCall) error {
// // 	//0. Preprocess the event
// // 	ec.preprocessContractCall(event)
// // 	//1. Convert into a RelayData instance then store to the db
// // 	contractCall, err := ec.ContractCallEvent2Model(event)
// // 	if err != nil {
// // 		return fmt.Errorf("failed to convert ContractCallEvent to RelayData: %w", err)
// // 	}
// // 	//2. update last checkpoint
// // 	lastCheckpoint, err := ec.dbAdapter.GetLastEventCheckPoint(ec.EvmConfig.GetId(), config.EVENT_EVM_CONTRACT_CALL)
// // 	if err != nil {
// // 		log.Debug().Str("chainId", ec.EvmConfig.GetId()).
// // 			Str("eventName", config.EVENT_EVM_CONTRACT_CALL).
// // 			Msg("[EvmClient] [handleContractCall] Get event from begining")
// // 	}
// // 	if event.Raw.BlockNumber > lastCheckpoint.BlockNumber ||
// // 		(event.Raw.BlockNumber == lastCheckpoint.BlockNumber && event.Raw.TxIndex > lastCheckpoint.LogIndex) {
// // 		lastCheckpoint.BlockNumber = event.Raw.BlockNumber
// // 		lastCheckpoint.TxHash = event.Raw.TxHash.String()
// // 		lastCheckpoint.LogIndex = event.Raw.Index
// // 		lastCheckpoint.EventKey = fmt.Sprintf("%s-%d-%d", event.Raw.TxHash.String(), event.Raw.BlockNumber, event.Raw.Index)
// // 	}
// // 	//3. store relay data to the db, update last checkpoint
// // 	err = ec.dbAdapter.CreateContractCall(contractCall, lastCheckpoint)
// // 	if err != nil {
// // 		return fmt.Errorf("failed to create evm contract call: %w", err)
// // 	}
// // 	//2. Send to the bus
// // 	confirmTxs := events.ConfirmTxsRequest{
// // 		ChainName: ec.EvmConfig.GetId(),
// // 		TxHashs:   map[string]string{contractCall.TxHash: contractCall.DestinationChain},
// // 	}
// // 	if ec.eventBus != nil {
// // 		ec.eventBus.BroadcastEvent(&events.EventEnvelope{
// // 			EventType:        events.EVENT_EVM_CONTRACT_CALL,
// // 			DestinationChain: events.SCALAR_NETWORK_NAME,
// // 			Data:             confirmTxs,
// // 		})
// // 	} else {
// // 		log.Warn().Msg("[EvmClient] [handleContractCall] event bus is undefined")
// // 	}
// // 	return nil
// // }
// // func (ec *EvmClient) preprocessContractCall(event *contracts.IScalarGatewayContractCall) error {
// // 	log.Info().
// // 		Str("sender", event.Sender.Hex()).
// // 		Str("destinationChain", event.DestinationChain).
// // 		Str("destinationContractAddress", event.DestinationContractAddress).
// // 		Str("payloadHash", hex.EncodeToString(event.PayloadHash[:])).
// // 		Str("txHash", event.Raw.TxHash.String()).
// // 		Uint("logIndex", event.Raw.Index).
// // 		Uint("txIndex", event.Raw.TxIndex).
// // 		Str("logData", hex.EncodeToString(event.Raw.Data)).
// // 		Msg("[EvmClient] [preprocessContractCall] Start handle Contract call")
// // 	//Todo: validate the event
// // 	return nil
// // }

// //	func (ec *EvmClient) HandleContractCallApproved(event *contracts.IScalarGatewayContractCallApproved) error {
// //		//0. Preprocess the event
// //		err := ec.preprocessContractCallApproved(event)
// //		if err != nil {
// //			return fmt.Errorf("failed to preprocess contract call approved: %w", err)
// //		}
// //		//1. Convert into a RelayData instance then store to the db
// //		contractCallApproved, err := ec.ContractCallApprovedEvent2Model(event)
// //		if err != nil {
// //			return fmt.Errorf("failed to convert ContractCallApprovedEvent to RelayData: %w", err)
// //		}
// //		err = ec.dbAdapter.SaveSingleValue(&contractCallApproved)
// //		if err != nil {
// //			return fmt.Errorf("failed to create contract call approved: %w", err)
// //		}
// //		// Find relayData from the db by combination (contractAddress, sourceAddress, payloadHash)
// //		// This contract call (initiated by the user call to the source chain) is approved by EVM network
// //		// So anyone can execute it on the EVM by broadcast the corresponding payload to protocol's smart contract on the destination chain
// //		destContractAddress := strings.TrimLeft(event.ContractAddress.Hex(), "0x")
// //		sourceAddress := strings.TrimLeft(event.SourceAddress, "0x")
// //		payloadHash := strings.TrimLeft(hex.EncodeToString(event.PayloadHash[:]), "0x")
// //		relayDatas, err := ec.dbAdapter.FindContractCallByParams(sourceAddress, destContractAddress, payloadHash)
// //		if err != nil {
// //			log.Error().Err(err).Msg("[EvmClient] [handleContractCallApproved] find relay data")
// //			return err
// //		}
// //		log.Debug().Str("contractAddress", event.ContractAddress.String()).
// //			Str("sourceAddress", event.SourceAddress).
// //			Str("payloadHash", hex.EncodeToString(event.PayloadHash[:])).
// //			Any("relayDatas count", len(relayDatas)).
// //			Msg("[EvmClient] [handleContractCallApproved] query relaydata by ContractCall")
// //		//3. Execute payload in the found relaydatas
// //		executeResults, err := ec.executeDestinationCall(event, relayDatas)
// //		if err != nil {
// //			log.Warn().Err(err).Any("executeResults", executeResults).Msg("[EvmClient] [handleContractCallApproved] execute destination call")
// //		}
// //		// Done; Don't need to send to the bus
// //		// TODO: Do we need to update relay data atomically?
// //		err = ec.dbAdapter.UpdateBatchContractCallStatus(executeResults, len(executeResults))
// //		if err != nil {
// //			return fmt.Errorf("failed to update relay data status to executed: %w", err)
// //		}
// //		return nil
// //	}
// // func (ec *EvmClient) executeDestinationCall(event *contracts.IScalarGatewayContractCallApproved, contractCalls []chains.ContractCall) ([]db.ContractCallExecuteResult, error) {
// // 	executeResults := []db.ContractCallExecuteResult{}
// // 	executed, err := ec.isContractCallExecuted(event)
// // 	if err != nil {
// // 		return executeResults, fmt.Errorf("[EvmClient] [executeDestinationCall] failed to check if contract call is approved: %w", err)
// // 	}
// // 	if executed {
// // 		//Update the relay data status to executed
// // 		for _, contractCall := range contractCalls {
// // 			executeResults = append(executeResults, db.ContractCallExecuteResult{
// // 				Status:  chains.ContractCallStatusSuccess,
// // 				EventId: contractCall.EventID,
// // 			})
// // 		}
// // 		return executeResults, fmt.Errorf("destination contract call is already executed")
// // 	}
// // 	if len(contractCalls) > 0 {
// // 		for _, contractCall := range contractCalls {
// // 			if len(contractCall.Payload) == 0 {
// // 				continue
// // 			}
// // 			log.Info().Str("payload", hex.EncodeToString(contractCall.Payload)).
// // 				Msg("[EvmClient] [executeDestinationCall]")
// // 			receipt, err := ec.ExecuteDestinationCall(event.ContractAddress, event.CommandId, event.SourceChain, event.SourceAddress, contractCall.Payload)
// // 			if err != nil {
// // 				return executeResults, fmt.Errorf("execute destination call with error: %w", err)
// // 			}

// // 			log.Info().Any("txReceipt", receipt).Msg("[EvmClient] [executeDestinationCall]")

// // 			if receipt.Hash() != (common.Hash{}) {
// // 				executeResults = append(executeResults, db.ContractCallExecuteResult{
// // 					Status:  chains.ContractCallStatusSuccess,
// // 					EventId: contractCall.EventID,
// // 				})
// // 			} else {
// // 				executeResults = append(executeResults, db.ContractCallExecuteResult{
// // 					Status:  chains.ContractCallStatusFailed,
// // 					EventId: contractCall.EventID,
// // 				})
// // 			}
// // 		}
// // 	}
// // 	return executeResults, nil
// // }

// // Check if ContractCall is already executed
// func (ec *EvmClient) isContractCallExecuted(event *contracts.IScalarGatewayContractCallApproved) (bool, error) {
// 	if ec.transactOpts == nil {
// 		log.Error().
// 			Str("commandId", hex.EncodeToString(event.CommandId[:])).
// 			Str("sourceChain", event.SourceChain).
// 			Str("sourceAddress", event.SourceAddress).
// 			Str("contractAddress", event.ContractAddress.String()).
// 			Msg("[EvmClient] [isContractCallExecuted] auth is nil")
// 		return false, fmt.Errorf("auth is nil")
// 	}
// 	callOpt := &bind.CallOpts{
// 		From:    ec.transactOpts.From,
// 		Context: context.Background(),
// 	}
// 	approved, err := ec.Gateway.IsContractCallApproved(callOpt, event.CommandId, event.SourceChain, event.SourceAddress, event.ContractAddress, event.PayloadHash)
// 	if err != nil {
// 		return false, fmt.Errorf("failed to check if contract call is approved: %w", err)
// 	}
// 	return !approved, nil
// }

// func (ec *EvmClient) preprocessContractCallApproved(event *contracts.IScalarGatewayContractCallApproved) error {
// 	log.Info().Any("event", event).Msgf("[EvmClient] [handleContractCallApproved]")
// 	//Todo: validate the event
// 	return nil
// }

// func (ec *EvmClient) HandleCommandExecuted(ctx context.Context, event *contracts.IScalarGatewayExecuted) error {
// 	//0. Preprocess the event
// 	// ec.preprocessCommandExecuted(event)
// 	// //1. Convert into a RelayData instance then store to the db
// 	// cmdExecuted := ec.CommandExecutedEvent2Model(event)
// 	// err := ec.dbAdapter.UpdateEvmCommandExecuted(ctx, &cmdExecuted)
// 	// if err != nil {
// 	// 	log.Error().Err(err).Msg("[EvmClient] HandleCommandExecuted failed")
// 	// }
// 	return nil
// }

// func (ec *EvmClient) preprocessCommandExecuted(event *contracts.IScalarGatewayExecuted) error {
// 	log.Info().Any("event", event).Msg("[EvmClient] [ExecutedHandler] Start processing evm command executed")
// 	//Todo: validate the event
// 	return nil
// }

// func (ec *EvmClient) SubmitTx(signedTx *ethtypes.Transaction, retryAttempt int) (*ethtypes.Receipt, error) {
// 	if retryAttempt >= ec.EvmConfig.MaxRetry {
// 		return nil, fmt.Errorf("max retry exceeded")
// 	}

// 	// Create a new context with timeout
// 	ctx, cancel := context.WithTimeout(context.Background(), ec.EvmConfig.TxTimeout)
// 	defer cancel()

// 	// Log transaction details
// 	log.Debug().
// 		Interface("tx", signedTx).
// 		Msg("Submitting transaction")

// 	// Send the transaction using the new context
// 	err := ec.Client.SendTransaction(ctx, signedTx)
// 	if err != nil {
// 		log.Error().
// 			Err(err).
// 			Str("rpcUrl", ec.EvmConfig.RPCUrl).
// 			Str("walletAddress", ec.transactOpts.From.String()).
// 			Str("to", signedTx.To().String()).
// 			Str("data", hex.EncodeToString(signedTx.Data())).
// 			Msg("[EvmClient.SubmitTx] Failed to submit transaction")

// 		// Sleep before retry
// 		time.Sleep(ec.EvmConfig.RetryDelay)

// 		log.Debug().
// 			Int("attempt", retryAttempt+1).
// 			Msg("Retrying transaction")

// 		return ec.SubmitTx(signedTx, retryAttempt+1)
// 	}

// 	// Wait for transaction receipt using the new context
// 	receipt, err := bind.WaitMined(ctx, ec.Client, signedTx)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to wait for transaction receipt: %w", err)
// 	}

// 	log.Debug().
// 		Interface("receipt", receipt).
// 		Msg("Transaction receipt received")

// 	return receipt, nil
// }

// func (ec *EvmClient) WaitForTransaction(hash string) (*ethtypes.Receipt, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), ec.EvmConfig.TxTimeout)
// 	defer cancel()

// 	txHash := common.HexToHash(hash)
// 	tx, _, err := ec.Client.TransactionByHash(ctx, txHash)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return bind.WaitMined(ctx, ec.Client, tx)
// }

// import (
// 	"fmt"
// 	"math/big"
// 	"reflect"
// 	"strings"

// 	"github.com/ethereum/go-ethereum/common"
// 	"github.com/ethereum/go-ethereum/core/types"
// 	eth_types "github.com/ethereum/go-ethereum/core/types"
// 	"github.com/rs/zerolog/log"
// 	contracts "github.com/scalarorg/scalar-healer/pkg/evm/contracts/generated"
// )

// type EvmEvent[T any] struct {
// 	Hash             string //TxHash
// 	BlockNumber      uint64
// 	TxIndex          uint
// 	LogIndex         uint
// 	WaitForFinality  func() (*types.Receipt, error)
// 	SourceChain      string
// 	DestinationChain string
// 	EventName        string
// 	Args             T
// }

// type AllEvmEvents struct {
// 	ContractCallApproved *EvmEvent[*contracts.IScalarGatewayContractCallApproved]
// 	ContractCall         *EvmEvent[*contracts.IScalarGatewayContractCall]
// 	Executed             *EvmEvent[*contracts.IScalarGatewayExecuted]
// 	TokenSent            *EvmEvent[*contracts.IScalarGatewayTokenSent]
// }

// type ValidEvmEvent interface {
// 	*contracts.IScalarGatewayContractCallApproved |
// 		*contracts.IScalarGatewayContractCall |
// 		*contracts.IScalarGatewayContractCallWithToken |
// 		*contracts.IScalarGatewayRedeemToken |
// 		*contracts.IScalarGatewayExecuted |
// 		*contracts.IScalarGatewayTokenSent |
// 		*contracts.IScalarGatewaySwitchPhase
// }

// // Todo: Check if this is the correct way to extract the destination chain
// // Maybe add destination chain to the event.Log
// func extractDestChainFromEvmGwContractCallApproved(event *contracts.IScalarGatewayContractCallApproved) string {
// 	return event.SourceChain
// }
// func parseLogIntoEventArgs(log eth_types.Log) (any, error) {
// 	// Try parsing as ContractCallApproved
// 	if eventArgs, err := parseContractCallApproved(log); err == nil {
// 		return eventArgs, nil
// 	}

// 	// Try parsing as ContractCall
// 	if eventArgs, err := parseContractCall(log); err == nil {
// 		return eventArgs, nil
// 	}

// 	// Try parsing as Execute
// 	if eventArgs, err := parseExecute(log); err == nil {
// 		return eventArgs, nil
// 	}

// 	return nil, fmt.Errorf("failed to parse log into any known event type")
// }

// // func parseEventIntoEnvelope(currentChainName string, eventArgs any, log eth_types.Log) (types.EventEnvelope, error) {
// // 	switch args := eventArgs.(type) {
// // 	case *contracts.IScalarGatewayContractCallApproved:
// // 		event, err := parseEventArgsIntoEvent[*contracts.IScalarGatewayContractCallApproved](args, currentChainName, log)
// // 		if err != nil {
// // 			return types.EventEnvelope{}, err
// // 		}
// // 		return types.EventEnvelope{
// // 			Component:        "DbAdapter",
// // 			SenderClientName: currentChainName,
// // 			Handler:          "FindCosmosToEvmCallContractApproved",
// // 			Data:             event,
// // 		}, nil

// // 	case *contracts.IScalarGatewayContractCall:
// // 		event, err := parseEventArgsIntoEvent[*contracts.IScalarGatewayContractCall](args, currentChainName, log)
// // 		if err != nil {
// // 			return types.EventEnvelope{}, err
// // 		}
// // 		return types.EventEnvelope{
// // 			Component:        "DbAdapter",
// // 			SenderClientName: currentChainName,
// // 			Handler:          "CreateEvmCallContractEvent",
// // 			Data:             event,
// // 		}, nil

// // 	case *contracts.IScalarGatewayExecuted:
// // 		event, err := parseEventArgsIntoEvent[*contracts.IScalarGatewayExecuted](args, currentChainName, log)
// // 		if err != nil {
// // 			return types.EventEnvelope{}, err
// // 		}
// // 		return types.EventEnvelope{
// // 			Component:        "DbAdapter",
// // 			SenderClientName: currentChainName,
// // 			Handler:          "CreateEvmExecutedEvent",
// // 			Data:             event,
// // 		}, nil

// // 	default:
// // 		return types.EventEnvelope{}, fmt.Errorf("unknown event type: %T", eventArgs)
// // 	}
// // }

// func parseEventArgsIntoEvent[T ValidEvmEvent](eventArgs T, currentChainName string, log eth_types.Log) (*EvmEvent[T], error) {
// 	// Get the value of eventArgs using reflection
// 	v := reflect.ValueOf(eventArgs).Elem()
// 	sourceChain := currentChainName
// 	if f := v.FieldByName("SourceChain"); f.IsValid() {
// 		sourceChain = f.String()
// 	}
// 	destinationChain := currentChainName
// 	if f := v.FieldByName("DestinationChain"); f.IsValid() {
// 		destinationChain = f.String()
// 	}

// 	return &EvmEvent[T]{
// 		Hash:             log.TxHash.Hex(),
// 		BlockNumber:      log.BlockNumber,
// 		LogIndex:         log.Index,
// 		SourceChain:      sourceChain,
// 		DestinationChain: destinationChain,
// 		Args:             eventArgs,
// 	}, nil
// }

// // parseAnyEvent is a generic function that parses any EVM event into a standardized EvmEvent structure
// func parseEvmEventContractCallApproved[T *contracts.IScalarGatewayContractCallApproved](
// 	currentChainName string,
// 	log eth_types.Log,
// ) (*EvmEvent[T], error) {
// 	eventArgs, err := parseContractCallApproved(log)
// 	if err != nil {
// 		return nil, err
// 	}

// 	event, err := parseEventArgsIntoEvent[T](eventArgs, currentChainName, log)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return event, nil
// }

// func parseContractCallApproved(
// 	log eth_types.Log,
// ) (*contracts.IScalarGatewayContractCallApproved, error) {
// 	event := struct {
// 		CommandId        [32]byte
// 		SourceChain      string
// 		SourceAddress    string
// 		ContractAddress  common.Address
// 		PayloadHash      [32]byte
// 		SourceTxHash     [32]byte
// 		SourceEventIndex *big.Int
// 	}{}

// 	abi, err := GetScalarGatewayAbi()
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to parse ABI: %w", err)
// 	}
// 	if err := abi.UnpackIntoInterface(&event, "ContractCallApproved", log.Data); err != nil {
// 		return nil, fmt.Errorf("failed to unpack event: %w", err)
// 	}

// 	// Add validation for required fields
// 	if len(event.SourceChain) == 0 || !isValidUTF8(event.SourceChain) {
// 		return nil, fmt.Errorf("invalid source chain value")
// 	}

// 	if len(event.SourceAddress) == 0 || !isValidUTF8(event.SourceAddress) {
// 		return nil, fmt.Errorf("invalid source address value")
// 	}

// 	var eventArgs contracts.IScalarGatewayContractCallApproved = contracts.IScalarGatewayContractCallApproved{
// 		CommandId:        event.CommandId,
// 		SourceChain:      event.SourceChain,
// 		SourceAddress:    event.SourceAddress,
// 		ContractAddress:  event.ContractAddress,
// 		PayloadHash:      event.PayloadHash,
// 		SourceTxHash:     event.SourceTxHash,
// 		SourceEventIndex: event.SourceEventIndex,
// 		Raw:              log,
// 	}

// 	fmt.Printf("[EVMListener] [parseContractCallApproved] eventArgs: %+v\n", eventArgs)

// 	return &eventArgs, nil
// }

// func parseEvmEventContractCall[T *contracts.IScalarGatewayContractCall](
// 	currentChainName string,
// 	log eth_types.Log,
// ) (*EvmEvent[T], error) {
// 	eventArgs, err := parseContractCall(log)
// 	if err != nil {
// 		return nil, err
// 	}

// 	event, err := parseEventArgsIntoEvent[T](eventArgs, currentChainName, log)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return event, nil
// }

// func parseContractCall(
// 	log eth_types.Log,
// ) (*contracts.IScalarGatewayContractCall, error) {
// 	event := struct {
// 		Sender                     common.Address
// 		DestinationChain           string
// 		DestinationContractAddress string
// 		PayloadHash                [32]byte
// 		Payload                    []byte
// 	}{}

// 	abi, err := GetScalarGatewayAbi()
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to parse ABI: %w", err)
// 	}
// 	if err := abi.UnpackIntoInterface(&event, "ContractCall", log.Data); err != nil {
// 		return nil, fmt.Errorf("failed to unpack event: %w", err)
// 	}

// 	// Add validation for required fields
// 	if len(event.DestinationChain) == 0 || !isValidUTF8(event.DestinationChain) {
// 		return nil, fmt.Errorf("invalid destination chain value")
// 	}

// 	if len(event.DestinationContractAddress) == 0 || !isValidUTF8(event.DestinationContractAddress) {
// 		return nil, fmt.Errorf("invalid destination address value")
// 	}

// 	var eventArgs contracts.IScalarGatewayContractCall = contracts.IScalarGatewayContractCall{
// 		Sender:                     event.Sender,
// 		DestinationChain:           event.DestinationChain,
// 		DestinationContractAddress: event.DestinationContractAddress,
// 		PayloadHash:                event.PayloadHash,
// 		Payload:                    event.Payload,
// 		Raw:                        log,
// 	}

// 	fmt.Printf("[EVMListener] [parseContractCall] eventArgs: %+v\n", eventArgs)

// 	return &eventArgs, nil
// }

// func parseEvmEventExecute[T *contracts.IScalarGatewayExecuted](
// 	currentChainName string,
// 	log eth_types.Log,
// ) (*EvmEvent[T], error) {
// 	eventArgs, err := parseExecute(log)
// 	if err != nil {
// 		return nil, err
// 	}

// 	event, err := parseEventArgsIntoEvent[T](eventArgs, currentChainName, log)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return event, nil
// }

// func parseExecute(
// 	log eth_types.Log,
// ) (*contracts.IScalarGatewayExecuted, error) {
// 	event := struct {
// 		CommandId [32]byte
// 	}{}
// 	abi, err := GetScalarGatewayAbi()
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to parse ABI: %w", err)
// 	}

// 	// Check if log data size matches exactly what we expect (32 bytes for CommandId)
// 	if len(log.Data) != 32 {
// 		return nil, fmt.Errorf("unexpected log data size: got %d bytes, want 32 bytes", len(log.Data))
// 	}

// 	if err := abi.UnpackIntoInterface(&event, "Executed", log.Data); err != nil {
// 		return nil, fmt.Errorf("failed to unpack event: %w", err)
// 	}

// 	// Add validation for required fields
// 	if len(event.CommandId) == 0 {
// 		return nil, fmt.Errorf("invalid command id value")
// 	}

// 	var eventArgs contracts.IScalarGatewayExecuted = contracts.IScalarGatewayExecuted{
// 		CommandId: event.CommandId,
// 		Raw:       log,
// 	}

// 	fmt.Printf("[EVMListener] [parseExecute] eventArgs: %+v\n", eventArgs)

// 	return &eventArgs, nil
// }

// // Add helper function to validate UTF-8 strings
// func isValidUTF8(s string) bool {
// 	return strings.ToValidUTF8(s, "") == s
// }

// func ParseEventData(receiptLog *eth_types.Log, eventName string, eventData any) error {
// 	gatewayAbi, err := GetScalarGatewayAbi()
// 	if err != nil {
// 		log.Error().Err(err).Msg("[EvmClient] ParseEventData")
// 		return err
// 	}
// 	if gatewayAbi.Events[eventName].ID != receiptLog.Topics[0] {
// 		return fmt.Errorf("receipt log topic 0 does not match %s event id", eventName)
// 	}
// 	// Unpack non-indexed arguments
// 	if err = gatewayAbi.UnpackIntoInterface(eventData, eventName, receiptLog.Data); err != nil {
// 		return fmt.Errorf("failed to unpack event: %w", err)
// 	}
// 	// Unpack indexed arguments
// 	// concat all topic data from second element into single buffer
// 	var buffer []byte
// 	for i := 1; i < len(receiptLog.Topics); i++ {
// 		buffer = append(buffer, receiptLog.Topics[i].Bytes()...)
// 	}
// 	indexedArgs, err := GetEventIndexedArguments(eventName)
// 	if err != nil {
// 		log.Error().Err(err).Msg("[EvmClient] ParseEventData")
// 		return err
// 	}
// 	if len(buffer) > 0 && len(indexedArgs) > 0 {
// 		unpacked, err := indexedArgs.Unpack(buffer)
// 		if err == nil {
// 			indexedArgs.Copy(eventData, unpacked)
// 		}
// 	}
// 	return nil
// }

// import (
// 	"context"
// 	"fmt"
// 	"math/big"

// 	"github.com/ethereum/go-ethereum"
// 	"github.com/ethereum/go-ethereum/accounts/abi"
// 	"github.com/ethereum/go-ethereum/common"
// 	"github.com/ethereum/go-ethereum/core/types"
// 	"github.com/rs/zerolog/log"
// 	"github.com/scalarorg/scalar-healer/pkg/db/sqlc"
// 	contracts "github.com/scalarorg/scalar-healer/pkg/evm/contracts/generated"
// )

// var (
// 	ALL_EVENTS = []string{
// 		EVENT_EVM_CONTRACT_CALL,
// 		EVENT_EVM_CONTRACT_CALL_WITH_TOKEN,
// 		EVENT_EVM_TOKEN_SENT,
// 		EVENT_EVM_CONTRACT_CALL_APPROVED,
// 		EVENT_EVM_COMMAND_EXECUTED,
// 	}
// )

// var abiEventsMap = map[string]abi.Event{}
// var switchPhaseEvent abi.Event
// var redeemTokenEvent abi.Event

// func init() {
// 	for _, event := range scalarGatewayAbi.Events {
// 		abiEventsMap[event.Name] = event
// 	}

// 	switchPhaseEvent = abiEventsMap[EVENT_EVM_SWITCHED_PHASE]
// 	redeemTokenEvent = abiEventsMap[EVENT_EVM_REDEEM_TOKEN]
// }

// func (c *EvmClient) RecoverRedeemSessions(ctx context.Context, groups []common.Hash) (*ChainRedeemSessions, error) {

// 	redeemSessions := NewChainRedeemSessions()

// 	// just get events in range 500 blocks to prevent rate limit
// 	currentBlockNumber, err := c.Client.BlockNumber(ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get current block number: %w", err)
// 	}

// 	baseQuery := ethereum.FilterQuery{
// 		Addresses: []common.Address{c.GatewayAddress},
// 		ToBlock:   big.NewInt(int64(currentBlockNumber)),
// 		FromBlock: big.NewInt(int64(currentBlockNumber - uint64(c.EvmConfig.RecoverRange) + 1)),
// 	}

// 	log.Info().
// 		Str("Number of groups", fmt.Sprintf("%d", len(groups))).
// 		Msg("[EvmClient] [RecoverRedeemSessions] start recovering redeem sessions")

// 	for _, gr := range groups {
// 		isEnoughEvents := false
// 		query := baseQuery
// 		query.Topics = [][]common.Hash{{switchPhaseEvent.ID}, {gr}}
// 		for !isEnoughEvents && query.FromBlock.Uint64() > c.EvmConfig.StartBlock {
// 			log.Info().
// 				Str("Chain", c.EvmConfig.ID).
// 				Str("Group", gr.String()).
// 				Msgf("[EvmClient] [RecoverRedeemSessions] querying logs fromBlock: %d, toBlock: %d", query.FromBlock.Uint64(), query.ToBlock.Uint64())
// 			events, err := c.Client.FilterLogs(ctx, query)
// 			if err != nil {
// 				return nil, fmt.Errorf("failed to filter logs: %w", err)
// 			}

// 			if len(events) > 0 {
// 				for i := len(events) - 1; i >= 0; i-- {
// 					event := events[i]
// 					switchedPhaseEvent, err := parseSwitchPhaseEvent(event)
// 					if err != nil {
// 						log.Error().Err(err).Msgf("[EvmClient] [RecoverRedeemSessions] failed to parse event %s", switchPhaseEvent.Name)
// 						return nil, fmt.Errorf("[EvmClient] [RecoverRedeemSessions]: failed to parse event %s", switchPhaseEvent.Name)
// 					}

// 					counter := redeemSessions.AppendSwitchPhaseEvent(gr.Hex(), switchedPhaseEvent)
// 					if counter == 2 || switchedPhaseEvent.To == sqlc.RedeemPhasePREPARING.Uint8() && switchedPhaseEvent.Sequence == 1 {
// 						log.Debug().
// 							Str("Chain", c.EvmConfig.ID).
// 							Str("Group", gr.String()).
// 							Msgf("[EvmClient] [RecoverRedeemSessions] found %d events", counter)
// 						isEnoughEvents = true
// 						break
// 					}
// 				}
// 			}

// 			query.ToBlock = big.NewInt(int64(query.ToBlock.Uint64() - uint64(c.EvmConfig.RecoverRange)))
// 			query.FromBlock = big.NewInt(int64(query.FromBlock.Uint64() - uint64(c.EvmConfig.RecoverRange)))
// 		}
// 	}

// 	return redeemSessions, nil
// }

// func parseSwitchPhaseEvent(log types.Log) (*contracts.IScalarGatewaySwitchPhase, error) {
// 	eventName := EVENT_EVM_SWITCHED_PHASE
// 	switchedPhase := &contracts.IScalarGatewaySwitchPhase{
// 		Raw: log,
// 	}
// 	err := ParseEventData(&log, eventName, switchedPhase)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to parse event %s: %w", eventName, err)
// 	}
// 	return switchedPhase, nil
// }

// func parseRedeemTokenEvent(log types.Log) (*contracts.IScalarGatewayRedeemToken, error) {
// 	eventName := EVENT_EVM_REDEEM_TOKEN
// 	redeemToken := &contracts.IScalarGatewayRedeemToken{
// 		Raw: log,
// 	}
// 	err := ParseEventData(&log, eventName, redeemToken)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to parse event %s: %w", eventName, err)
// 	}
// 	return redeemToken, nil
// }

// import (
// 	"encoding/hex"
// 	"fmt"
// 	"math"
// 	"strings"
// 	"sync"
// 	"sync/atomic"
// 	"time"

// 	"github.com/ethereum/go-ethereum/accounts/abi"
// 	"github.com/ethereum/go-ethereum/common"
// 	ethTypes "github.com/ethereum/go-ethereum/core/types"
// 	"github.com/rs/zerolog/log"
// 	"github.com/scalarorg/bitcoin-vault/go-utils/encode"
// 	"github.com/scalarorg/bitcoin-vault/go-utils/types"
// 	"github.com/scalarorg/scalar-healer/pkg/db/sqlc"
// 	contracts "github.com/scalarorg/scalar-healer/pkg/evm/contracts/generated"
// )

// type ValidWatchEvent interface {
// 	*contracts.IScalarGatewayTokenSent |
// 		*contracts.IScalarGatewayContractCallWithToken |
// 		*contracts.IScalarGatewayExecuted |
// 		*contracts.IScalarGatewaySwitchPhase |
// 		*contracts.IScalarGatewayRedeemToken
// 	//*contracts.IScalarGatewayContractCall |
// 	//*contracts.IScalarGatewayContractCallApproved |
// 	//*contracts.IScalarGatewayTokenDeployed
// }

// const (
// 	HashLen      = 32
// 	baseDelay    = 5 * time.Second
// 	maxDelay     = 2 * time.Minute
// 	maxAttempts  = math.MaxUint64
// 	jitterFactor = 0.2 // 20% jitter
// )

// type UTXO struct {
// 	TxID         [HashLen]byte
// 	Vout         uint32
// 	ScriptPubkey []byte
// 	AmountInSats uint64
// 	Reserved     map[string]uint64
// }

// type ExecuteParams struct {
// 	SourceChain      string
// 	SourceAddress    string
// 	ContractAddress  common.Address
// 	PayloadHash      [32]byte
// 	SourceTxHash     [32]byte
// 	SourceEventIndex uint64
// }

// const (
// 	COMPONENT_NAME = "EvmClient"
// 	// Initial retry interval
// )

// type DecodedExecuteData struct {
// 	//Data
// 	ChainId    uint64
// 	CommandIds [][32]byte
// 	Commands   []string
// 	Params     [][]byte
// 	//Proof
// 	Operators  []common.Address
// 	Weights    []uint64
// 	Threshold  uint64
// 	Signatures []string
// 	//Input
// 	Input []byte
// }

// type ExecuteData[T any] struct {
// 	//Data
// 	Data T
// 	//Proof
// 	Operators  []common.Address
// 	Weights    []uint64
// 	Threshold  uint64
// 	Signatures []string
// 	//Input
// 	Input []byte
// }
// type ApproveContractCall struct {
// 	ChainId    uint64
// 	CommandIds [][32]byte
// 	Commands   []string
// 	Params     [][]byte
// }
// type DeployToken struct {
// 	//Data
// 	Name         string
// 	Symbol       string
// 	Decimals     uint8
// 	Cap          uint64
// 	TokenAddress common.Address
// 	MintLimit    uint64
// }

// type RedeemPhase struct {
// 	Sequence uint64
// 	Phase    uint8
// }

// type MissingLogs struct {
// 	chainId                     string
// 	mutex                       sync.Mutex
// 	logs                        []ethTypes.Log
// 	lastPreparingSwitchedEvents map[string]*contracts.IScalarGatewaySwitchPhase //Map each custodian group to the last switched phase event
// 	lastExecutingSwitchedEvents map[string]*contracts.IScalarGatewaySwitchPhase //Map each custodian group to the last switched phase event
// 	Recovered                   atomic.Bool                                     //True if logs are recovered
// 	RedeemTxs                   map[string][]string                             //Map each destination chain to the array of redeem txs
// }

// func (m *MissingLogs) IsRecovered() bool {
// 	return m.Recovered.Load()
// }
// func (m *MissingLogs) SetRecovered(recovered bool) {
// 	m.Recovered.Store(recovered)
// }

// func (m *MissingLogs) AppendLogs(logs []ethTypes.Log) {
// 	m.mutex.Lock()
// 	defer m.mutex.Unlock()
// 	log.Info().Str("Chain", m.chainId).Int("Number of event logs", len(logs)).Msgf("[EvmClient] [AppendLogs] appending logs")
// 	redeemTokenEvent, ok := GetEventByName(EVENT_EVM_REDEEM_TOKEN)
// 	if !ok {
// 		log.Error().Msgf("RedeemToken event not found")
// 		return
// 	}
// 	switchedPhaseEvent, ok := GetEventByName(EVENT_EVM_SWITCHED_PHASE)
// 	if !ok {
// 		log.Error().Msgf("SwitchedPhase event not found")
// 		return
// 	}
// 	for _, eventLog := range logs {
// 		if eventLog.Topics[0] == redeemTokenEvent.ID {
// 			m.appendRedeemLog(redeemTokenEvent, eventLog)
// 		} else if eventLog.Topics[0] != switchedPhaseEvent.ID {
// 			m.logs = append(m.logs, eventLog)
// 		}
// 	}
// 	log.Info().Str("Chain", m.chainId).Int("Number of logs", len(m.logs)).Msgf("[EvmClient] [AppendLogs] appended logs")
// }
// func (m *MissingLogs) appendRedeemLog(redeemEvent *abi.Event, eventLog ethTypes.Log) error {
// 	redeemToken := &contracts.IScalarGatewayRedeemToken{
// 		Raw: eventLog,
// 	}
// 	err := ParseEventData(&eventLog, redeemEvent.Name, redeemToken)
// 	if err != nil {
// 		log.Error().Err(err).Msgf("failed to parse event %s", redeemEvent.Name)
// 		return err
// 	}
// 	groupUid := hex.EncodeToString(redeemToken.CustodianGroupUID[:])
// 	if m.isRedeemLogInLastPhase(groupUid, redeemToken) {
// 		log.Info().Str("Chain", m.chainId).Any("eventLog", eventLog).Msgf("[EvmClient] [AppendLogs] appending redeem log")
// 		m.logs = append(m.logs, eventLog)
// 		if m.RedeemTxs == nil {
// 			m.RedeemTxs = make(map[string][]string)
// 		}
// 		txs := m.RedeemTxs[redeemToken.DestinationChain]
// 		m.RedeemTxs[redeemToken.DestinationChain] = append(txs, redeemToken.Raw.TxHash.Hex())
// 	} else {
// 		log.Info().Str("Chain", m.chainId).Any("eventLog", eventLog).Msgf("[EvmClient] [AppendLogs] skipping redeem log due to not in last phase")
// 		return fmt.Errorf("redeem log not in last phase")
// 	}
// 	return nil
// }
// func (m *MissingLogs) isRedeemLogInLastPhase(groupUid string, redeemToken *contracts.IScalarGatewayRedeemToken) bool {
// 	lastSwitchedPhase, ok := m.lastPreparingSwitchedEvents[groupUid]
// 	if !ok {
// 		log.Error().Str("groupUid", groupUid).Msgf("Last switched phase event not found")
// 		return false
// 	}
// 	if redeemToken.Sequence != lastSwitchedPhase.Sequence {
// 		log.Error().Str("groupUid", groupUid).Msgf("Redeem event sequence %d is not equal to last switched phase sequence %d", redeemToken.Sequence, lastSwitchedPhase.Sequence)
// 		return false
// 	}
// 	return true
// }
// func (m *MissingLogs) SetLastSwitchedEvents(mapPreparingEvents map[string]*contracts.IScalarGatewaySwitchPhase,
// 	mapExecutingEvents map[string]*contracts.IScalarGatewaySwitchPhase) {
// 	m.mutex.Lock()
// 	defer m.mutex.Unlock()
// 	m.lastPreparingSwitchedEvents = mapPreparingEvents
// 	m.lastExecutingSwitchedEvents = mapExecutingEvents
// 	_, ok := GetEventByName(EVENT_EVM_SWITCHED_PHASE)
// 	if !ok {
// 		log.Error().Msgf("SwitchedPhase event not found")
// 		return
// 	}
// 	// for groupUid, eventLog := range mapSwitchedPhaseEvents {
// 	// 	params, err := event.Inputs.Unpack(eventLog.Data)
// 	// 	if err != nil {
// 	// 		log.Error().Msgf("Failed to unpack switched phase event: %v", err)
// 	// 		return
// 	// 	}
// 	// 	log.Info().Any("params", params).Msgf("Switched phase event: %v", eventLog)

// 	// 	// fromPhase := params[0].(uint8)
// 	// 	// toPhase := params[1].(uint8)
// 	// 	m.lastSwitchedPhaseEvent[eventLog.Topics[2].String()] = eventLog
// 	// }

// }
// func (m *MissingLogs) GetExecutingEvents() map[string]*contracts.IScalarGatewaySwitchPhase {
// 	m.mutex.Lock()
// 	defer m.mutex.Unlock()
// 	return m.lastExecutingSwitchedEvents
// }
// func (m *MissingLogs) GetLogs(count int) []ethTypes.Log {
// 	m.mutex.Lock()
// 	defer m.mutex.Unlock()
// 	var logs []ethTypes.Log
// 	if len(m.logs) <= count {
// 		logs = m.logs
// 		m.logs = []ethTypes.Log{}
// 	} else {
// 		logs = m.logs[:count]
// 		m.logs = m.logs[count:]
// 	}
// 	log.Info().Str("Chain", m.chainId).Int("Number of logs", len(logs)).
// 		Int("Remaining logs", len(m.logs)).
// 		Msgf("[MissingLogs] [GetLogs] returned logs")
// 	return logs
// }

// // Each chain store switch phase events array with 1 or 2 elements of the form:
// // 1. [Preparing]
// // 2. [Preparing, Executing] in the same sequence
// type ChainRedeemSessions struct {
// 	SwitchPhaseEvents map[string][]*contracts.IScalarGatewaySwitchPhase //Map by custodian group uid
// 	RedeemTokenEvents map[string][]*contracts.IScalarGatewayRedeemToken
// }

// func NewChainRedeemSessions() *ChainRedeemSessions {
// 	return &ChainRedeemSessions{
// 		SwitchPhaseEvents: make(map[string][]*contracts.IScalarGatewaySwitchPhase),
// 		RedeemTokenEvents: make(map[string][]*contracts.IScalarGatewayRedeemToken),
// 	}
// }

// // Return number of events added
// func (s *ChainRedeemSessions) AppendSwitchPhaseEvent(groupUid string, event *contracts.IScalarGatewaySwitchPhase) int {
// 	//Put switch phase event in the first position
// 	switchPhaseEvents, ok := s.SwitchPhaseEvents[groupUid]
// 	if !ok {
// 		s.SwitchPhaseEvents[groupUid] = []*contracts.IScalarGatewaySwitchPhase{event}
// 		return 1
// 	}
// 	if len(switchPhaseEvents) >= 2 {
// 		log.Warn().Str("groupUid", groupUid).Msg("[ChainRedeemSessions] [AppendSwitchPhaseEvent] switch phase events already has 2 elements")
// 		return 2
// 	}
// 	currentPhase := switchPhaseEvents[0]
// 	log.Warn().Str("groupUid", groupUid).Any("current element", currentPhase).
// 		Any("incomming element", event).
// 		Msg("[ChainRedeemSessions] [AppendSwitchPhaseEvent] switch phase event has the same sequence")
// 	if currentPhase.Sequence == event.Sequence {
// 		if currentPhase.To == sqlc.RedeemPhasePREPARING.Uint8() && event.To == sqlc.RedeemPhaseEXECUTING.Uint8() {
// 			s.SwitchPhaseEvents[groupUid] = append(switchPhaseEvents, event)
// 			return 2
// 		} else if currentPhase.To == sqlc.RedeemPhaseEXECUTING.Uint8() && event.To == sqlc.RedeemPhasePREPARING.Uint8() {
// 			s.SwitchPhaseEvents[groupUid] = []*contracts.IScalarGatewaySwitchPhase{event, currentPhase}
// 			return 2
// 		} else {
// 			log.Warn().Msg("[ChainRedeemSessions] [AppendSwitchPhaseEvent] event is already in the list")
// 			return 1
// 		}
// 	} else if event.Sequence < currentPhase.Sequence {
// 		if event.Sequence == currentPhase.Sequence-1 && event.To == sqlc.RedeemPhaseEXECUTING.Uint8() && currentPhase.To == sqlc.RedeemPhasePREPARING.Uint8() {
// 			s.SwitchPhaseEvents[groupUid] = []*contracts.IScalarGatewaySwitchPhase{event, currentPhase}
// 			return 2
// 		}
// 		log.Warn().Msg("[ChainRedeemSessions] [AppendSwitchPhaseEvent] incomming event is too old")
// 		return 1
// 	} else {
// 		//event.Sequence > currentPhase.Sequence
// 		if currentPhase.Sequence == event.Sequence-1 && currentPhase.To == sqlc.RedeemPhaseEXECUTING.Uint8() && event.To == sqlc.RedeemPhasePREPARING.Uint8() {
// 			s.SwitchPhaseEvents[groupUid] = []*contracts.IScalarGatewaySwitchPhase{currentPhase, event}
// 			return 2
// 		}
// 		log.Warn().Msg("[ChainRedeemSessions] [AppendSwitchPhaseEvent] Incomming event is too high, set it as an unique item in the list")
// 		s.SwitchPhaseEvents[groupUid] = []*contracts.IScalarGatewaySwitchPhase{event}
// 		return 1
// 	}
// }

// // Put redeem token event in the first position
// // we keep only the redeem transaction of the max session's sequence
// func (s *ChainRedeemSessions) AppendRedeemTokenEvent(groupUid string, event *contracts.IScalarGatewayRedeemToken) {
// 	redeemEvents, ok := s.RedeemTokenEvents[groupUid]
// 	if !ok {
// 		s.RedeemTokenEvents[groupUid] = []*contracts.IScalarGatewayRedeemToken{event}
// 	} else if len(redeemEvents) > 0 {
// 		lastInsertedEvent := redeemEvents[0]
// 		if event.Sequence > lastInsertedEvent.Sequence {
// 			s.RedeemTokenEvents[groupUid] = []*contracts.IScalarGatewayRedeemToken{event}
// 		} else if lastInsertedEvent.Sequence == event.Sequence {
// 			s.RedeemTokenEvents[groupUid] = append([]*contracts.IScalarGatewayRedeemToken{event}, redeemEvents...)
// 		} else {
// 			// log.Warn().Str("groupUid", groupUid).
// 			// Any("last inserted event", lastInsertedEvent).
// 			// Any("incomming event", event).
// 			// 	Msg("[ChainRedeemSessions] [AppendRedeemTokenEvent] ignore redeem token tx with lower sequence")
// 		}
// 	}
// }
