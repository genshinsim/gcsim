package db

import (
	context "context"

	"github.com/genshinsim/gcsim/pkg/model"
)

func (s *Server) ReplaceConfig(ctx context.Context, req *ReplaceConfigRequest) (*ReplaceConfigResponse, error) {
	old, err := s.DBStore.ReplaceConfig(ctx, req.GetId(), req.GetConfig(), req.GetSourceTag())
	if err != nil {
		return nil, err
	}
	s.notify(TopicReplace, &model.EntryReplaceEvent{
		DbId:      req.GetId(),
		Config:    req.GetConfig(),
		OldConfig: old,
	})

	return &ReplaceConfigResponse{
		Id: req.GetId(),
	}, nil
}
