package db

import (
	context "context"

	"github.com/genshinsim/gcsim/pkg/model"
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
	// submitter must match sender
	if req.GetSender() != e.Submitter {
		return nil, status.Error(codes.PermissionDenied, "delete failed; this submission is not owned by you")
	}

	// can only delete if it's not dbvalid
	if e.IsDbValid {
		return nil, status.Error(codes.PermissionDenied, "submission already added to db; cannot be deleted")
	}

	err = s.DBStore.Delete(ctx, e.Id)
	if err != nil {
		return nil, err
	}

	s.notify(TopicSubmissionDelete, &model.SubmissionDeleteEvent{
		DbId:      req.GetId(),
		Config:    e.Config,
		Submitter: e.Submitter,
	})

	return &DeletePendingResponse{Id: e.Id}, nil
}
