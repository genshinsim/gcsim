package db

import (
	context "context"

	"github.com/genshinsim/gcsim/pkg/model"
	"google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	status "google.golang.org/grpc/status"
)

type ClientCfg struct {
	Addr string
}

type Client struct {
	cfg  ClientCfg
	conn DBStoreClient
}

func NewClient(cfg ClientCfg, cust ...func(*Client) error) (*Client, error) {
	c := &Client{cfg: cfg}

	conn, err := grpc.Dial(cfg.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	c.conn = NewDBStoreClient(conn)

	return c, nil
}

func (c *Client) Get(ctx context.Context, query *model.DBQueryOpt) (*model.DBEntries, error) {
	req := &GetRequest{
		Query: query,
	}
	resp, err := c.conn.Get(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.GetData(), nil
}

func (c *Client) GetUnfiltered(ctx context.Context, query *model.DBQueryOpt) (*model.DBEntries, error) {
	req := &GetUnfilteredRequest{
		Query: query,
	}
	resp, err := c.conn.GetUnfiltered(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.GetData(), nil
}

func (c *Client) Update(ctx context.Context, id string, result *model.SimulationResult) error {
	_, err := c.conn.Update(ctx, &UpdateRequest{
		Id:     id,
		Result: result,
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) GetWork(ctx context.Context) ([]*model.ComputeWork, error) {
	res, err := c.conn.GetWork(ctx, &GetWorkRequest{})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			return nil, nil
		}
		return nil, err
	}
	return res.GetData(), nil
}

func (c *Client) ApproveTag(ctx context.Context, id string, tag model.DBTag) error {
	_, err := c.conn.ApproveTag(ctx, &ApproveTagRequest{
		Id:  id,
		Tag: tag,
	})
	return err
}

func (c *Client) RejectTag(ctx context.Context, id string, tag model.DBTag) error {
	_, err := c.conn.RejectTag(ctx, &RejectTagRequest{
		Id:  id,
		Tag: tag,
	})
	return err
}
