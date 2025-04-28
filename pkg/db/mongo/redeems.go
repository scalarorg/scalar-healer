package mongo

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/scalarorg/data-models/chains"
	"github.com/scalarorg/scalar-healer/pkg/db/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (m *MongoRepository) FindPendingRedeemsTransaction(ctx context.Context, chainId string, expectedConfirmBlock int) ([]*chains.RedeemTx, error) {
	return nil, nil
}
func (m *MongoRepository) UpdateRedeemExecutedCommands(ctx context.Context, chainId string, txHashes []string) error {
	return nil
}

func (m *MongoRepository) SaveRedeemTxs(ctx context.Context, redeemTxs []*chains.RedeemTx) error {
	return nil
}

func (m *MongoRepository) SaveRedeemRequest(ctx context.Context, chainId uint64, address common.Address, signature []byte, amount *big.Int, symbol string, nonce uint64) error {
	currentTime := time.Now().Unix()

	redeemRequest := models.RedeemRequest{
		Address:   address.Bytes(),
		Amount:    amount.String(),
		Symbol:    symbol,
		Nonce:     nonce,
		Signature: signature,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		ChainID:   chainId,
	}

	m.RedeemRequests.InsertOne(ctx, redeemRequest)
	return nil
}

func (m *MongoRepository) ListRedeemRequests(ctx context.Context, address common.Address, page, size int) ([]*models.RedeemRequest, error) {

	filter := bson.D{{Key: "address", Value: address.Bytes()}}

	var redeemRequests []*models.RedeemRequest
	opts := options.Find().SetSkip(int64(page * size)).SetLimit(int64(size))
	cursor, err := m.RedeemRequests.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &redeemRequests); err != nil {
		return nil, err
	}

	return redeemRequests, nil
}
