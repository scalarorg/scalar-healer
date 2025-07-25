// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package sqlc

import (
	"database/sql/driver"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type BatchStatus string

const (
	BatchStatusPENDING  BatchStatus = "PENDING"
	BatchStatusSIGNED   BatchStatus = "SIGNED"
	BatchStatusEXECUTED BatchStatus = "EXECUTED"
)

func (e *BatchStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = BatchStatus(s)
	case string:
		*e = BatchStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for BatchStatus: %T", src)
	}
	return nil
}

type NullBatchStatus struct {
	BatchStatus BatchStatus `json:"batch_status"`
	Valid       bool        `json:"valid"` // Valid is true if BatchStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullBatchStatus) Scan(value interface{}) error {
	if value == nil {
		ns.BatchStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.BatchStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullBatchStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.BatchStatus), nil
}

type CommandType string

const (
	CommandTypeBurnToken                   CommandType = "burnToken"
	CommandTypeDeployToken                 CommandType = "deployToken"
	CommandTypeMintToken                   CommandType = "mintToken"
	CommandTypeApproveContractCall         CommandType = "approveContractCall"
	CommandTypeApproveContractCallWithMint CommandType = "approveContractCallWithMint"
	CommandTypeTransferOperatorship        CommandType = "transferOperatorship"
	CommandTypeSwitchPhase                 CommandType = "switchPhase"
	CommandTypeRegisterCustodianGroup      CommandType = "registerCustodianGroup"
	CommandTypeRedeemToken                 CommandType = "redeemToken"
)

func (e *CommandType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = CommandType(s)
	case string:
		*e = CommandType(s)
	default:
		return fmt.Errorf("unsupported scan type for CommandType: %T", src)
	}
	return nil
}

type NullCommandType struct {
	CommandType CommandType `json:"command_type"`
	Valid       bool        `json:"valid"` // Valid is true if CommandType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullCommandType) Scan(value interface{}) error {
	if value == nil {
		ns.CommandType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.CommandType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullCommandType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.CommandType), nil
}

type RedeemPhase string

const (
	RedeemPhasePREPARING RedeemPhase = "PREPARING"
	RedeemPhaseEXECUTING RedeemPhase = "EXECUTING"
)

func (e *RedeemPhase) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = RedeemPhase(s)
	case string:
		*e = RedeemPhase(s)
	default:
		return fmt.Errorf("unsupported scan type for RedeemPhase: %T", src)
	}
	return nil
}

type NullRedeemPhase struct {
	RedeemPhase RedeemPhase `json:"redeem_phase"`
	Valid       bool        `json:"valid"` // Valid is true if RedeemPhase is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullRedeemPhase) Scan(value interface{}) error {
	if value == nil {
		ns.RedeemPhase, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.RedeemPhase.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullRedeemPhase) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.RedeemPhase), nil
}

type BridgeRequest struct {
	ID        int64            `json:"id"`
	Address   []byte           `json:"address"`
	Signature []byte           `json:"signature"`
	Chain     string           `json:"chain"`
	TxHash    []byte           `json:"tx_hash"`
	Nonce     pgtype.Numeric   `json:"nonce"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
}

type ChainRedeemSession struct {
	ID                int64            `json:"id"`
	Chain             string           `json:"chain"`
	CustodianGroupUid []byte           `json:"custodian_group_uid"`
	Sequence          int64            `json:"sequence"`
	CurrentPhase      RedeemPhase      `json:"current_phase"`
	CreatedAt         pgtype.Timestamp `json:"created_at"`
	UpdatedAt         pgtype.Timestamp `json:"updated_at"`
}

type Command struct {
	ID             []byte           `json:"id"`
	Chain          string           `json:"chain"`
	CommandBatchID []byte           `json:"command_batch_id"`
	Payload        []byte           `json:"payload"`
	Params         []byte           `json:"params"`
	Status         pgtype.Int4      `json:"status"`
	CommandType    CommandType      `json:"command_type"`
	CreatedAt      pgtype.Timestamp `json:"created_at"`
	UpdatedAt      pgtype.Timestamp `json:"updated_at"`
}

type CommandBatch struct {
	ID        []byte           `json:"id"`
	Chain     string           `json:"chain"`
	Data      []byte           `json:"data"`
	SigHash   []byte           `json:"sig_hash"`
	Signature []byte           `json:"signature"`
	Status    BatchStatus      `json:"status"`
	ExtraData []byte           `json:"extra_data"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
}

type CustodianGroup struct {
	ID            int64            `json:"id"`
	Uid           []byte           `json:"uid"`
	Name          string           `json:"name"`
	BitcoinPubkey []byte           `json:"bitcoin_pubkey"`
	Quorum        int64            `json:"quorum"`
	Custodians    []byte           `json:"custodians"`
	CreatedAt     pgtype.Timestamp `json:"created_at"`
	UpdatedAt     pgtype.Timestamp `json:"updated_at"`
}

