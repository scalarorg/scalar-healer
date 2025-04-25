package db

import (
	"github.com/scalarorg/data-models/chains"
	"github.com/scalarorg/data-models/scalarnet"
	"github.com/scalarorg/scalar-healer/pkg/db/models"
)

type DbAdapter interface {
	GetAllCustodianGroups() ([]string, error)
	GetChainName(chainType string, chainId uint64) (string, error)
	GetTokenSymbolByAddress(chainType string, chainId uint64, tokenAddress string) (string, error)
	GetLastEventCheckPoint(chainName, eventName string, fromBlock uint64) (*scalarnet.EventCheckPoint, error)
	UpdateLastEventCheckPoint(lastCheckPoint *scalarnet.EventCheckPoint) error
	GetRedeemSession(chainId string, groupUid string) models.RedeemSession
	SaveSwitchPhaseValue(event *chains.SwitchedPhase) error
	SaveTokenSent(tokenSent *chains.TokenSent, eventCheckPoint *scalarnet.EventCheckPoint) error
	SaveTokenSents(tokenSents []*chains.TokenSent) error
	SaveRedeemTxs(redeemTxs []*chains.RedeemTx) error
	FindPendingBtcTokenSent(chainId string, expectedConfirmBlock int) ([]*chains.TokenSent, error)
	FindPendingRedeemsTransaction(chainId string, expectedConfirmBlock int) ([]*chains.RedeemTx, error)
	UpdateRedeemExecutedCommands(chainId string, txHashes []string) error
	SaveContractCallWithToken(contractCallWithToken *chains.ContractCallWithToken, eventCheckPoint *scalarnet.EventCheckPoint) error
	UpdateEvmCommandExecuted(cmdExecuted *chains.CommandExecuted) error
}

func NewDbAdapter(connectionString *string) (DbAdapter, error) {
	return nil, nil
}
