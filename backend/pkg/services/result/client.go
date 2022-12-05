package result

import (
	"context"

	"github.com/genshinsim/gcsim/backend/pkg/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// Client implements api.ResultStore
type Client struct {
	cfg       ClientCfg
	srvClient ResultStoreClient
}

type ClientCfg struct {
	Addr string
}

func NewClient(cfg ClientCfg, cust ...func(*Client) error) (*Client, error) {
	c := &Client{cfg: cfg}

	conn, err := grpc.Dial(cfg.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	c.srvClient = NewResultStoreClient(conn)

	return c, nil
}

func parseError(err error) error {
	st := status.Convert(err)
	switch st.Code() {
	case codes.NotFound:
		return api.ErrKeyNotFound
	case codes.Internal:
		return api.ErrServerError
	case codes.OK:
		return nil
	default:
		return st.Err()
	}
}

func (c *Client) Create(data []byte, ctx context.Context) (string, error) {
	req := &CreateRequest{
		Result: data,
		Ttl:    extractTTL(ctx),
	}
	resp, err := c.srvClient.Create(ctx, req)
	if err != nil {
		return "", parseError(err)

	}
	return resp.GetKey(), nil
}

func (c *Client) Read(id string, ctx context.Context) ([]byte, uint64, error) {
	req := &ReadRequest{
		Key: id,
	}
	resp, err := c.srvClient.Read(ctx, req)
	if err != nil {
		return nil, 0, parseError(err)
	}
	return resp.GetResult(), resp.GetTtl(), nil
}

func (c *Client) Update(id string, data []byte, ctx context.Context) error {
	req := &UpdateRequest{
		Key:    id,
		Result: data,
		Ttl:    extractTTL(ctx),
	}
	_, err := c.srvClient.Update(ctx, req)
	if err != nil {
		return parseError(err)
	}
	return nil
}

func (c *Client) SetTTL(id string, ctx context.Context) error {
	req := &SetTTLRequest{
		Key: id,
		Ttl: extractTTL(ctx),
	}
	_, err := c.srvClient.SetTTL(ctx, req)
	if err != nil {
		return parseError(err)
	}
	return nil
}

func (c *Client) Delete(id string, ctx context.Context) error {
	req := &DeleteRequest{
		Key: id,
	}
	_, err := c.srvClient.Delete(ctx, req)
	if err != nil {
		return parseError(err)
	}
	return nil
}

func (c *Client) Random(ctx context.Context) (string, error) {
	req := &RandomRequest{}
	resp, err := c.srvClient.Random(ctx, req)
	if err != nil {
		return "", parseError(err)

	}
	return resp.GetKey(), nil
}

func extractTTL(ctx context.Context) uint64 {
	x := ctx.Value("ttl")
	//expecting ttl to be an integer value >= 0; if not int then default to
	//14 days; if ttl = 0 then assume to be permanent
	ttl := api.DefaultTLL
	if val, ok := x.(int); ok {
		if val >= 0 {
			ttl = val
		}
	}
	return uint64(ttl)
}
