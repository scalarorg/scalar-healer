package postgres

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/scalarorg/scalar-healer/pkg/db"
)

func (m *PostgresRepository) SaveTransferRequest(ctx context.Context, chainId uint64, address common.Address, signature []byte, amount *big.Int, destChain string, destAddress *common.Address, symbol string, nonce uint64) error {
	// currentTime := time.Now().Unix()

	// req := TransferRequest{
	// 	Address:            address.Bytes(),
	// 	Nonce:              nonce,
	// 	Signature:          signature,
	// 	ChainID:            chainId,
	// 	DestinationChain:   destChain,
	// 	DestinationAddress: destAddress.Bytes(),
	// 	Symbol:             symbol,
	// 	Amount:             amount.String(),
	// 	CreatedAt:          currentTime,
	// 	UpdatedAt:          currentTime,
	// }

	// m.TransferRequests.InsertOne(ctx, req)
	return nil
}

func (m *PostgresRepository) ListTransferRequests(ctx context.Context, address common.Address, page, size int32) ([]db.TransferRequest, error) {

	// filter := bson.D{{Key: "address", Value: address.Bytes()}}

	// var req []*TransferRequest
	// opts := options.Find().SetSkip(int64(page * size)).SetLimit(int64(size))
	// cursor, err := m.TransferRequests.Find(ctx, filter, opts)
	// if err != nil {
	// 	return nil, err
	// }
	// defer cursor.Close(ctx)
	// if err = cursor.All(ctx, &req); err != nil {
	// 	return nil, err
	// }

	// return req, nil
	return nil, nil
}
