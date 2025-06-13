package tofnd

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"time"

	ec "github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/scalarorg/scalar-core/x/tss/tofnd"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	Conn    *grpc.ClientConn
	Service tofnd.MultisigClient
	PartyID string
	KeyID   string
	Weight  int
}

func NewClient(cfg *ClientConfig, timeout time.Duration) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, cfg.Address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	return &Client{
		Conn:    conn,
		Service: tofnd.NewMultisigClient(conn),
		PartyID: cfg.PartyID,
		KeyID:   cfg.KeyUID,
		Weight:  cfg.Weight,
	}, nil
}

type SigningResponse struct {
	Sig     *Signature       `json:"signature"`
	Pubkey  *ecdsa.PublicKey `json:"pubkey"`
	PartyID string           `json:"party_id"`
}

func (c *Client) Sign(ctx context.Context, hashRaw []byte) (*SigningResponse, error) {
	if len(hashRaw) != common.HashLength {
		return nil, fmt.Errorf("hash to sign must be 32 bytes")
	}

	hash := common.BytesToHash(hashRaw)

	req := &tofnd.SignRequest{
		KeyUid:    c.KeyID,
		MsgToSign: hashRaw,
		PartyUid:  c.PartyID,
	}

	log.Info().Interface("req", req).Msg("Signing message")

	resp, err := c.Service.Sign(ctx, req)
	if err != nil {
		return nil, err
	}

	ecdsaSig, err := ec.ParseDERSignature(resp.GetSignature())
	if err != nil {
		return nil, err
	}

	sig, pk, err := ToSignature(*ecdsaSig, hash)
	if err != nil {
		return nil, err
	}

	return &SigningResponse{
		Sig:     sig,
		Pubkey:  pk,
		PartyID: c.PartyID,
	}, nil
}
