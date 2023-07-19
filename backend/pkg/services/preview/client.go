package preview

import (
	context "context"

	model "github.com/genshinsim/gcsim/pkg/model"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientCfg struct {
	Addr string
}

type Client struct {
	cfg  ClientCfg
	conn EmbedClient
}

func NewClient(cfg ClientCfg, cust ...func(*Client) error) (*Client, error) {
	c := &Client{cfg: cfg}

	conn, err := grpc.Dial(cfg.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	c.conn = NewEmbedClient(conn)

	return c, nil
}

func (c *Client) Get(ctx context.Context, id string, data *model.SimulationResult) ([]byte, error) {
	req := &GetRequest{
		Id:   id,
		Data: data,
	}
	resp, err := c.conn.Get(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.GetData(), nil
}
