package sqlc

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
	"github.com/scalarorg/scalar-healer/pkg/utils/chains"
)

func (e RedeemPhase) Uint8() uint8 {
	switch e {
	case RedeemPhasePREPARING:
		return 0
	case RedeemPhaseEXECUTING:
		return 1
	default:
		return math.MaxUint8
	}
}

func (e RedeemPhase) Bytes() byte {
	return byte(e.Uint8())
}

func PhaseFromUint8(phase uint8) RedeemPhase {
	switch phase {
	case 0:
		return RedeemPhasePREPARING
	case 1:
		return RedeemPhaseEXECUTING
	default:
		panic("invalid phase")
	}
}

func (s *ChainRedeemSession) Cmp(other *ChainRedeemSession) int64 {
	if other == nil {
		return math.MaxInt64
	}
	diffSeq := s.Sequence - other.Sequence

	if diffSeq != 0 {
		return diffSeq
	}

	return int64(s.CurrentPhase.Uint8()) - int64(other.CurrentPhase.Uint8())
}

func (c CommandType) String() string {
	return string(c)
}

type CommandStatus uint8

const (
	COMMAND_STATUS_PENDING  CommandStatus = 0
	COMMAND_STATUS_QUEUED   CommandStatus = 1
	COMMAND_STATUS_EXECUTED CommandStatus = 2
)

func (s CommandStatus) Int32() int32 {
	return int32(s)
}

func (c CommandStatus) ToPgType() pgtype.Int4 {
	return pgtype.Int4{
		Int32: c.Int32(),
		Valid: true,
	}
}

type RedeemRequestWithCommand struct {
	*RedeemRequest
	Status      BatchStatus `json:"status"`
	ExecuteData []byte      `json:"execute_data"`
}

type ChainRedeemSessionUpdate struct {
	Chain             string
	CustodianGroupUid []byte
	Sequence          int64
	CurrentPhase      RedeemPhase
	NewPhase          RedeemPhase
}

const commandIDSize = 32

type CommandID [commandIDSize]byte

func NewCommandID(data []byte, chainID string) CommandID {
	var commandID CommandID
	copy(commandID[:], crypto.Keccak256(append(data, []byte(chainID)...))[:commandIDSize])

	return commandID
}

func (c CommandID) Bytes() []byte {
	return c[:]
}

func (c *CommandBatch) GetExtraData() ([][]byte, error) {
	var extraData [][]byte
	err := json.Unmarshal(c.ExtraData, &extraData)
	if err != nil {
		return nil, err
	}
	return extraData, nil
}

type TokenDetails struct {
	Address []byte `json:"address"`
	ChainID int64  `json:"chain_id"`
	Chain   string `json:"chain"`
}

type Custodian struct {
	Name          string `json:"name"`
	ValAddress    string `json:"val_address"`
	BitcoinPubkey []byte `json:"bitcoin_pubkey"`
}

type Custodians []Custodian

func (c *Custodians) FromJson(data []byte) error {
	return json.Unmarshal(data, c)
}

type ProtocolWithTokenDetails struct {
	*Protocol
	Tokens          []TokenDetails `json:"token_details"`
	Custodians      []Custodian    `json:"custodians"`
	CustodianQuorum int64          `json:"custodian_quorum"`
}

func (p *ProtocolWithTokenDetails) GetTokenDetailsByChain(chain string) (TokenDetails, error) {
	for _, token := range p.Tokens {
		if token.Chain == chain {
			return token, nil
		}
	}
	return TokenDetails{}, errors.New("token not found")
}

type UtxoWithReservations struct {
	*Utxo
	Reservations []*UtxoReservation `json:"reservations"`
}

func (utxo *UtxoWithReservations) IsReserved(requestID []byte) bool {
	for _, reserved := range utxo.Reservations {
		if bytes.Equal(reserved.ReservationID, requestID) {
			return true
		}
	}
	return false
}

func (utxo *UtxoWithReservations) AvailableAmount() uint64 {
	reservedAmount := utxo.GetReservedAmount()
	result := ConvertNumericToUint64(utxo.AmountInSats) - reservedAmount
	return result
}

