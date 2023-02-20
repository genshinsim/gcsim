package share

import (
	context "context"
	"fmt"

	"github.com/genshinsim/gcsim/pkg/model"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientCfg struct {
	Addr                string
	DefaultTTLInSeconds uint64
}

type Client struct {
	cfg       ClientCfg
	srvClient ShareStoreClient
}

func NewClient(cfg ClientCfg, cust ...func(*Client) error) (*Client, error) {
	c := &Client{cfg: cfg}

	conn, err := grpc.Dial(cfg.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	c.srvClient = NewShareStoreClient(conn)

	return c, nil
}

func (c *Client) Create(ctx context.Context, data *model.SimulationResult) (string, error) {
	//TODO: check ctx for perma settings
	resp, err := c.srvClient.Create(ctx, &CreateRequest{
		Result: data,
	})
	if err != nil {
		return "", err
	}
	return resp.GetKey(), nil
}

func (c *Client) Update(ctx context.Context, id string, data *model.SimulationResult) error {
	return fmt.Errorf("not implemented")
}

func (c *Client) SetTTL(ctx context.Context, id string) error {
	return fmt.Errorf("not implemented")
}

func (c *Client) Delete(ctx context.Context, id string) error {
	return fmt.Errorf("not implemented")
}

func (c *Client) Random(ctx context.Context) (string, error) {
	return "", fmt.Errorf("not implemented")
}

func (c *Client) Read(ctx context.Context, id string) (*model.SimulationResult, uint64, error) {
	resp, err := c.srvClient.Read(ctx, &ReadRequest{
		Key: id,
	})

	if err != nil {
		return nil, 0, err
	}

	return resp.GetResult(), resp.GetExpiresAt(), nil
}
