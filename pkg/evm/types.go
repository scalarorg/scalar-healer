package evm

import (
	"encoding/hex"
	"fmt"
	"math"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	"github.com/scalarorg/bitcoin-vault/go-utils/encode"
	"github.com/scalarorg/bitcoin-vault/go-utils/types"
	"github.com/scalarorg/scalar-healer/pkg/db"
	contracts "github.com/scalarorg/scalar-healer/pkg/evm/contracts/generated"
)

type ValidWatchEvent interface {
	*contracts.IScalarGatewayTokenSent |
		*contracts.IScalarGatewayContractCallWithToken |
		*contracts.IScalarGatewayExecuted |
		*contracts.IScalarGatewaySwitchPhase |
		*contracts.IScalarGatewayRedeemToken
	//*contracts.IScalarGatewayContractCall |
	//*contracts.IScalarGatewayContractCallApproved |
	//*contracts.IScalarGatewayTokenDeployed
}

const (
	HashLen      = 32
	baseDelay    = 5 * time.Second
	maxDelay     = 2 * time.Minute
	maxAttempts  = math.MaxUint64
	jitterFactor = 0.2 // 20% jitter
)

type UTXO struct {
	TxID         [HashLen]byte
	Vout         uint32
	ScriptPubkey []byte
	AmountInSats uint64
	Reserved     map[string]uint64
}

type ExecuteParams struct {
	SourceChain      string
	SourceAddress    string
	ContractAddress  common.Address
	PayloadHash      [32]byte
	SourceTxHash     [32]byte
	SourceEventIndex uint64
}

const (
	COMPONENT_NAME = "EvmClient"
	RETRY_INTERVAL = time.Second * 12 // Initial retry interval
)

type Byte32 [32]uint8
type Bytes []byte
type EvmNetworkConfig struct {
	ChainID      uint64        `mapstructure:"chain_id"`
	ID           string        `mapstructure:"id"`
	Name         string        `mapstructure:"name"`
	RPCUrl       string        `mapstructure:"rpc_url"`
	AuthWeighted string        `mapstructure:"auth_weighted"`
	Gateway      string        `mapstructure:"gateway"`
	Finality     int           `mapstructure:"finality"`
	StartBlock   uint64        `mapstructure:"start_block"`
	PrivateKey   string        `mapstructure:"private_key"`
	GasLimit     uint64        `mapstructure:"gas_limit"`
	BlockTime    time.Duration `mapstructure:"blockTime"` //Timeout im ms for pending txs
	MaxRetry     int
	RecoverRange uint64 `mapstructure:"recover_range"` //Max block range to recover events in single query
	RetryDelay   time.Duration
	TxTimeout    time.Duration `mapstructure:"tx_timeout"` //Timeout for send txs (~3s)
}

func (c *EvmNetworkConfig) GetChainId() uint64 {
	return c.ChainID
}
func (c *EvmNetworkConfig) GetId() string {
	return c.ID
}
func (c *EvmNetworkConfig) GetName() string {
	return c.Name
}
func (c *EvmNetworkConfig) GetFamily() string {
	return types.ChainTypeEVM.String()
}

type DecodedExecuteData struct {
	//Data
	ChainId    uint64
	CommandIds [][32]byte
	Commands   []string
	Params     [][]byte
	//Proof
	Operators  []common.Address
	Weights    []uint64
	Threshold  uint64
	Signatures []string
	//Input
	Input []byte
}

type ExecuteData[T any] struct {
	//Data
	Data T
	//Proof
	Operators  []common.Address
	Weights    []uint64
	Threshold  uint64
	Signatures []string
	//Input
	Input []byte
}
type ApproveContractCall struct {
	ChainId    uint64
	CommandIds [][32]byte
	Commands   []string
	Params     [][]byte
}
type DeployToken struct {
	//Data
	Name         string
	Symbol       string
	Decimals     uint8
	Cap          uint64
	TokenAddress common.Address
	MintLimit    uint64
}

type RedeemPhase struct {
	Sequence uint64
	Phase    uint8
}

