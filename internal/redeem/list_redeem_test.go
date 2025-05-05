package redeem_test

// func TestListRedeem(t *testing.T) {
// 	tests := []struct {
// 		name          string
// 		address       string
// 		page          int
// 		size          int
// 		setup         func(t *testing.T)
// 		checkResponse func(recoder *httptest.ResponseRecorder)
// 	}{
// 		{
// 			name:    "valid request",
// 			address: "0x24a1dB57Fa3ecAFcbaD91d6Ef068439acEeAe090",
// 			page:    0,
// 			size:    10,
// 			setup: func(t *testing.T) {
// 				chainId := uint64(1)
// 				address := common.HexToAddress("0x24a1dB57Fa3ecAFcbaD91d6Ef068439acEeAe090")
// 				signature, _ := hex.DecodeString("f6c5691b0cd1120058f8a4ed75cd67065a8cdcefaa34ff55678ce1fcab07e0c91357e525c94b97e78b558e3cfe44eb66e3de28cc0d65a6c11c910fff0fabad0100")
// 				amount, _ := utils.StringToBigInt("123456")
// 				symbol := "ETH"
// 				nonce := uint64(0) // First request should have nonce 0

// 				err := db.SaveRedeemRequest(context.Background(), chainId, address, signature, amount, symbol, nonce)
// 				assert.NoError(t, err)
// 			},
// 			checkResponse: func(recoder *httptest.ResponseRecorder) {
// 				assert.Equal(t, http.StatusOK, recoder.Code)
// 				data, err := io.ReadAll(recoder.Body)
// 				assert.NoError(t, err)
// 				var list []*models.RedeemRequest
// 				json.Unmarshal(data, &list)
// 				assert.Equal(t, 1, len(list))
// 				address := common.HexToAddress("0x24a1dB57Fa3ecAFcbaD91d6Ef068439acEeAe090")
// 				signature, _ := hex.DecodeString("f6c5691b0cd1120058f8a4ed75cd67065a8cdcefaa34ff55678ce1fcab07e0c91357e525c94b97e78b558e3cfe44eb66e3de28cc0d65a6c11c910fff0fabad0100")
// 				amount, _ := utils.StringToBigInt("123456")
// 				symbol := "ETH"
// 				nonce := uint64(0) // First request should have nonce 0
// 				assert.Equal(t, address, common.BytesToAddress(list[0].Address))
// 				assert.Equal(t, signature, list[0].Signature)
// 				assert.Equal(t, amount.String(), list[0].Amount)
// 				assert.Equal(t, symbol, list[0].Symbol)
// 				assert.Equal(t, nonce, list[0].Nonce)
// 			},
// 		},
// 		{
// 			name:    "invalid address",
// 			address: "invalid_address",
// 			page:    0,
// 			size:    10,
// 			setup:   func(t *testing.T) {},
// 			checkResponse: func(recoder *httptest.ResponseRecorder) {
// 				assert.Equal(t, http.StatusBadRequest, recoder.Code)
// 			},
// 		},
// 		{
// 			name:    "negative page",
// 			address: "0x24a1dB57Fa3ecAFcbaD91d6Ef068439acEeAe090",
// 			page:    -1,
// 			size:    10,
// 			setup:   func(t *testing.T) {},
// 			checkResponse: func(recoder *httptest.ResponseRecorder) {
// 				assert.Equal(t, http.StatusBadRequest, recoder.Code)
// 			},
// 		},
// 		{
// 			name:    "negative size",
// 			address: "0x24a1dB57Fa3ecAFcbaD91d6Ef068439acEeAe090",
// 			page:    0,
// 			size:    -1,
// 			setup:   func(t *testing.T) {},
// 			checkResponse: func(recoder *httptest.ResponseRecorder) {
// 				assert.Equal(t, http.StatusBadRequest, recoder.Code)
// 			},
// 		},
// 	}

// 	for _, tc := range tests {
// 		tc := tc // capture range variable
// 		t.Run(tc.name, func(t *testing.T) {
// 			tc.setup(t)
// 			req, rec := utils.Request(&utils.RequestOption{
// 				Method: http.MethodGet,
// 				URL:    "/api/redeem/" + tc.address,
// 				QueryParams: map[string]string{
// 					"page": utils.IntToString(tc.page),
// 					"size": utils.IntToString(tc.size),
// 				},
// 			})
// 			testServer.Raw.ServeHTTP(rec, req)
// 			tc.checkResponse(rec)
// 		})
// 	}
// }
