package models

type CustodianGroup struct {
	UID           []byte   `json:"uid" bson:"uid"`
	Name          string   `json:"name" bson:"name"`
	BitcoinPubkey []byte   `json:"bitcoin_pubkey" bson:"bitcoin_pubkey"`
	Quorum        uint32   `json:"quorum" bson:"quorum"`
	Custodians    [][]byte `json:"custodians" bson:"custodians"`
}
