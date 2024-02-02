package backend

import (
	"context"
	"errors"

	"github.com/genshinsim/gcsim/backend/pkg/services/db"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

func (s *Store) GetBySubmitter(id string, page int) ([]*db.Entry, error) {
	sort := bson.M{
		"create_date": -1,
	}
	var err error
	opt := &db.QueryOpt{
		Skip:  int64(page-1) * 5,
		Limit: 5,
	}
	opt.Sort, err = structpb.NewStruct(sort)
	if err != nil {
		s.Log.Infow("error generating sort options", "err", err)
		return nil, err
	}
	resp, err := s.DBClient.GetBySubmitter(context.TODO(), &db.GetBySubmitterRequest{
		Submitter: id,
		Query:     opt,
	})
	if err != nil {
		return nil, err
	}
	return resp.GetData().GetData(), nil
}

func (s *Store) DeletePending(id, sender string) error {
	_, err := s.DBClient.DeletePending(context.TODO(), &db.DeletePendingRequest{
		Id:     id,
		Sender: sender,
	})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return err
		}
		switch st.Code() {
		case codes.PermissionDenied:
			return errors.New(st.Message())
		default:
			return st.Err()
		}
	}
	return nil
}
