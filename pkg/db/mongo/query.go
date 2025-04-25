package mongo

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/scalarorg/data-models/chains"
	"github.com/scalarorg/data-models/scalarnet"
	"github.com/scalarorg/scalar-healer/pkg/db/models"
)

func (m *MongoRepository) GetLastEventCheckPoint(chainName, eventName string, fromBlock uint64) (*scalarnet.EventCheckPoint, error) {
	return nil, nil
}

func (m *MongoRepository) GetRedeemSession(chainId string, groupUid string) *models.RedeemSession {
	return nil
}

func (m *MongoRepository) SaveSwitchPhaseValue(event *chains.SwitchedPhase) error {
	return nil
}

func (m *MongoRepository) SaveTokenSent(tokenSent *chains.TokenSent, eventCheckPoint *scalarnet.EventCheckPoint) error {
	return nil
}

func (m *MongoRepository) SaveContractCallWithToken(contractCallWithToken *chains.ContractCallWithToken, eventCheckPoint *scalarnet.EventCheckPoint) error {
	return nil
}

func (m *MongoRepository) GetRedeemNonce(address common.Address) uint64 {
	log.Info().Msgf("GetRedeemNonce: %s", address.Hex())
	return 0
}

func (m *MongoRepository) SaveRedeemRequest(chainId uint64, address common.Address, signature []byte, amount *big.Int, symbol string, nonce uint64) error {
	return nil
}

func (m *MongoRepository) FindPendingBtcTokenSent(chainId string, expectedConfirmBlock int) ([]*chains.TokenSent, error) {
	return nil, nil
}
func (m *MongoRepository) FindPendingRedeemsTransaction(chainId string, expectedConfirmBlock int) ([]*chains.RedeemTx, error) {
	return nil, nil
}
func (m *MongoRepository) UpdateRedeemExecutedCommands(chainId string, txHashes []string) error {
	return nil
}

func (m *MongoRepository) GetAllCustodianGroups() ([]string, error) {
	return nil, nil
}

func (m *MongoRepository) GetChainName(chainType string, chainId uint64) (string, error) {
	return "", nil
}

func (m *MongoRepository) UpdateLastEventCheckPoint(lastCheckPoint *scalarnet.EventCheckPoint) error {
	return nil
}

func (m *MongoRepository) SaveTokenSents(tokenSents []*chains.TokenSent) error {
	return nil
}
func (m *MongoRepository) SaveRedeemTxs(redeemTxs []*chains.RedeemTx) error {
	return nil
}

func (m *MongoRepository) UpdateEvmCommandExecuted(cmdExecuted *chains.CommandExecuted) error {
	return nil
}

func (m *MongoRepository) GetTokenSymbolByAddress(chainType string, chainId uint64, tokenAddress string) (string, error) {
	return "", nil
}

func (m *MongoRepository) CheckTokenExists(symbol string) bool {
	return false
}

func (m *MongoRepository) GetGatewayAddress(chainId uint64) *common.Address {
	return nil
}
