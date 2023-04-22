package backend

import (
	"context"

	"github.com/genshinsim/gcsim/backend/pkg/services/db"
	"github.com/genshinsim/gcsim/pkg/model"
)

func (s *Store) GetPending(tag model.DBTag) ([]*db.Entry, error) {
	resp, err := s.DBClient.GetPending(context.TODO(), &db.GetPendingRequest{
		Tag: tag,
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