type MissingLogs struct {
	chainId                     string
	mutex                       sync.Mutex
	logs                        []ethTypes.Log
	lastPreparingSwitchedEvents map[string]*contracts.IScalarGatewaySwitchPhase //Map each custodian group to the last switched phase event
	lastExecutingSwitchedEvents map[string]*contracts.IScalarGatewaySwitchPhase //Map each custodian group to the last switched phase event
	Recovered                   atomic.Bool                                     //True if logs are recovered
	RedeemTxs                   map[string][]string                             //Map each destination chain to the array of redeem txs
}

func (m *MissingLogs) IsRecovered() bool {
	return m.Recovered.Load()
}
func (m *MissingLogs) SetRecovered(recovered bool) {
	m.Recovered.Store(recovered)
}

func (m *MissingLogs) AppendLogs(logs []ethTypes.Log) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	log.Info().Str("Chain", m.chainId).Int("Number of event logs", len(logs)).Msgf("[EvmClient] [AppendLogs] appending logs")
	redeemTokenEvent, ok := GetEventByName(EVENT_EVM_REDEEM_TOKEN)
	if !ok {
		log.Error().Msgf("RedeemToken event not found")
		return
	}
	switchedPhaseEvent, ok := GetEventByName(EVENT_EVM_SWITCHED_PHASE)
	if !ok {
		log.Error().Msgf("SwitchedPhase event not found")
		return
	}
	for _, eventLog := range logs {
		if eventLog.Topics[0] == redeemTokenEvent.ID {
			m.appendRedeemLog(redeemTokenEvent, eventLog)
		} else if eventLog.Topics[0] != switchedPhaseEvent.ID {
			m.logs = append(m.logs, eventLog)
		}
	}
	log.Info().Str("Chain", m.chainId).Int("Number of logs", len(m.logs)).Msgf("[EvmClient] [AppendLogs] appended logs")
}
func (m *MissingLogs) appendRedeemLog(redeemEvent *abi.Event, eventLog ethTypes.Log) error {
	redeemToken := &contracts.IScalarGatewayRedeemToken{
		Raw: eventLog,
	}
	err := ParseEventData(&eventLog, redeemEvent.Name, redeemToken)
	if err != nil {
		log.Error().Err(err).Msgf("failed to parse event %s", redeemEvent.Name)
		return err
	}
	groupUid := hex.EncodeToString(redeemToken.CustodianGroupUID[:])
	if m.isRedeemLogInLastPhase(groupUid, redeemToken) {
		log.Info().Str("Chain", m.chainId).Any("eventLog", eventLog).Msgf("[EvmClient] [AppendLogs] appending redeem log")
		m.logs = append(m.logs, eventLog)
		if m.RedeemTxs == nil {
			m.RedeemTxs = make(map[string][]string)
		}
		txs := m.RedeemTxs[redeemToken.DestinationChain]
		m.RedeemTxs[redeemToken.DestinationChain] = append(txs, redeemToken.Raw.TxHash.Hex())
	} else {
		log.Info().Str("Chain", m.chainId).Any("eventLog", eventLog).Msgf("[EvmClient] [AppendLogs] skipping redeem log due to not in last phase")
		return fmt.Errorf("redeem log not in last phase")
	}
	return nil
}
func (m *MissingLogs) isRedeemLogInLastPhase(groupUid string, redeemToken *contracts.IScalarGatewayRedeemToken) bool {
	lastSwitchedPhase, ok := m.lastPreparingSwitchedEvents[groupUid]
	if !ok {
		log.Error().Str("groupUid", groupUid).Msgf("Last switched phase event not found")
		return false
	}
	if redeemToken.Sequence != lastSwitchedPhase.Sequence {
		log.Error().Str("groupUid", groupUid).Msgf("Redeem event sequence %d is not equal to last switched phase sequence %d", redeemToken.Sequence, lastSwitchedPhase.Sequence)
		return false
	}
	return true
}
func (m *MissingLogs) SetLastSwitchedEvents(mapPreparingEvents map[string]*contracts.IScalarGatewaySwitchPhase,
	mapExecutingEvents map[string]*contracts.IScalarGatewaySwitchPhase) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.lastPreparingSwitchedEvents = mapPreparingEvents
	m.lastExecutingSwitchedEvents = mapExecutingEvents
	_, ok := GetEventByName(EVENT_EVM_SWITCHED_PHASE)
	if !ok {
		log.Error().Msgf("SwitchedPhase event not found")
		return
	}
	// for groupUid, eventLog := range mapSwitchedPhaseEvents {
	// 	params, err := event.Inputs.Unpack(eventLog.Data)
	// 	if err != nil {
	// 		log.Error().Msgf("Failed to unpack switched phase event: %v", err)
	// 		return
	// 	}
	// 	log.Info().Any("params", params).Msgf("Switched phase event: %v", eventLog)

	// 	// fromPhase := params[0].(uint8)
	// 	// toPhase := params[1].(uint8)
	// 	m.lastSwitchedPhaseEvent[eventLog.Topics[2].String()] = eventLog
	// }

}
func (m *MissingLogs) GetExecutingEvents() map[string]*contracts.IScalarGatewaySwitchPhase {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.lastExecutingSwitchedEvents
}
func (m *MissingLogs) GetLogs(count int) []ethTypes.Log {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	var logs []ethTypes.Log
	if len(m.logs) <= count {
		logs = m.logs
		m.logs = []ethTypes.Log{}
	} else {
		logs = m.logs[:count]
		m.logs = m.logs[count:]
	}
	log.Info().Str("Chain", m.chainId).Int("Number of logs", len(logs)).
		Int("Remaining logs", len(m.logs)).
		Msgf("[MissingLogs] [GetLogs] returned logs")
	return logs
}

