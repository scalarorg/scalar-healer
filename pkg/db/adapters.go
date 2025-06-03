package db

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/scalarorg/data-models/chains"
	"github.com/scalarorg/data-models/scalarnet"
	"github.com/scalarorg/scalar-healer/pkg/db/sqlc"
)

type HealderAdapter interface {

	// protocols
	GetProtocol(ctx context.Context, asset string) (*sqlc.ProtocolWithTokenDetails, error)
	SaveProtocols(ctx context.Context, protocols []sqlc.Protocol) error
	GetProtocols(ctx context.Context) ([]sqlc.ProtocolWithTokenDetails, error)

	// tokens
	SaveTokens(ctx context.Context, tokens []sqlc.Token) error
	GetTokenAddressBySymbol(ctx context.Context, chain string, symbol string) (*common.Address, error)
	GetTokenSymbolByAddress(ctx context.Context, chain string, tokenAddress *common.Address) (string, error)

	// gateways
	GetGatewayAddress(ctx context.Context, chain string) (*common.Address, error)
	CreateGatewayAddresses(ctx context.Context, addresses [][]byte, chains []string) error

	// chain
	GetChainName(ctx context.Context, chainType string, chainId uint64) (string, error)

	// custodian groups
	SaveCustodianGroups(ctx context.Context, custodianGroups []sqlc.CustodianGroup) error
	GetCustodianGroup(ctx context.Context, uid []byte) (sqlc.CustodianGroup, error)
	GetAllCustodianGroups(ctx context.Context) ([]sqlc.CustodianGroup, error)

	// utxo snapshots
	SaveUtxoSnapshot(ctx context.Context, utxoSnapshot []sqlc.Utxo) error

	// checkpoints
	GetLastEventCheckPoint(ctx context.Context, chainName, eventName string, fromBlock uint64) (*scalarnet.EventCheckPoint, error)
	UpdateLastEventCheckPoint(ctx context.Context, lastCheckPoint *scalarnet.EventCheckPoint) error

	// token-sent
	SaveTokenSent(ctx context.Context, tokenSent *chains.TokenSent, eventCheckPoint *scalarnet.EventCheckPoint) error
	SaveTokenSents(ctx context.Context, tokenSents []chains.TokenSent) error
	FindPendingBtcTokenSent(ctx context.Context, chainId string, expectedConfirmBlock int32) ([]chains.TokenSent, error)

	// bridge
	SaveBridgeRequest(ctx context.Context, chain string, address common.Address, signature []byte, txHash []byte, nonce uint64) error
	ListBridgeRequests(ctx context.Context, address common.Address, page, size int32) ([]sqlc.BridgeRequest, error)

	// transfer
	SaveTransferRequest(ctx context.Context, chain string, address common.Address, signature []byte, amount *big.Int, destChain string, destAddress *common.Address, symbol string, nonce uint64) error
	ListTransferRequests(ctx context.Context, address common.Address, page, size int32) ([]sqlc.TransferRequest, error)

	// redeem
	SaveRedeemTxs(ctx context.Context, redeemTxs []chains.RedeemTx) error
	FindPendingRedeemsTransaction(ctx context.Context, chainId string, expectedConfirmBlock int32) ([]chains.RedeemTx, error)
	UpdateRedeemExecutedCommands(ctx context.Context, chainId string, txHashes []string) error

	SaveRedeemRequest(ctx context.Context, sourceChain, destChain string, address common.Address, amount *big.Int, symbol string, lockingScript []byte) error
	ListRedeemRequests(ctx context.Context, address common.Address, page, size int32) ([]sqlc.RedeemRequest, int64, error)

	// contract calls
	SaveContractCallWithToken(ctx context.Context, contractCallWithToken *chains.ContractCallWithToken, eventCheckPoint *scalarnet.EventCheckPoint) error

	// command executed
	// UpdateEvmCommandExecuted(ctx context.Context, cmdExecuted *chains.CommandExecuted) error

	// accounts
	GetNonce(ctx context.Context, address common.Address) uint64

	// redeem sessions
	SaveRedeemSessionAndChainRedeemSessionsTx(ctx context.Context, chainRedeemSessions []sqlc.ChainRedeemSession) (outdatedSessionsByGroup map[string][]sqlc.ChainRedeemSessionUpdate, err error)
	GetRedeemSession(ctx context.Context, uid []byte) (*sqlc.RedeemSession, error)
	GetCompletedRedeemSessions(ctx context.Context) ([]sqlc.RedeemSession, error)

	// chain redeem sessions
	GetChainRedeemSession(ctx context.Context, grUID []byte, chain string) (*sqlc.ChainRedeemSession, error)

	// commands
	SaveCommands(ctx context.Context, commands []*sqlc.Command) error
	SaveCommandBatches(ctx context.Context, commandBatches []*sqlc.CommandBatch) error
	SaveCommandsAndBatchCommandsTx(ctx context.Context, commands []sqlc.Command) error

	GetCommandBatchByID(ctx context.Context, commandBatchID []byte) (sqlc.CommandBatch, error)
	GetCommandBatches(ctx context.Context) ([]sqlc.CommandBatch, error)

	Close()
}

type IndexerAdapter interface {
	GetNumberOfLatestSwitchedPhaseEvents(ctx context.Context, numberOfEvents int, chain string, grUID string) ([]chains.SwitchedPhase, error)
	GetBatchNumberOfLatestSwitchedPhaseEvents(ctx context.Context, numberOfEvents int, chain string, grUID []string) (map[string][]chains.SwitchedPhase, error)
	GetBatchLastestSwitchedPhaseEvents(ctx context.Context, chain string, grUID []string) (map[string]chains.SwitchedPhase, error)
}

type CombinedAdapter interface {
	HealderAdapter
	IndexerAdapter
	Close()
}
