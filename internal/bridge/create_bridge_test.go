package bridge_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/scalarorg/scalar-healer/internal/bridge"
	"github.com/scalarorg/scalar-healer/pkg/crypto/eip712"
	testutils "github.com/scalarorg/scalar-healer/pkg/test_utils"
	"github.com/zeebo/assert"
)

func TestCreateBridgeRequest(t *testing.T) {
	tests := []struct {
		name          string
		request       bridge.CreateBridgeRequest
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "valid request",
			request: bridge.CreateBridgeRequest{
				BaseRequest: eip712.BaseRequest{
					Address:   "0x24a1dB57Fa3ecAFcbaD91d6Ef068439acEeAe090",
					Signature: "8e3bad2520fc46b7f78653a92745812c046df00dee0b29e0a01d670f6de9351a2e6bdd1bd471e95e0a94fd6b4262d173eb50fcc7e6fb3ea3b27823c2d893476b00",
					ChainID:   1,
					Nonce:     0,
				},
				TxHash: "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, recorder.Code, http.StatusOK)
			},
		},
		{
			name: "address binding error",
			request: bridge.CreateBridgeRequest{
				BaseRequest: eip712.BaseRequest{
					Address:   "0x24a1dB57Fa3ecAFcbaD91d6Ef068439acEeAe09k",
					Signature: "9efeb92deab95a1183de55fe6be80afec2ab5f4766515eb706c16a1dafb2b6186da17b036134738509fc23ac4e9882bf66c4bbefb9b41d73ed72568de261036901",
					ChainID:   1,
					Nonce:     0,
				},
				TxHash: "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, rec.Code, http.StatusBadRequest)
			},
		},
		{
			name: "invalid chain id",
			request: bridge.CreateBridgeRequest{
				BaseRequest: eip712.BaseRequest{
					Address:   "0x24a1dB57Fa3ecAFcbaD91d6Ef068439acEeAe090",
					Signature: "9efeb92deab95a1183de55fe6be80afec2ab5f4766515eb706c16a1dafb2b6186da17b036134738509fc23ac4e9882bf66c4bbefb9b41d73ed72568de261036901",
					ChainID:   1000,
					Nonce:     0,
				},
				TxHash: "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, rec.Code, http.StatusBadRequest)
				var response map[string]string
				err := json.NewDecoder(rec.Body).Decode(&response)
				assert.NoError(t, err)
				data := response["message"]
				assert.Equal(t, data, "not found gateway address for chain: 1000")
			},
		},
		{
			name: "invalid signature",
			request: bridge.CreateBridgeRequest{
				BaseRequest: eip712.BaseRequest{
					Address:   "0x24a1dB57Fa3ecAFcbaD91d6Ef068439acEeAe090",
					Signature: "7381418bca4505e78251271a98c1e1a44bfe4a1ea4884f80dbf0b1b5a12d542639bd98b71e787f6a19566424f8c90a874c93d4788467d7e9bbfb65ec10a602a900",
					ChainID:   1,
					Nonce:     0,
				},
				TxHash: "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, rec.Code, http.StatusBadRequest)
				var response map[string]string
				err := json.NewDecoder(rec.Body).Decode(&response)
				assert.NoError(t, err)
				data := response["message"]
				assert.Equal(t, data, "invalid signature")
			},
		},
	}

	for _, tc := range tests {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			req, rec := testutils.Request(&testutils.RequestOption{
				Method: http.MethodPost,
				URL:    "/api/bridge",
				Body:   tc.request,
			})

			testServer.Raw.ServeHTTP(rec, req)
			tc.checkResponse(rec)
		})
	}
}