// Each chain store switch phase events array with 1 or 2 elements of the form:
// 1. [Preparing]
// 2. [Preparing, Executing] in the same sequence
type ChainRedeemSessions struct {
	SwitchPhaseEvents map[string][]*contracts.IScalarGatewaySwitchPhase //Map by custodian group uid
	RedeemTokenEvents map[string][]*contracts.IScalarGatewayRedeemToken
}

// Return number of events added
func (s *ChainRedeemSessions) AppendSwitchPhaseEvent(groupUid string, event *contracts.IScalarGatewaySwitchPhase) int {
	//Put switch phase event in the first position
	switchPhaseEvents, ok := s.SwitchPhaseEvents[groupUid]
	if !ok {
		s.SwitchPhaseEvents[groupUid] = []*contracts.IScalarGatewaySwitchPhase{event}
		return 1
	}
	if len(switchPhaseEvents) >= 2 {
		log.Warn().Str("groupUid", groupUid).Msg("[ChainRedeemSessions] [AppendSwitchPhaseEvent] switch phase events already has 2 elements")
		return 2
	}
	currentPhase := switchPhaseEvents[0]
	log.Warn().Str("groupUid", groupUid).Any("current element", currentPhase).
		Any("incomming element", event).
		Msg("[ChainRedeemSessions] [AppendSwitchPhaseEvent] switch phase event has the same sequence")
	if currentPhase.Sequence == event.Sequence {
		if currentPhase.To == uint8(db.Preparing) && event.To == uint8(db.Executing) {
			s.SwitchPhaseEvents[groupUid] = append(switchPhaseEvents, event)
			return 2
		} else if currentPhase.To == uint8(db.Executing) && event.To == uint8(db.Preparing) {
			s.SwitchPhaseEvents[groupUid] = []*contracts.IScalarGatewaySwitchPhase{event, currentPhase}
			return 2
		} else {
			log.Warn().Msg("[ChainRedeemSessions] [AppendSwitchPhaseEvent] event is already in the list")
			return 1
		}
	} else if event.Sequence < currentPhase.Sequence {
		if event.Sequence == currentPhase.Sequence-1 && event.To == uint8(db.Executing) && currentPhase.To == uint8(db.Preparing) {
			s.SwitchPhaseEvents[groupUid] = []*contracts.IScalarGatewaySwitchPhase{event, currentPhase}
			return 2
		}
		log.Warn().Msg("[ChainRedeemSessions] [AppendSwitchPhaseEvent] incomming event is too old")
		return 1
	} else {
		//event.Sequence > currentPhase.Sequence
		if currentPhase.Sequence == event.Sequence-1 && currentPhase.To == uint8(db.Executing) && event.To == uint8(db.Preparing) {
			s.SwitchPhaseEvents[groupUid] = []*contracts.IScalarGatewaySwitchPhase{currentPhase, event}
			return 2
		}
		log.Warn().Msg("[ChainRedeemSessions] [AppendSwitchPhaseEvent] Incomming event is too high, set it as an unique item in the list")
		s.SwitchPhaseEvents[groupUid] = []*contracts.IScalarGatewaySwitchPhase{event}
		return 1
	}
}

