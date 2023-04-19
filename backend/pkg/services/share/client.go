package share

import (
	context "context"
	"fmt"
	"time"

	"github.com/genshinsim/gcsim/pkg/model"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientCfg struct {
	Addr string
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

func (c *Client) Create(ctx context.Context, data *model.SimulationResult, ttl uint64, submitter string) (string, uint64, error) {
	var expiry uint64
	if ttl > 0 {
		expiry = uint64(time.Now().Unix()) + ttl
	}

	resp, err := c.srvClient.Create(ctx, &CreateRequest{
		Result:    data,
		ExpiresAt: expiry,
		Submitter: submitter,
	})
	if err != nil {
		return "", 0, err
	}
	return resp.GetId(), expiry, nil
}

func (c *Client) Replace(ctx context.Context, id string, data *model.SimulationResult) error {
	//TODO: handle ttl; should be ok for now since we only call this for db update..
	_, err := c.srvClient.Update(ctx, &UpdateRequest{
		Id:     id,
		Result: data,
	})
	if err != nil {
		return err
	}
	return nil
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
		Id: id,
	})

	if err != nil {
		return nil, 0, err
	}

	return resp.GetResult(), resp.GetExpiresAt(), nil
}
