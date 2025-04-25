package db

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/scalarorg/data-models/chains"
	"github.com/scalarorg/data-models/scalarnet"
	"github.com/scalarorg/scalar-healer/pkg/db/models"
)

type DbAdapter interface {
	GetGatewayAddress(chainId uint64) *common.Address
	GetAllCustodianGroups() ([]string, error)
	GetChainName(chainType string, chainId uint64) (string, error)
	CheckTokenExists(symbol string) bool
	GetTokenSymbolByAddress(chainType string, chainId uint64, tokenAddress string) (string, error)
	GetLastEventCheckPoint(chainName, eventName string, fromBlock uint64) (*scalarnet.EventCheckPoint, error)
	UpdateLastEventCheckPoint(lastCheckPoint *scalarnet.EventCheckPoint) error
	GetRedeemSession(chainId string, groupUid string) *models.RedeemSession
	SaveSwitchPhaseValue(event *chains.SwitchedPhase) error
	SaveTokenSent(tokenSent *chains.TokenSent, eventCheckPoint *scalarnet.EventCheckPoint) error
	SaveTokenSents(tokenSents []*chains.TokenSent) error
	SaveRedeemTxs(redeemTxs []*chains.RedeemTx) error
	FindPendingBtcTokenSent(chainId string, expectedConfirmBlock int) ([]*chains.TokenSent, error)
	FindPendingRedeemsTransaction(chainId string, expectedConfirmBlock int) ([]*chains.RedeemTx, error)
	UpdateRedeemExecutedCommands(chainId string, txHashes []string) error
	SaveContractCallWithToken(contractCallWithToken *chains.ContractCallWithToken, eventCheckPoint *scalarnet.EventCheckPoint) error
	UpdateEvmCommandExecuted(cmdExecuted *chains.CommandExecuted) error
	GetRedeemNonce(address common.Address) uint64
	SaveRedeemRequest(chainId uint64, address common.Address, signature []byte, amount *big.Int, symbol string, nonce uint64) error
	Close()
}
