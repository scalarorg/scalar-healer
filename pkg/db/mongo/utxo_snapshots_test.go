package mongo_test

import (
	"context"
	"testing"

	"github.com/scalarorg/scalar-healer/pkg/db/models"
)

func TestSaveUtxoSnapshot(t *testing.T) {
	utxos := make([]*models.UTXO, 0)
	for i := 0; i < 10; i++ {
		utxos = append(utxos, &models.UTXO{
			TxID:         []byte("txid"),
			Vout:         uint32(i),
			ScriptPubkey: []byte("12345678"),
			AmountInSats: 123456,
			Reservations: []*models.Reservation{},
		})
	}

	utxoSnapshot := &models.UTXOSnapshot{
		BlockHeight:       10000,
		CustodianGroupUID: []byte("custodian_group_uid"),
		UTXOs:             utxos,
	}

	err := repo.SaveUtxoSnapshot(context.Background(), utxoSnapshot)
	if err != nil {
		t.Fatal(err)
	}
}
