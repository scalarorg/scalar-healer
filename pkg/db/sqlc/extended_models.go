package sqlc

import (
	"encoding/json"
	"errors"
	"math"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/jackc/pgx/v5/pgtype"
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

type CommandBatchStatus uint8

const (
	COMMAND_BATCH_STATUS_PENDING  CommandBatchStatus = 0
	COMMAND_BATCH_STATUS_EXECUTED CommandBatchStatus = 1
)

func (c CommandBatchStatus) ToPgType() pgtype.Int4 {
	return pgtype.Int4{
		Int32: c.Int32(),
		Valid: true,
	}
}

func (s CommandBatchStatus) Int32() int32 {
	return int32(s)
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

// func CreateRedeemParams(ctx context.Context, vsizeLimit uint64, amount *big.Int, sequence uint64, utxos []sqlc.Utxo) ([]byte, *CommandID, error) {

// 	now := time.Now().UnixNano()

// 	bz := make([]byte, 8)
// 	binary.BigEndian.PutUint64(bz, uint64(now))

// 	dataHash := crypto.Keccak256(bz, r.Address, []byte(r.SourceChain), []byte(r.DestChain), []byte(r.Symbol), amount.Bytes())

// 	dataHex := hex.EncodeToString(dataHash)

// 	// reservedUtxos, err := k.reserveUtxos(ctx, r.CustodianGroupUid, dataHex, amount,
// 	// 	vsizeLimit)
// 	// if err != nil {
// 	// 	return nil, nil, err
// 	// }

// 	cmdId := NewCommandID(dataHash, r.SourceChain)
// 	payload := &cov.RedeemTokenPayloadWithType{
// 		RedeemTokenPayload: cov.RedeemTokenPayload{
// 			Amount:        req.Amount,
// 			LockingScript: req.LockingScript,
// 			Utxos:         reservedUtxos,
// 			RequestId:     cmdId.Bytes(),
// 		},
// 		PayloadType: encode.ContractCallWithTokenPayloadType_CustodianOnly,
// 	}
// 	redeemTokenParams := &cov.RedeemTokenParams{
// 		DestinationChain:   req.DestChain.String(),
// 		DestinationAddress: req.Address,
// 		Payload:            *payload,
// 		Symbol:             req.Symbol,
// 		Amount:             req.Amount,
// 		CustodianGroupUID:  custodianGrUID,
// 		SessionSequence:    sequence,
// 	}
// 	params, err := redeemTokenParams.AbiPack()
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	return params, &cmdId, err
// }

// func (k Keeper) reserveUtxos(ctx sdk.Context, custodianGroupUID [32]byte, requestID string, amount uint64, vsizeLimit uint64) ([]*cov.UTXO, error) {
// 	// We've already validated redeem session in the outer function
// 	// redeemSession, ok := k.GetRedeemSession(ctx, custodianGroupUID)
// 	// if !ok {
// 	// 	return nil, fmt.Errorf("redeem session not found")
// 	// }
// 	// log.Info().Any("RedeemSession", redeemSession).
// 	// 	Hex("CustodianGroupUid", custodianGroupUID[:]).
// 	// 	Msg("[x/covernant] reserveUtxos found redeem session")
// 	// if redeemSession.IsSwitching {
// 	// 	return nil, fmt.Errorf("redeem phase is switching")
// 	// }
// 	// if redeemSession.CurrentPhase != covExported.Preparing {
// 	// 	return nil, fmt.Errorf("redeem session is not in preparing phase")
// 	// }
// 	group, ok := k.GetCustodianGroup(ctx, custodianGroupUID)
// 	if !ok {
// 		return nil, fmt.Errorf("custodian group not found")
// 	}
// 	utxoSnapshot, ok := k.GetUtxoSnapshot(ctx, custodianGroupUID)
// 	if !ok {
// 		return nil, fmt.Errorf("utxo snapshot not found")
// 	}

// 	if len(utxoSnapshot.Utxos) == 0 {
// 		return nil, fmt.Errorf("no utxos found")
// 	}

// 	// Find optimal UTXO combination using knapsack algorithm
// 	reserveUtxos, err := utxoSnapshot.ReserveUtxos(requestID, amount, uint64(group.Quorum), vsizeLimit)
// 	if err != nil {
// 		k.Logger(ctx).Error("failed to reserve utxos", "error", err)
// 		return nil, err
// 	}

// 	// Update the redeem session in storage
// 	k.setUtxoSnapshot(ctx, utxoSnapshot)

// 	return reserveUtxos, nil
// }
