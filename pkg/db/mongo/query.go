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

func (m *MongoRepository) SaveRedeemRequest(address common.Address, signature []byte, amount *big.Int, symbol string, nonce uint64) error {
	return nil
}
