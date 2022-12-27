package db

import (
	context "context"

	"github.com/genshinsim/gcsim/pkg/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/structpb"
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
	req := &CreateRequest{
		Data: e,
	}
	resp, err := c.srvClient.Create(ctx, req)
	return resp.GetKey(), err
}

const limit = 30

func (c *Client) Get(ctx context.Context, query *structpb.Struct, page int64) ([]*model.DBEntry, error) {
	req := &GetRequest{
		Query: query,
		Limit: limit,
		Page:  page,
	}
	resp, err := c.srvClient.Get(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.GetData(), nil
}
