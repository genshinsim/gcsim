package db

import (
	context "context"

	"github.com/genshinsim/gcsim/pkg/model"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

func (s *Server) ApproveTag(ctx context.Context, req *ApproveTagRequest) (*ApproveTagResponse, error) {
	//TODO: tag validation should be done at API gateway lvl?? need to check both auth and tag is valid
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id cannot be blank")
	}
	if req.GetTag() == model.DBTag_DB_TAG_INVALID {
		return nil, status.Error(codes.InvalidArgument, "tag cannot be blank")
	}
	err := s.DBStore.ApproveTag(ctx, req.GetId(), req.GetTag())
	if err != nil {
		return nil, err
	}

	return &ApproveTagResponse{Id: req.GetId()}, nil
}

func (s *Server) RejectTag(ctx context.Context, req *RejectTagRequest) (*RejectTagResponse, error) {
	//TODO: tag validation should be done at API gateway lvl?? need to check both auth and tag is valid
	if req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id cannot be blank")
	}
	if req.GetTag() == model.DBTag_DB_TAG_INVALID {
		return nil, status.Error(codes.InvalidArgument, "tag cannot be blank")
	}
	err := s.DBStore.RejectTag(ctx, req.GetId(), req.GetTag())
	if err != nil {
		return nil, err
	}

	return &RejectTagResponse{Id: req.GetId()}, nil
}
