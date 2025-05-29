package chains

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

const SEPARATOR = "|"

type ChainName string

func (c ChainName) String() string {
	return string(c)
}

func (m ChainName) GetParts() (*string, *string, error) {
	parts := strings.Split(m.String(), SEPARATOR)
	if len(parts) != 2 {
		return nil, nil, fmt.Errorf("invalid chain name: %s", m.String())
	}
	return &parts[0], &parts[1], nil
}

func (m ChainName) GetChainID() (*big.Int, error) {
	_, chainId, err := m.GetParts()
	if err != nil {
		return nil, err
	}
	chainID, err := strconv.ParseUint(*chainId, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid chain id: %s", *chainId)
	}
	id := big.Int{}
	id.SetUint64(chainID)
	return &id, nil
}

func (m ChainName) IsEvmChain() bool {
	chainFamily, _, err := m.GetParts()
	if err != nil {
		return false
	}
	return *chainFamily == "evm"
}
