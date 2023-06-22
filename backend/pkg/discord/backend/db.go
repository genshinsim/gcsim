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

func (s *Store) RejectAll(tag model.DBTag) (int64, error) {
	res, err := s.DBClient.RejectTagAllUnapproved(context.TODO(), &db.RejectTagAllUnapprovedRequest{
		Tag: tag,
	})
	if err != nil {
		return 0, err
	}
	return res.GetCount(), nil
}

func (s *Store) GetRandomSim() string {
	id, err := s.ShareStore.Random(context.TODO())
	if err != nil {
		s.Log.Infow("error getting a random sim", "err", err)
		return ""
	}
	return id
}

func (s *Store) GetDBStatus() (*model.DBStatus, error) {
	resp, err := s.DBClient.WorkStatus(context.TODO(), &db.WorkStatusRequest{})
	if err != nil {
		s.Log.Infow("error getting work status", "err", err)
		return nil, err
	}
	return &model.DBStatus{
		DbTotalCount: resp.GetTotalCount(),
		ComputeCount: resp.GetTodoCount(),
	}, nil
}

func (s *Store) ReplaceConfig(id string, link string) error {
	s.Log.Infow("replace config request received", "id", id, "link", link)

	linkid, err := s.validateLink(link)
	if err != nil {
		s.Log.Infow("replace config link validation failed", "err", err)
		return err
	}

	res, _, err := s.ShareStore.Read(context.TODO(), linkid)
	if err != nil {
		s.Log.Infow("replace config getting share failed", "err", err)
		return err
	}

	resp, err := s.DBClient.ReplaceConfig(context.TODO(), &db.ReplaceConfigRequest{
		Id:     id,
		Config: res.Config,
	})

	if err != nil {
		s.Log.Infow("replace config failed", "err", err)
		return err
	}

	s.Log.Infow("replace config completed", "resp", resp.String())

	return nil
}
