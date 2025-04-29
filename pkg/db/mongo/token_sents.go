package mongo

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/scalarorg/data-models/chains"
	"github.com/scalarorg/data-models/scalarnet"
	"github.com/scalarorg/scalar-healer/pkg/db/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (m *MongoRepository) SaveTokenSent(ctx context.Context, tokenSent *chains.TokenSent, eventCheckPoint *scalarnet.EventCheckPoint) error {
	return nil
}

func (m *MongoRepository) SaveTokenSents(ctx context.Context, tokenSents []*chains.TokenSent) error {
	return nil
}

func (m *MongoRepository) FindPendingBtcTokenSent(ctx context.Context, chainId string, expectedConfirmBlock int) ([]*chains.TokenSent, error) {
	return nil, nil
}

func (m *MongoRepository) SaveBridgeRequest(ctx context.Context, chainId uint64, address common.Address, signature []byte, txHash []byte, nonce uint64) error {
	currentTime := time.Now().Unix()

	req := models.BridgeRequest{
		Address:   address.Bytes(),
		TxHash:    txHash,
		Nonce:     nonce,
		Signature: signature,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		ChainID:   chainId,
	}

	m.BridgeRequests.InsertOne(ctx, req)
	return nil
}

func (m *MongoRepository) ListBridgeRequests(ctx context.Context, address common.Address, page, size int) ([]*models.BridgeRequest, error) {

	filter := bson.D{{Key: "address", Value: address.Bytes()}}

	var req []*models.BridgeRequest
	opts := options.Find().SetSkip(int64(page * size)).SetLimit(int64(size))
	cursor, err := m.BridgeRequests.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &req); err != nil {
		return nil, err
	}

	return req, nil
}

func (m *MongoRepository) SaveTransferRequest(ctx context.Context, chainId uint64, address common.Address, signature []byte, amount *big.Int, destChain string, destAddress *common.Address, symbol string, nonce uint64) error {
	currentTime := time.Now().Unix()

	req := models.TransferRequest{
		Address:            address.Bytes(),
		Nonce:              nonce,
		Signature:          signature,
		ChainID:            chainId,
		DestinationChain:   destChain,
		DestinationAddress: destAddress.Bytes(),
		Symbol:             symbol,
		Amount:             amount.String(),
		CreatedAt:          currentTime,
		UpdatedAt:          currentTime,
	}

	m.TransferRequests.InsertOne(ctx, req)
	return nil
}

func (m *MongoRepository) ListTransferRequests(ctx context.Context, address common.Address, page, size int) ([]*models.TransferRequest, error) {

	filter := bson.D{{Key: "address", Value: address.Bytes()}}

	var req []*models.TransferRequest
	opts := options.Find().SetSkip(int64(page * size)).SetLimit(int64(size))
	cursor, err := m.TransferRequests.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &req); err != nil {
		return nil, err
	}

	return req, nil
}
