package evm

import (
	"encoding/hex"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/scalarorg/bitcoin-vault/go-utils/encode"
	"github.com/scalarorg/scalar-healer/pkg/db/sqlc"
)

type RedeemTokenPayload struct {
	Amount        uint64
	LockingScript []byte
	Utxos         []sqlc.Utxo
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
		vouts[i] = uint32(utxo.Vout)
		amounts[i] = sqlc.ConvertNumericToUint64(utxo.AmountInSats)
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
	p.Utxos = make([]sqlc.Utxo, len(txIds))
	for i, txId := range txIds {
		hash, err := hex.DecodeString(strings.TrimPrefix(txId, "0x"))
		if err != nil {
			log.Error().Err(err).Msg("txId hash error")
			return err
		}
		p.Utxos[i] = sqlc.Utxo{TxID: hash, Vout: int64(vouts[i]), AmountInSats: sqlc.ConvertUint64ToNumeric(amounts[i])}
	}
	p.RequestId = unpacked[5].([32]byte)
	return nil
}
