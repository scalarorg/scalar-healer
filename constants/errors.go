package constants

import "fmt"

var (
	ErrInternal         = fmt.Errorf("internal error")
	ErrTokenNotExists   = fmt.Errorf("token not exists")
	ErrInvalidAmount    = fmt.Errorf("invalid amount")
	ErrInvalidSignature = fmt.Errorf("invalid signature")
	ErrNotFoundGateway  = fmt.Errorf("not found gateway address for chain")
	ErrInvalidNonce     = fmt.Errorf("invalid nonce")
)