// Put redeem token event in the first position
// we keep only the redeem transaction of the max session's sequence
func (s *ChainRedeemSessions) AppendRedeemTokenEvent(groupUid string, event *contracts.IScalarGatewayRedeemToken) {
	redeemEvents, ok := s.RedeemTokenEvents[groupUid]
	if !ok {
		s.RedeemTokenEvents[groupUid] = []*contracts.IScalarGatewayRedeemToken{event}
	} else if len(redeemEvents) > 0 {
		lastInsertedEvent := redeemEvents[0]
		if event.Sequence > lastInsertedEvent.Sequence {
			s.RedeemTokenEvents[groupUid] = []*contracts.IScalarGatewayRedeemToken{event}
		} else if lastInsertedEvent.Sequence == event.Sequence {
			s.RedeemTokenEvents[groupUid] = append([]*contracts.IScalarGatewayRedeemToken{event}, redeemEvents...)
		} else {
			log.Warn().Str("groupUid", groupUid).Any("last inserted event", lastInsertedEvent).
				Any("incomming event", event).
				Msg("[ChainRedeemSessions] [AppendRedeemTokenEvent] ignore redeem token tx with lower sequence")
		}
	}
}

type RedeemTokenPayload struct {
	Amount        uint64
	LockingScript []byte
	Utxos         []*UTXO
	RequestId     [32]byte
}

type RedeemTokenPayloadWithType struct {
	RedeemTokenPayload
	PayloadType encode.ContractCallWithTokenPayloadType
}

func (p *RedeemTokenPayloadWithType) AbiPack() ([]byte, error) {
	payload, err := p.RedeemTokenPayload.AbiPack()
	if err != nil {
		return nil, err
	}
	return encode.AppendPayload(encode.ContractCallWithTokenPayloadType(p.PayloadType), payload), nil
}

func (p *RedeemTokenPayloadWithType) AbiUnpack(data []byte) error {
	p.PayloadType = encode.ContractCallWithTokenPayloadType(data[0])
	var payload RedeemTokenPayload
	err := payload.AbiUnpack(data[1:])
	if err != nil {
		return err
	}
	p.RedeemTokenPayload = payload
	return nil
}

func (p *RedeemTokenPayload) AbiPack() ([]byte, error) {
	txIds := make([]string, len(p.Utxos))
	vouts := make([]uint32, len(p.Utxos))
	amounts := make([]uint64, len(p.Utxos))
	for i, utxo := range p.Utxos {
		txIds[i] = hex.EncodeToString(utxo.TxID[:])
		vouts[i] = utxo.Vout
		amounts[i] = utxo.AmountInSats
	}
	return RedeemTokenPayloadArguments.Pack(
		p.Amount,
		p.LockingScript,
		txIds,
		vouts,
		amounts,
		p.RequestId,
	)
}

func (p *RedeemTokenPayload) AbiUnpack(data []byte) error {
	unpacked, err := RedeemTokenPayloadArguments.Unpack(data)
	if err != nil {
		log.Error().Err(err).Msg("redeem token payload abi unpack error")
		return err
	}
	p.Amount = unpacked[0].(uint64)
	p.LockingScript = unpacked[1].([]byte)
	txIds := unpacked[2].([]string)
	vouts := unpacked[3].([]uint32)
	amounts := unpacked[4].([]uint64)
	p.Utxos = make([]*UTXO, len(txIds))
	for i, txId := range txIds {
		hash, err := hex.DecodeString(strings.TrimPrefix(txId, "0x"))
		if err != nil {
			log.Error().Err(err).Msg("txId hash error")
			return err
		}
		var txId [HashLen]byte
		copy(txId[:], hash)
		p.Utxos[i] = &UTXO{TxID: txId, Vout: vouts[i], AmountInSats: amounts[i]}
	}
	p.RequestId = unpacked[5].([32]byte)
	return nil
}
