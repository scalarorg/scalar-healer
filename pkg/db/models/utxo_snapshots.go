package models

type Reservation struct {
	Request string `json:"request" bson:"request"`
	Amount  uint64 `json:"amount" bson:"amount"`
}

type UTXO struct {
	TxID         []byte         `json:"txid" bson:"txid"` // in reverse-order
	Vout         uint32         `json:"vout" bson:"vout"`
	ScriptPubkey []byte         `json:"script_pubkey" bson:"script_pubkey"`
	AmountInSats uint64         `json:"amount_in_sats" bson:"amount_in_sats"`
	Reservations []*Reservation `json:"reservations" bson:"reservations"`
}

type UTXOSnapshot struct {
	CustodianGroupUID []byte  `json:"custodian_group_uid" bson:"custodian_group_uid"`
	BlockHeight       uint64  `json:"block_height" bson:"block_height"`
	UTXOs             []*UTXO `json:"utxos" bson:"utxos"`
}
