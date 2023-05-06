package backend

import (
	"context"

	"github.com/genshinsim/gcsim/backend/pkg/services/db"
	"github.com/genshinsim/gcsim/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/protobuf/types/known/structpb"
)

func (s *Store) GetPending(tag model.DBTag, page int) ([]*db.Entry, error) {
	sort := bson.M{
		"create_date": 1,
	}
	var err error
	opt := &db.QueryOpt{
		Skip:  int64(page-1) * 10,
		Limit: 10,
	}
	opt.Sort, err = structpb.NewStruct(sort)
	if err != nil {
		s.Log.Infow("error generating sort options", "err", err)
		return nil, err
	}
	resp, err := s.DBClient.GetPending(context.TODO(), &db.GetPendingRequest{
		Tag:   tag,
		Query: opt,
	})
	if err != nil {
		return nil, err
	}
	return resp.GetData().GetData(), nil
}

func (s *Store) Approve(id string, tag model.DBTag) error {
	_, err := s.DBClient.ApproveTag(context.TODO(), &db.ApproveTagRequest{
		Id:  id,
		Tag: tag,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) Reject(id string, tag model.DBTag) error {
	_, err := s.DBClient.RejectTag(context.TODO(), &db.RejectTagRequest{
		Id:  id,
		Tag: tag,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetRandomSim() string {
	id, err := s.ShareStore.Random(context.TODO())
	if err != nil {
		s.Log.Infow("error getting a random sim", "err", err)
		return ""
	}
	return id
}
