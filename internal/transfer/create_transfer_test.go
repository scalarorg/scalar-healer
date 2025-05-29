package transfer_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/scalarorg/scalar-healer/constants"
	"github.com/scalarorg/scalar-healer/internal/transfer"
	"github.com/scalarorg/scalar-healer/pkg/crypto/eip712"
	testutils "github.com/scalarorg/scalar-healer/pkg/test_utils"
	"github.com/zeebo/assert"
)

func TestCreateTransferRequest(t *testing.T) {
	tests := []struct {
		name          string
		request       transfer.CreateTransferRequest
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "valid request",
			request: transfer.CreateTransferRequest{
				BaseRequest: eip712.BaseRequest{
					Address:   common.HexToAddress("0x24a1dB57Fa3ecAFcbaD91d6Ef068439acEeAe090"),
					Signature: "c5b79a60eba368238f3628dd747fe091366735c9edf5d52afb1091e2d90ce8647a7eea63b31b324c19f55ac081ac487c2d0ded7bd7b66e26002b56b8e31eb60201",
					Chain:     "evm|1",
					Nonce:     0,
				},
				DestinationChain:   "evm|11155111",
				DestinationAddress: common.MaxAddress.Hex(),
				Symbol:             "ETH",
				Amount:             "123456",
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, rec.Code, http.StatusOK)
			},
		},
		{
			name: "binding error",
			request: transfer.CreateTransferRequest{
				BaseRequest: eip712.BaseRequest{
					Address:   common.HexToAddress("0x24a1dB57Fa3ecAFcbaD91d6Ef068439acEeAe09k"),
					Signature: "9efeb92deab95a1183de55fe6be80afec2ab5f4766515eb706c16a1dafb2b6186da17b036134738509fc23ac4e9882bf66c4bbefb9b41d73ed72568de261036901",
					Chain:     "evm|1",
					Nonce:     0,
				},
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, rec.Code, http.StatusBadRequest)
			},
		},
		{
			name: "invalid chain id",
			request: transfer.CreateTransferRequest{
				BaseRequest: eip712.BaseRequest{
					Address:   common.HexToAddress("0x24a1dB57Fa3ecAFcbaD91d6Ef068439acEeAe090"),
					Signature: "9efeb92deab95a1183de55fe6be80afec2ab5f4766515eb706c16a1dafb2b6186da17b036134738509fc23ac4e9882bf66c4bbefb9b41d73ed72568de261036901",
					Chain:     "evm|1000",
					Nonce:     0,
				},
				DestinationChain:   "evm|11155111",
				DestinationAddress: common.MaxAddress.Hex(),
				Symbol:             "ETH",
				Amount:             "123456",
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, rec.Code, http.StatusBadRequest)
				var response map[string]string
				err := json.NewDecoder(rec.Body).Decode(&response)
				assert.NoError(t, err)
				data := response["message"]
				assert.Equal(t, data, constants.ErrNotFoundGateway.Error())
			},
		},
		{
			name: "invalid signature",
			request: transfer.CreateTransferRequest{
				BaseRequest: eip712.BaseRequest{
					Address:   common.HexToAddress("0x24a1dB57Fa3ecAFcbaD91d6Ef068439acEeAe090"),
					Signature: "7381418bca4505e78251271a98c1e1a44bfe4a1ea4884f80dbf0b1b5a12d542639bd98b71e787f6a19566424f8c90a874c93d4788467d7e9bbfb65ec10a602a900",
					Chain:     "evm|1",
					Nonce:     0,
				},
				DestinationChain:   "evm|11155111",
				DestinationAddress: common.MaxAddress.Hex(),
				Symbol:             "ETH",
				Amount:             "123456",
			},
			checkResponse: func(rec *httptest.ResponseRecorder) {
				assert.Equal(t, rec.Code, http.StatusBadRequest)
				var response map[string]string
				err := json.NewDecoder(rec.Body).Decode(&response)
				assert.NoError(t, err)
				data := response["message"]
				assert.Equal(t, data, constants.ErrInvalidSignature.Error())
			},
		},
	}

	for _, tc := range tests {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			req, rec := testutils.Request(&testutils.RequestOption{
				Method: http.MethodPost,
				URL:    "/api/transfer",
				Body:   tc.request,
			})

			testServer.Raw.ServeHTTP(rec, req)
			tc.checkResponse(rec)
			cleanup()
		})
	}
}
