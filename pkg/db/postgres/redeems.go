package postgres

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/scalarorg/data-models/chains"
	"github.com/scalarorg/scalar-healer/pkg/db"
)

func (m *PostgresRepository) FindPendingRedeemsTransaction(ctx context.Context, chainId string, expectedConfirmBlock int32) ([]chains.RedeemTx, error) {
	return nil, nil
}
func (m *PostgresRepository) UpdateRedeemExecutedCommands(ctx context.Context, chainId string, txHashes []string) error {
	return nil
}

func (m *PostgresRepository) SaveRedeemTxs(ctx context.Context, redeemTxs []chains.RedeemTx) error {
	return nil
}

func (m *PostgresRepository) SaveRedeemRequest(ctx context.Context, chainId uint64, address common.Address, signature []byte, amount *big.Int, symbol string, nonce uint64) error {
	// currentTime := time.Now().Unix()

	// redeemRequest := RedeemRequest{
	// 	Address:   address.Bytes(),
	// 	Amount:    amount.String(),
	// 	Symbol:    symbol,
	// 	Nonce:     nonce,
	// 	Signature: signature,
	// 	CreatedAt: currentTime,
	// 	UpdatedAt: currentTime,
	// 	ChainID:   chainId,
	// }

	// m.RedeemRequests.InsertOne(ctx, redeemRequest)
	return nil
}

func (m *PostgresRepository) ListRedeemRequests(ctx context.Context, address common.Address, page, size int32) ([]db.RedeemRequest, error) {

	// filter := bson.D{{Key: "address", Value: address.Bytes()}}

	// var redeemRequests []*RedeemRequest
	// opts := options.Find().SetSkip(int64(page * size)).SetLimit(int64(size))
	// cursor, err := m.RedeemRequests.Find(ctx, filter, opts)
	// if err != nil {
	// 	return nil, err
	// }
	// defer cursor.Close(ctx)
	// if err = cursor.All(ctx, &redeemRequests); err != nil {
	// 	return nil, err
	// }

	// return redeemRequests, nil
	return nil, nil
}
