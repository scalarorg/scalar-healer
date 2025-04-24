package db

import (
	"github.com/scalarorg/data-models/chains"
	"github.com/scalarorg/data-models/scalarnet"
)

type DbAdapter interface {
	GetLastEventCheckPoint(chainName, eventName string, fromBlock uint64) (*scalarnet.EventCheckPoint, error)
	GetRedeemSession(chainId string, groupUid string) RedeemSession
	SaveSwitchPhaseValue(event *chains.SwitchedPhase) error
	SaveTokenSent(tokenSent *chains.TokenSent, eventCheckPoint *scalarnet.EventCheckPoint) error
	SaveContractCallWithToken(contractCallWithToken *chains.ContractCallWithToken, eventCheckPoint *scalarnet.EventCheckPoint) error
}

func NewDbAdapter(connectionString *string) (DbAdapter, error) {
	return nil, nil
}
