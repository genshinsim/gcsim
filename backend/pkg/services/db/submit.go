package db

import (
	context "context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) Submit(ctx context.Context, req *SubmitRequest) (*SubmitResponse, error) {

	entry := &Entry{
		Config:      req.GetConfig(),
		Description: req.GetDescription(),
		Submitter:   req.GetSubmitter(),
	}

	id, err := s.DBStore.Create(ctx, entry)
	if err != nil {
		return nil, err
	}

	return &SubmitResponse{
		Id: id,
	}, nil
}

func (s *Server) DeletePending(ctx context.Context, req *DeletePendingRequest) (*DeletePendingResponse, error) {

	e, err := s.DBStore.GetById(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	if e.Summary != nil {
		return nil, status.Error(codes.NotFound, "submission no longer pending")
	}

	err = s.DBStore.Delete(ctx, e.Id)
	if err != nil {
		return nil, err
	}

	return &DeletePendingResponse{Id: e.Id}, nil
}
