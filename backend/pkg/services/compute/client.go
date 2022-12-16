package compute

import (
	"context"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientCfg struct {
	Addr   string
	APIKey string
}

type Client struct {
	cfg       ClientCfg
	srvClient ComputeClient
}

func NewClient(cfg ClientCfg, cust ...func(*Client) error) (*Client, error) {
	c := &Client{cfg: cfg}

	conn, err := grpc.Dial(cfg.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	c.srvClient = NewComputeClient(conn)

	return c, nil
}

func (c *Client) Run(key, cfg string, ctx context.Context) error {
	req := &RunRequest{
		Key:    key,
		Config: cfg,
		ApiKey: c.cfg.APIKey,
	}
	_, err := c.srvClient.Run(ctx, req)

	//TODO: read error out
	return err
}
