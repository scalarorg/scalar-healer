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

	// protocols
	SaveProtocols(ctx context.Context, protocols []models.Protocol) error

	// tokens
	SaveTokenInfos(ctx context.Context, tokens []models.Token) error
	GetTokenAddressBySymbol(ctx context.Context, chainId uint64, symbol string) (*common.Address, error)
	GetTokenSymbolByAddress(ctx context.Context, chainId uint64, tokenAddress common.Address) (string, error)

	// gateways
	GetGatewayAddress(ctx context.Context, chainId uint64) (*common.Address, error)

	// chain
	GetChainName(ctx context.Context, chainType string, chainId uint64) (string, error)

	// custodian groups
	GetAllCustodianGroups(ctx context.Context) ([]string, error)

	// checkpoints
	GetLastEventCheckPoint(ctx context.Context, chainName, eventName string, fromBlock uint64) (*scalarnet.EventCheckPoint, error)
	UpdateLastEventCheckPoint(ctx context.Context, lastCheckPoint *scalarnet.EventCheckPoint) error

	// session
	GetRedeemSession(ctx context.Context, chainId string, groupUid string) *models.RedeemSession
	SaveSwitchPhaseValue(ctx context.Context, event *chains.SwitchedPhase) error

	// token-sent
	SaveTokenSent(ctx context.Context, tokenSent *chains.TokenSent, eventCheckPoint *scalarnet.EventCheckPoint) error
	SaveTokenSents(ctx context.Context, tokenSents []*chains.TokenSent) error
	FindPendingBtcTokenSent(ctx context.Context, chainId string, expectedConfirmBlock int) ([]*chains.TokenSent, error)
	SaveBridgeRequest(ctx context.Context, chainId uint64, address common.Address, signature []byte, txHash []byte, nonce uint64) error
	ListBridgeRequests(ctx context.Context, address common.Address, page, size int) ([]*models.BridgeRequest, error)
	SaveTransferRequest(ctx context.Context, chainId uint64, address common.Address, signature []byte, amount *big.Int, destChain string, destAddress *common.Address, symbol string, nonce uint64) error
	ListTransferRequests(ctx context.Context, address common.Address, page, size int) ([]*models.TransferRequest, error)

	// redeem
	SaveRedeemTxs(ctx context.Context, redeemTxs []*chains.RedeemTx) error
	FindPendingRedeemsTransaction(ctx context.Context, chainId string, expectedConfirmBlock int) ([]*chains.RedeemTx, error)
	UpdateRedeemExecutedCommands(ctx context.Context, chainId string, txHashes []string) error
	SaveRedeemRequest(ctx context.Context, chainId uint64, address common.Address, signature []byte, amount *big.Int, symbol string, nonce uint64) error
	ListRedeemRequests(ctx context.Context, address common.Address, page, size int) ([]*models.RedeemRequest, error)

	// contract calls
	SaveContractCallWithToken(ctx context.Context, contractCallWithToken *chains.ContractCallWithToken, eventCheckPoint *scalarnet.EventCheckPoint) error

	// command executed
	UpdateEvmCommandExecuted(ctx context.Context, cmdExecuted *chains.CommandExecuted) error

	// accounts
	GetNonce(ctx context.Context, address common.Address) uint64

	Close()
}
