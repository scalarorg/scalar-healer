package bridge_test

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/scalarorg/scalar-healer/pkg/db/sqlc"
	testutils "github.com/scalarorg/scalar-healer/pkg/test_utils"
	"github.com/scalarorg/scalar-healer/pkg/utils"
	"github.com/zeebo/assert"
)

func TestListBridge(t *testing.T) {
	tests := []struct {
		name          string
		address       string
		page          int
		size          int
		setup         func(t *testing.T)
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name:    "valid request",
			address: "0x24a1dB57Fa3ecAFcbaD91d6Ef068439acEeAe090",
			page:    0,
			size:    10,
			setup: func(t *testing.T) {
				chain := "evm|1"
				address := common.HexToAddress("0x24a1dB57Fa3ecAFcbaD91d6Ef068439acEeAe090")
				signature, _ := hex.DecodeString("8e3bad2520fc46b7f78653a92745812c046df00dee0b29e0a01d670f6de9351a2e6bdd1bd471e95e0a94fd6b4262d173eb50fcc7e6fb3ea3b27823c2d893476b00")
				nonce := uint64(0)
				txHash := common.HexToAddress("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")

				err := dbAdapter.SaveBridgeRequest(context.Background(), chain, address, signature, txHash.Bytes(), nonce)
				assert.NoError(t, err)
			},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recoder.Code)
				data, err := io.ReadAll(recoder.Body)
				assert.NoError(t, err)
				var list []*sqlc.BridgeRequest
				json.Unmarshal(data, &list)
				assert.Equal(t, 1, len(list))
				address := common.HexToAddress("0x24a1dB57Fa3ecAFcbaD91d6Ef068439acEeAe090")
				signature, _ := hex.DecodeString("8e3bad2520fc46b7f78653a92745812c046df00dee0b29e0a01d670f6de9351a2e6bdd1bd471e95e0a94fd6b4262d173eb50fcc7e6fb3ea3b27823c2d893476b00")
				txHash := common.HexToAddress("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
				nonce := sqlc.ConvertUint64ToNumeric(0) // First request should have nonce 0
				assert.Equal(t, address, common.BytesToAddress(list[0].Address))
				assert.Equal(t, signature, list[0].Signature)
				assert.Equal(t, nonce, list[0].Nonce)
				assert.Equal(t, txHash.Bytes(), list[0].TxHash)
			},
		},
		{
			name:    "invalid address",
			address: "invalid_address",
			page:    0,
			size:    10,
			setup:   func(t *testing.T) {},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recoder.Code)
			},
		},
		{
			name:    "negative page",
			address: "0x24a1dB57Fa3ecAFcbaD91d6Ef068439acEeAe090",
			page:    -1,
			size:    10,
			setup:   func(t *testing.T) {},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recoder.Code)
			},
		},
		{
			name:    "negative size",
			address: "0x24a1dB57Fa3ecAFcbaD91d6Ef068439acEeAe090",
			page:    0,
			size:    -1,
			setup:   func(t *testing.T) {},
			checkResponse: func(recoder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, recoder.Code)
			},
		},
	}

	for _, tc := range tests {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			tc.setup(t)
			req, rec := testutils.Request(&testutils.RequestOption{
				Method: http.MethodGet,
				URL:    "/api/bridge/" + tc.address,
				QueryParams: map[string]string{
					"page": utils.IntToString(tc.page),
					"size": utils.IntToString(tc.size),
				},
			})
			testServer.Raw.ServeHTTP(rec, req)
			tc.checkResponse(rec)
			cleanup()
		})
	}
}