type GatewayAddress struct {
	ID        int64            `json:"id"`
	Address   []byte           `json:"address"`
	Chain     string           `json:"chain"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
}

type Nonce struct {
	ID        int64            `json:"id"`
	Address   []byte           `json:"address"`
	Nonce     pgtype.Numeric   `json:"nonce"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
}

type Protocol struct {
	ID                 int64            `json:"id"`
	Symbol             string           `json:"symbol"`
	Name               string           `json:"name"`
	BitcoinPubkey      []byte           `json:"bitcoin_pubkey"`
	Decimals           int64            `json:"decimals"`
	Avatar             string           `json:"avatar"`
	CustodianGroupName string           `json:"custodian_group_name"`
	CustodianGroupUid  []byte           `json:"custodian_group_uid"`
	Tag                string           `json:"tag"`
	LiquidityModel     string           `json:"liquidity_model"`
	Capacity           pgtype.Numeric   `json:"capacity"`
	DailyMintLimit     pgtype.Numeric   `json:"daily_mint_limit"`
	CreatedAt          pgtype.Timestamp `json:"created_at"`
	UpdatedAt          pgtype.Timestamp `json:"updated_at"`
}

type RedeemCommand struct {
	ID          []byte           `json:"id"`
	RequestID   int64            `json:"request_id"`
	Chain       string           `json:"chain"`
	Status      BatchStatus      `json:"status"`
	Params      []byte           `json:"params"`
	Data        []byte           `json:"data"`
	SigHash     []byte           `json:"sig_hash"`
	Signature   []byte           `json:"signature"`
	ExecuteData []byte           `json:"execute_data"`
	CreatedAt   pgtype.Timestamp `json:"created_at"`
	UpdatedAt   pgtype.Timestamp `json:"updated_at"`
}

type RedeemRequest struct {
	ID                int64            `json:"id"`
	Address           []byte           `json:"address"`
	SourceChain       string           `json:"source_chain"`
	DestChain         string           `json:"dest_chain"`
	Symbol            string           `json:"symbol"`
	Amount            string           `json:"amount"`
	LockingScript     []byte           `json:"locking_script"`
	CustodianGroupUid []byte           `json:"custodian_group_uid"`
	CreatedAt         pgtype.Timestamp `json:"created_at"`
	UpdatedAt         pgtype.Timestamp `json:"updated_at"`
}

type RedeemSession struct {
	ID                int64            `json:"id"`
	CustodianGroupUid []byte           `json:"custodian_group_uid"`
	Sequence          int64            `json:"sequence"`
	CurrentPhase      RedeemPhase      `json:"current_phase"`
	LastRedeemTx      []byte           `json:"last_redeem_tx"`
	IsSwitching       pgtype.Bool      `json:"is_switching"`
	PhaseExpiredAt    pgtype.Timestamp `json:"phase_expired_at"`
	CreatedAt         pgtype.Timestamp `json:"created_at"`
	UpdatedAt         pgtype.Timestamp `json:"updated_at"`
}

type Reservation struct {
	RequestID []byte           `json:"request_id"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
}

type Token struct {
	ID        int64            `json:"id"`
	Symbol    string           `json:"symbol"`
	Chain     string           `json:"chain"`
	ChainID   pgtype.Numeric   `json:"chain_id"`
	Active    bool             `json:"active"`
	Address   []byte           `json:"address"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
}

type TransferRequest struct {
	ID                 int64            `json:"id"`
	Address            []byte           `json:"address"`
	Signature          []byte           `json:"signature"`
	Chain              string           `json:"chain"`
	DestinationChain   string           `json:"destination_chain"`
	DestinationAddress []byte           `json:"destination_address"`
	Symbol             string           `json:"symbol"`
	Amount             string           `json:"amount"`
	Nonce              pgtype.Numeric   `json:"nonce"`
	CreatedAt          pgtype.Timestamp `json:"created_at"`
	UpdatedAt          pgtype.Timestamp `json:"updated_at"`
}

type Utxo struct {
	ID                int64            `json:"id"`
	TxID              []byte           `json:"tx_id"`
	Vout              int64            `json:"vout"`
	ScriptPubkey      []byte           `json:"script_pubkey"`
	AmountInSats      pgtype.Numeric   `json:"amount_in_sats"`
	CustodianGroupUid []byte           `json:"custodian_group_uid"`
	BlockHeight       int64            `json:"block_height"`
	CreatedAt         pgtype.Timestamp `json:"created_at"`
	UpdatedAt         pgtype.Timestamp `json:"updated_at"`
}

type UtxoReservation struct {
	UtxoTxID      []byte           `json:"utxo_tx_id"`
	UtxoVout      int64            `json:"utxo_vout"`
	ReservationID []byte           `json:"reservation_id"`
	Amount        pgtype.Numeric   `json:"amount"`
	CreatedAt     pgtype.Timestamp `json:"created_at"`
	UpdatedAt     pgtype.Timestamp `json:"updated_at"`
}
