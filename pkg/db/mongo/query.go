package mongo

import (
	"context"
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
	return nil, nil
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
	return nil, nil
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

func (m *MongoRepository) GetTokenSymbolByAddress(ctx context.Context, chainType string, chainId uint64, tokenAddress string) (string, error) {
	return "", nil
}

func (m *MongoRepository) CheckTokenExists(ctx context.Context, symbol string) bool {
	result := m.Tokens.FindOne(ctx, map[string]interface{}{
		"symbol": symbol,
	}, options.FindOne().SetProjection(bson.M{
		"symbol": 1,
	}))
	return result.Err() == nil
}

func (m *MongoRepository) GetGatewayAddress(ctx context.Context, chainId uint64) *common.Address {
	var data struct {
		Address []byte `bson:"address"`
	}

	opts := options.FindOne().SetProjection(bson.M{
		"address": 1,
	})
	m.GatewayAddresses.FindOne(ctx, map[string]interface{}{
		"chain_id": chainId,
	}, opts).Decode(&data)
	add := common.BytesToAddress(data.Address)
	return &add
}
