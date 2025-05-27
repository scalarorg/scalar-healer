package chains

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

type ChainName string

func (c ChainName) String() string {
	return string(c)
}

func (m ChainName) GetChainID() (*big.Int, error) {
	parts := strings.Split(m.String(), "|")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid chain name: %s", m.String())
	}
	chainID, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid chain id: %s", parts[1])
	}
	id := big.Int{}
	id.SetUint64(chainID)
	return &id, nil
}
