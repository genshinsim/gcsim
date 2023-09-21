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

func (s *Server) ReplaceDesc(ctx context.Context, req *ReplaceDescRequest) (*ReplaceDescResponse, error) {
	old, err := s.DBStore.ReplaceDesc(ctx, req.GetId(), req.GetDesc(), req.GetSourceTag())
	if err != nil {
		return nil, err
	}
	s.notify(TopicReplace, &model.DescReplaceEvent{
		DbId:    req.GetId(),
		Desc:    req.GetDesc(),
		OldDesc: old,
	})

	return &ReplaceDescResponse{
		Id: req.GetId(),
	}, nil
}
