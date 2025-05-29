package redeem_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/scalarorg/scalar-healer/constants"
	"github.com/scalarorg/scalar-healer/internal/redeem"
	"github.com/scalarorg/scalar-healer/pkg/crypto/eip712"
	testutils "github.com/scalarorg/scalar-healer/pkg/test_utils"
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
				BaseRequest: eip712.BaseRequest{
					Address:   common.HexToAddress("0x24a1dB57Fa3ecAFcbaD91d6Ef068439acEeAe090"),
					Signature: "0xf6c5691b0cd1120058f8a4ed75cd67065a8cdcefaa34ff55678ce1fcab07e0c91357e525c94b97e78b558e3cfe44eb66e3de28cc0d65a6c11c910fff0fabad0100",
					Nonce:     0, // First request should have nonce 0
					Chain:     "evm|1",
				},
				Symbol: "ETH",
				Amount: "123456",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid chain ID",
			request: redeem.CreateRedeemRequest{
				BaseRequest: eip712.BaseRequest{
					Address:   common.HexToAddress("0x24a1dB57Fa3ecAFcbaD91d6Ef068439acEeAe090"),
					Signature: "f6c5691b0cd1120058f8a4ed75cd67065a8cdcefaa34ff55678ce1fcab07e0c91357e525c94b97e78b558e3cfe44eb66e3de28cc0d65a6c11c910fff0fabad0100",
					Nonce:     0,
					Chain:     "evm|10000",
				},
				Symbol: "ETH",
				Amount: "1000000000000000000",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  constants.ErrNotFoundGateway.Error(),
		},
		{
			name: "bind error",
			request: redeem.CreateRedeemRequest{
				BaseRequest: eip712.BaseRequest{
					Address:   common.HexToAddress("D91d6Ef068439acEeAe090"),
					Signature: "aaaa",
					Nonce:     0,
					Chain:     "evm|1",
				},
				Symbol: "ETH",
				Amount: "1000000000000000000",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "",
		},
		{
			name: "invalid token",
			request: redeem.CreateRedeemRequest{
				BaseRequest: eip712.BaseRequest{
					Address:   common.HexToAddress("0x1234567890123456789012345678901234567890"),
					Signature: "0xf6c5691b0cd112005",
					Nonce:     0,
					Chain:     "evm|1",
				},
				Symbol: "BTC",
				Amount: "100",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  constants.ErrTokenNotExists.Error(),
		},
		{
			name: "invalid amount",
			request: redeem.CreateRedeemRequest{
				BaseRequest: eip712.BaseRequest{
					Address:   common.HexToAddress("0x1234567890123456789012345678901234567890"),
					Signature: "0xf6c5691b0cd112005",
					Nonce:     0,
					Chain:     "evm|1",
				},
				Symbol: "ETH",
				Amount: "12312321aaaa",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  constants.ErrInvalidAmount.Error(),
		},
		{
			name: "invalid signature",
			request: redeem.CreateRedeemRequest{
				BaseRequest: eip712.BaseRequest{
					Address:   common.HexToAddress("0x24a1dB57Fa3ecAFcbaD91d6Ef068439acEeAe090"),
					Signature: "0xf6c5691b0cd1120058f8a4ed75cd67065a8cdcefaa34ff55678ce1fcab07e0c91357e525c94b97e78b558e3cfe44eb66e3de28cc0d65a6c11c910fff0fabad0101",
					Nonce:     0, // First request should have nonce 0
					Chain:     "evm|1",
				},
				Symbol: "ETH",
				Amount: "123456",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  constants.ErrInvalidSignature.Error(),
		},
	}

	for _, tc := range tests {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			req, rec := testutils.Request(&testutils.RequestOption{
				Method: http.MethodPost,
				URL:    "/api/redeem",
				Body:   tc.request,
			})

			testServer.Raw.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			if tc.expectedError != "" {
				var response map[string]string
				err := json.NewDecoder(rec.Body).Decode(&response)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedError, response["message"])
			}

			cleanup()
		})
	}
}
