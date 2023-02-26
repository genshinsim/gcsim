package submission

import (
	context "context"

	"github.com/genshinsim/gcsim/pkg/model"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type Client struct {
	conn SubmissionStoreClient
}

func NewClient(addr string) (*Client, error) {
	c := &Client{}

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	c.conn = NewSubmissionStoreClient(conn)

	return c, nil
}

func (c *Client) Submit(ctx context.Context, s *model.Submission) (string, error) {
	res, err := c.conn.Submit(ctx, &SubmitRequest{
		Config:      s.GetConfig(),
		Submitter:   s.GetConfig(),
		Description: s.GetDescription(),
	})
	if err != nil {
		return "", err
	}

	return res.GetId(), nil
}

func (c *Client) Delete(ctx context.Context, id string) error {
	_, err := c.conn.DeletePending(ctx, &DeletePendingRequest{
		Id: id,
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Complete(ctx context.Context, result *model.SimulationResult) error {
	_, err := c.conn.CompletePending(ctx, &CompletePendingRequest{
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
