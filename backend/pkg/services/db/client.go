package db

import (
	context "context"

	"github.com/genshinsim/gcsim/pkg/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientCfg struct {
	Addr string
}

type Client struct {
	cfg       ClientCfg
	srvClient DBStoreClient
}

func NewClient(cfg ClientCfg, cust ...func(*Client) error) (*Client, error) {
	c := &Client{cfg: cfg}

	conn, err := grpc.Dial(cfg.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	c.srvClient = NewDBStoreClient(conn)

	return c, nil
}

func (c *Client) Create(ctx context.Context, e *model.DBEntry) (string, error) {
	req := &CreateOrUpdateDBEntryRequest{
		Data: e,
	}
	resp, err := c.srvClient.CreateOrUpdateDBEntry(ctx, req)
	return resp.GetKey(), err
}

const limit = 30

func (c *Client) Get(ctx context.Context, query *model.DBQueryOpt) (*model.DBEntries, error) {
	req := &GetRequest{
		Query: query,
	}
	resp, err := c.srvClient.Get(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.GetData(), nil
}

func (c *Client) GetComputeWork(ctx context.Context) (*model.ComputeWork, error) {
	req := &GetComputeWorkRequest{}
	resp, err := c.srvClient.GetComputeWork(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.GetWork(), nil
}