func (utxo *UtxoWithReservations) GetReservedAmount() uint64 {
	amount := uint64(0)
	for _, reserved := range utxo.Reservations {
		if !reserved.Amount.Valid {
			continue
		}
		amount += ConvertNumericToUint64(reserved.Amount)
	}
	return amount
}

type UtxoSnapshot []*UtxoWithReservations

func (snapshot UtxoSnapshot) ReserveUtxos(requestID []byte, amount uint64, quorum uint64, vSizeLimit uint64) ([]Utxo, error) {
	currentInputs, currentOutputs := snapshot.CountInputOutput()
	newInput := 0
	newOutput := 1
	remainingAmount := amount
	reserveUtxos := make([]Utxo, 0)
	mapNewResevations := map[int]uint64{}
	for ind, utxo := range snapshot {
		if utxo.IsReserved(requestID) {
			return nil, fmt.Errorf("requestID %s is already reserved in utxo %x", requestID, utxo.TxID)
		}
		availableAmount := utxo.AvailableAmount()

		log.Info().Msgf("availableAmount %d, remainingAmount %d", availableAmount, remainingAmount)

		if availableAmount > 0 {
			//Reserve amount is min(availableAmount, remainingAmount)
			reserveAmount := availableAmount
			if reserveAmount > remainingAmount {
				reserveAmount = remainingAmount
			}
			remainingAmount -= reserveAmount
			mapNewResevations[ind] = reserveAmount
			reserveUtxos = append(reserveUtxos, Utxo{
				TxID:         utxo.TxID,
				Vout:         utxo.Vout,
				AmountInSats: ConvertUint64ToNumeric(reserveAmount),
				ScriptPubkey: utxo.ScriptPubkey,
			})
			//First reservation
			if len(utxo.Reservations) == 0 {
				newInput += 1
			}
		}
		if remainingAmount == 0 {
			break
		}
	}
	if remainingAmount > 0 {
		return nil, fmt.Errorf("not enough utxos to reserve, remainingAmount %d", remainingAmount)
	}

	//Add extra input and output for collect change amount
	newVsize := chains.CalculateVsize(currentInputs+newInput+1, currentOutputs+newOutput+1, quorum)
	if newVsize > vSizeLimit {
		return nil, fmt.Errorf("new virtual size exceeds the limit %d > %d", newVsize, vSizeLimit)
	}

	for ind, reserveAmount := range mapNewResevations {
		err := snapshot[ind].AppendReserved(requestID, reserveAmount)
		if err != nil {
			return nil, err
		}
	}

	return reserveUtxos, nil
}

func (utxo *UtxoWithReservations) AppendReserved(requestID []byte, amount uint64) error {
	if utxo.Reservations == nil {
		utxo.Reservations = []*UtxoReservation{}
	}
	totalReserved := uint64(0)
	for _, reserved := range utxo.Reservations {
		if bytes.Equal(reserved.ReservationID, requestID) {
			return fmt.Errorf("requestID already reserved in this utxo %x", utxo.TxID)
		}
		totalReserved += ConvertNumericToUint64(reserved.Amount)
	}
	if totalReserved+amount > ConvertNumericToUint64(utxo.AmountInSats) {
		return fmt.Errorf("amount exceeds utxo amount, totalReserved %d, amount %d, utxo.AmountInSats %d", totalReserved, amount, ConvertNumericToUint64(utxo.AmountInSats))
	}
	utxo.Reservations = append(utxo.Reservations, &UtxoReservation{
		ReservationID: requestID,
		Amount:        ConvertUint64ToNumeric(amount),
		UtxoTxID:      utxo.TxID,
		UtxoVout:      utxo.Vout,
	})
	return nil
}

func (s UtxoSnapshot) CountInputOutput() (int, int) {
	inputs := 0
	mapRequests := map[string]bool{}
	for _, utxo := range s {
		if len(utxo.Reservations) > 0 {
			inputs += 1
			for _, reservation := range utxo.Reservations {
				mapRequests[hex.EncodeToString(reservation.ReservationID)] = true
			}
		}
	}
	return inputs, len(mapRequests)
}
