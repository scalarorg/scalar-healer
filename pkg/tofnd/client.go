package tofnd

import (
	"context"
	"time"

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
		KeyID:   cfg.KeyID,
		Weight:  cfg.Weight,
	}, nil
}

type SigningResponse struct {
	Signature []byte `json:"signature"`
	PartyID   string `json:"party_id"`
}

func (r *SigningResponse) GetSignature() []byte {
	return r.Signature
}

func (c *Client) Sign(ctx context.Context, msg []byte) (*SigningResponse, error) {
	// Call the Sign RPC, fill in as needed
	req := &tofnd.SignRequest{
		KeyUid:    c.KeyID,
		MsgToSign: msg,
		PartyUid:  c.PartyID,
	}
	resp, err := c.Service.Sign(ctx, req)
	if err != nil {
		return nil, err
	}

	return &SigningResponse{
		Signature: resp.GetSignature(),
		PartyID:   c.PartyID,
	}, nil
}
