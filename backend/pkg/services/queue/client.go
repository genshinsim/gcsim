package queue

import (
	context "context"

	"github.com/genshinsim/gcsim/pkg/model"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn WorkQueueClient
}

func NewClient(addr string) (*Client, error) {
	c := &Client{}

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	c.conn = NewWorkQueueClient(conn)

	return c, nil
}

func (c *Client) Add(ctx context.Context, w []*model.ComputeWork) ([]string, error) {
	res, err := c.conn.Populate(ctx, &PopulateReq{
		Data: w,
	})
	if err != nil {
		return nil, err
	}

	return res.GetIds(), nil
}

func (c *Client) Get(ctx context.Context) (*model.ComputeWork, error) {
	res, err := c.conn.Get(ctx, &GetReq{})
	if err != nil {
		return nil, err
	}
	return res.Data, nil
}

func (c *Client) Complete(ctx context.Context, id string) error {
	_, err := c.conn.Complete(ctx, &CompleteReq{Id: id})
	if err != nil {
		return err
	}
	return nil
}
