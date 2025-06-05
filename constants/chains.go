package constants

type Threshold struct {
	Numerator   int64 `json:"numerator"`
	Denominator int64 `json:"denominator"`
}

type Params struct {
	Chain              string            `json:"chain,omitempty"`
	ConfirmationHeight uint64            `json:"confirmation_height,omitempty"`
	VotingThreshold    Threshold         `json:"voting_threshold"`
	Metadata           map[string]string `json:"metadata,omitempty"`
	// Redeem limit
	RedeemSessionAmountLimit uint64 `json:"redeem_session_amount_limit,omitempty"`
	RedeemTxsVsizeLimit      uint64 `json:"redeem_txs_vsize_limit,omitempty"`
}

var EVM_CHAIN_PARAMS = Params{
	Chain:                    "evm",
	ConfirmationHeight:       2,
	VotingThreshold:          Threshold{Numerator: 51, Denominator: 100},
	RedeemSessionAmountLimit: 100_000_000,
	RedeemTxsVsizeLimit:      200_000,
}
