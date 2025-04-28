package db

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/scalarorg/data-models/chains"
	"github.com/scalarorg/data-models/scalarnet"
	"github.com/scalarorg/scalar-healer/pkg/db/models"
)

type DbAdapter interface {
	SaveProtocols(ctx context.Context, protocols []models.Protocol) error
	SaveTokenInfos(ctx context.Context, tokens []models.Token) error
	GetTokenAddressBySymbol(ctx context.Context, chainId uint64, symbol string) (*common.Address, error)
	GetTokenSymbolByAddress(ctx context.Context, chainId uint64, tokenAddress common.Address) (string, error)
	GetGatewayAddress(ctx context.Context, chainId uint64) (*common.Address, error)
	GetAllCustodianGroups(ctx context.Context) ([]string, error)
	GetChainName(ctx context.Context, chainType string, chainId uint64) (string, error)
	CheckTokenExists(ctx context.Context, symbol string) bool
	GetLastEventCheckPoint(ctx context.Context, chainName, eventName string, fromBlock uint64) (*scalarnet.EventCheckPoint, error)
	UpdateLastEventCheckPoint(ctx context.Context, lastCheckPoint *scalarnet.EventCheckPoint) error
	GetRedeemSession(ctx context.Context, chainId string, groupUid string) *models.RedeemSession
	SaveSwitchPhaseValue(ctx context.Context, event *chains.SwitchedPhase) error
	SaveTokenSent(ctx context.Context, tokenSent *chains.TokenSent, eventCheckPoint *scalarnet.EventCheckPoint) error
	SaveTokenSents(ctx context.Context, tokenSents []*chains.TokenSent) error
	SaveRedeemTxs(ctx context.Context, redeemTxs []*chains.RedeemTx) error
	FindPendingBtcTokenSent(ctx context.Context, chainId string, expectedConfirmBlock int) ([]*chains.TokenSent, error)
	FindPendingRedeemsTransaction(ctx context.Context, chainId string, expectedConfirmBlock int) ([]*chains.RedeemTx, error)
	UpdateRedeemExecutedCommands(ctx context.Context, chainId string, txHashes []string) error
	SaveContractCallWithToken(ctx context.Context, contractCallWithToken *chains.ContractCallWithToken, eventCheckPoint *scalarnet.EventCheckPoint) error
	UpdateEvmCommandExecuted(ctx context.Context, cmdExecuted *chains.CommandExecuted) error
	GetRedeemNonce(ctx context.Context, address common.Address) uint64
	SaveRedeemRequest(ctx context.Context, chainId uint64, address common.Address, signature []byte, amount *big.Int, symbol string, nonce uint64) error
	ListRedeemRequests(ctx context.Context, address common.Address, page, size int) ([]*models.RedeemRequest, error)
	Close()
}
