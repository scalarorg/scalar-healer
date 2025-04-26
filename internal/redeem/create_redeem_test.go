package redeem_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/scalarorg/scalar-healer/internal/redeem"
	"github.com/scalarorg/scalar-healer/pkg/utils"
	"github.com/zeebo/assert"
)

func TestCreateRedeem(t *testing.T) {
	tests := []struct {
		name           string
		request        redeem.CreateRedeemRequest
		expectedStatus int
		expectedError  string
	}{
		{
			name: "valid request",
			request: redeem.CreateRedeemRequest{
				Address:   "0x24a1dB57Fa3ecAFcbaD91d6Ef068439acEeAe090",
				Signature: "0xf6c5691b0cd1120058f8a4ed75cd67065a8cdcefaa34ff55678ce1fcab07e0c91357e525c94b97e78b558e3cfe44eb66e3de28cc0d65a6c11c910fff0fabad0100",
				ChainID:   1,
				Symbol:    "ETH",
				Amount:    "123456",
				Nonce:     0, // First request should have nonce 0
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid chain ID",
			request: redeem.CreateRedeemRequest{
				Address:   "0x24a1dB57Fa3ecAFcbaD91d6Ef068439acEeAe090",
				Signature: "0xf6c5691b0cd1120058f8a4ed75cd67065a8cdcefaa34ff55678ce1fcab07e0c91357e525c94b97e78b558e3cfe44eb66e3de28cc0d65a6c11c910fff0fabad0100",
				ChainID:   999,
				Symbol:    "ETH",
				Amount:    "1000000000000000000",
				Nonce:     0,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "not found gateway address for chain: 999",
		},
		{
			name: "invalid nonce",
			request: redeem.CreateRedeemRequest{
				Address:   "0x1234567890123456789012345678901234567890",
				Signature: "0xf6c5691b0cd1120058f8a4ed75cd67065a8cdcefaa34ff55678ce1fcab07e0c91357e525c94b97e78b558e3cfe44eb66e3de28cc0d65a6c11c910fff0fabad0100",
				ChainID:   1,
				Symbol:    "ETH",
				Amount:    "1000000000000000000",
				Nonce:     2,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid nonce",
		},
		{
			name: "bind error",
			request: redeem.CreateRedeemRequest{
				Address:   "D91d6Ef068439acEeAe090",
				Signature: "aaaa",
				ChainID:   1,
				Symbol:    "ETH",
				Amount:    "1000000000000000000",
				Nonce:     0,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "",
		},
		{
			name: "invalid token",
			request: redeem.CreateRedeemRequest{
				Address:   "0x1234567890123456789012345678901234567890",
				Signature: "0xf6c5691b0cd112005",
				ChainID:   1,
				Symbol:    "BTC",
				Amount:    "100",
				Nonce:     0,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "token not exists",
		},
		{
			name: "invalid amount",
			request: redeem.CreateRedeemRequest{
				Address:   "0x1234567890123456789012345678901234567890",
				Signature: "0xf6c5691b0cd112005",
				ChainID:   1,
				Symbol:    "ETH",
				Amount:    "12312321aaaa",
				Nonce:     0,
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid amount",
		},
		{
			name: "invalid signature",
			request: redeem.CreateRedeemRequest{
				Address:   "0x24a1dB57Fa3ecAFcbaD91d6Ef068439acEeAe090",
				Signature: "0xf6c5691b0cd1120058f8a4ed75cd67065a8cdcefaa34ff55678ce1fcab07e0c91357e525c94b97e78b558e3cfe44eb66e3de28cc0d65a6c11c910fff0fabad0101",
				ChainID:   1,
				Symbol:    "ETH",
				Amount:    "123456",
				Nonce:     0, // First request should have nonce 0
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "failed to verify signature: invalid signature",
		},
	}

	for _, tc := range tests {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			req, rec := utils.Request(&utils.RequestOption{
				Method: http.MethodPost,
				URL:    "/api/redeem",
				Body:   tc.request,
			})

			testServer.Raw.ServeHTTP(rec, req)

			if tc.expectedError != "" {
				assert.Equal(t, tc.expectedStatus, rec.Code)
				var response map[string]string
				err := json.NewDecoder(rec.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedError, response["message"])
			} else {
				assert.Equal(t, tc.expectedStatus, rec.Code)
			}
		})
	}
}
