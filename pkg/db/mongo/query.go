package mongo

import (
	"context"
	"encoding/hex"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/scalarorg/data-models/chains"
	"github.com/scalarorg/data-models/scalarnet"
	"github.com/scalarorg/scalar-healer/pkg/db/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (m *MongoRepository) SaveProtocols(ctx context.Context, protocols []models.Protocol) error {
	protocolDocs := make([]interface{}, len(protocols))
	protocolNames := bson.A{}
	for i, protocol := range protocols {
		protocolDocs[i] = protocol
		protocolNames = append(protocolNames, protocol.Name)
	}
	filter := bson.M{
		"name": bson.M{
			"$in": protocolNames,
		},
	}
	collection := m.DB.Collection(COLLECTION_PROTOCOLS)
	collection.DeleteMany(ctx, filter)
	_, err := collection.InsertMany(ctx, protocolDocs)
	return err
}
func (m *MongoRepository) SaveTokenInfos(ctx context.Context, tokens []models.Token) error {
	tokenDocs := make([]interface{}, len(tokens))
	tokenSymbols := bson.A{}
	for i, token := range tokens {
		tokenDocs[i] = token
		tokenSymbols = append(tokenSymbols, token.Symbol)
	}
	_, err := m.Tokens.DeleteMany(ctx, bson.M{
		"symbol": bson.M{
			"$in": tokenSymbols,
		},
	})
	_, err = m.Tokens.InsertMany(ctx, tokenDocs)
	return err
}
func (m *MongoRepository) GetLastEventCheckPoint(ctx context.Context,
	chainName, eventName string, fromBlock uint64) (*scalarnet.EventCheckPoint, error) {
	return nil, nil
}

func (m *MongoRepository) GetRedeemSession(ctx context.Context, chainId string, groupUid string) *models.RedeemSession {
	return nil
}

func (m *MongoRepository) SaveSwitchPhaseValue(ctx context.Context, event *chains.SwitchedPhase) error {
	return nil
}

func (m *MongoRepository) SaveTokenSent(ctx context.Context, tokenSent *chains.TokenSent, eventCheckPoint *scalarnet.EventCheckPoint) error {
	return nil
}

func (m *MongoRepository) SaveContractCallWithToken(ctx context.Context, contractCallWithToken *chains.ContractCallWithToken, eventCheckPoint *scalarnet.EventCheckPoint) error {
	return nil
}

func (m *MongoRepository) GetRedeemNonce(ctx context.Context, address common.Address) uint64 {
	log.Info().Msgf("GetRedeemNonce: %s", address.Hex())
	return 0
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

func (m *MongoRepository) FindPendingBtcTokenSent(ctx context.Context, chainId string, expectedConfirmBlock int) ([]*chains.TokenSent, error) {
	return nil, nil
}
func (m *MongoRepository) FindPendingRedeemsTransaction(ctx context.Context, chainId string, expectedConfirmBlock int) ([]*chains.RedeemTx, error) {
	return nil, nil
}
func (m *MongoRepository) UpdateRedeemExecutedCommands(ctx context.Context, chainId string, txHashes []string) error {
	return nil
}

func (m *MongoRepository) GetAllCustodianGroups(ctx context.Context) ([]string, error) {
	protocolCollection := m.DB.Collection(COLLECTION_PROTOCOLS)
	cursor, err := protocolCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var result []string
	for cursor.Next(ctx) {
		var d struct {
			CustodianGroupUid [32]byte `bson:"custodian_group_uid"`
		}
		err := cursor.Decode(&d)
		if err != nil {
			return nil, err
		}
		result = append(result, hex.EncodeToString(d.CustodianGroupUid[:]))
	}
	return result, nil
}

func (m *MongoRepository) GetChainName(ctx context.Context, chainType string, chainId uint64) (string, error) {
	return "", nil
}

func (m *MongoRepository) UpdateLastEventCheckPoint(ctx context.Context, lastCheckPoint *scalarnet.EventCheckPoint) error {
	return nil
}

func (m *MongoRepository) SaveTokenSents(ctx context.Context, tokenSents []*chains.TokenSent) error {
	return nil
}
func (m *MongoRepository) SaveRedeemTxs(ctx context.Context, redeemTxs []*chains.RedeemTx) error {
	return nil
}

func (m *MongoRepository) UpdateEvmCommandExecuted(ctx context.Context, cmdExecuted *chains.CommandExecuted) error {
	return nil
}

func (m *MongoRepository) GetTokenSymbolByAddress(ctx context.Context, chainId uint64, tokenAddress common.Address) (string, error) {

	filter := bson.M{
		"chain_id": chainId,
		"address":  tokenAddress.Bytes(),
	}
	var data struct {
		Symbol string `bson:"symbol"`
	}
	opts := options.FindOne().SetProjection(bson.M{
		"symbol": 1,
	})
	err := m.Tokens.FindOne(ctx, filter, opts).Decode(&data)
	if err != nil {
		return "", err
	}
	return data.Symbol, nil
}

func (m *MongoRepository) GetTokenAddressBySymbol(ctx context.Context, chainId uint64, tokenSymbol string) (*common.Address, error) {
	filter := bson.M{
		"chain_id": chainId,
		"symbol":   tokenSymbol,
	}
	var data struct {
		Address []byte `bson:"address"`
	}
	opts := options.FindOne().SetProjection(bson.M{
		"address": 1,
	})
	err := m.Tokens.FindOne(ctx, filter, opts).Decode(&data)
	if err != nil {
		return nil, err
	}
	add := common.BytesToAddress(data.Address)
	return &add, nil
}
func (m *MongoRepository) CheckTokenExists(ctx context.Context, symbol string) bool {
	result := m.Tokens.FindOne(ctx, map[string]interface{}{
		"symbol": symbol,
	}, options.FindOne().SetProjection(bson.M{
		"symbol": 1,
	}))
	return result.Err() == nil
}

func (m *MongoRepository) GetGatewayAddress(ctx context.Context, chainId uint64) (*common.Address, error) {
	var data struct {
		Address []byte `bson:"address"`
	}

	opts := options.FindOne().SetProjection(bson.M{
		"address": 1,
	})
	err := m.GatewayAddresses.FindOne(ctx, map[string]interface{}{
		"chain_id": chainId,
	}, opts).Decode(&data)
	if err != nil {
		return nil, err
	}

	add := common.BytesToAddress(data.Address)
	return &add, nil
}
