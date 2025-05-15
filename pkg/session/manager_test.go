package session_test

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/scalarorg/scalar-healer/pkg/session"
	"github.com/stretchr/testify/assert"
)

var (
	secretKey = []byte("test-secret-key")
	expiry    = 24 * time.Hour
)

func TestInit(t *testing.T) {
	err := session.Init(secretKey, expiry)
	assert.NoError(t, err)
	err = session.Init(secretKey, expiry)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already initialized")
}

func TestCreateToken(t *testing.T) {
	err := session.Init(secretKey, expiry)
	if err != nil {
		t.Log("failed to initialize session manager")
	}

	address := common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc454e4438f44e")
	token, err := session.CreateToken(address)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	validatedAddr, err := session.ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, address, *validatedAddr)
}
