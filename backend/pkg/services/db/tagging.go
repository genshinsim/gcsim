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
	if req.GetTag() == model.DBTag_DB_TAG_ADMIN_DO_NOT_USE {
		return nil, status.Error(codes.InvalidArgument, "reserved tag cannot be used")
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
	if req.GetTag() == model.DBTag_DB_TAG_ADMIN_DO_NOT_USE {
		return nil, status.Error(codes.InvalidArgument, "reserved tag cannot be used")
	}
	err := s.DBStore.RejectTag(ctx, req.GetId(), req.GetTag())
	if err != nil {
		return nil, err
	}

	return &RejectTagResponse{Id: req.GetId()}, nil
}

func (s *Server) RejectTagAllUnapproved(ctx context.Context, req *RejectTagAllUnapprovedRequest) (*RejectTagAllUnapprovedResponse, error) {
	if req.GetTag() == model.DBTag_DB_TAG_INVALID {
		return nil, status.Error(codes.InvalidArgument, "tag cannot be blank")
	}
	if req.GetTag() == model.DBTag_DB_TAG_ADMIN_DO_NOT_USE {
		return nil, status.Error(codes.InvalidArgument, "reserved tag cannot be used")
	}
	count, err := s.DBStore.RejectTagAllUnapproved(ctx, req.GetTag())
	if err != nil {
		return nil, err
	}

	return &RejectTagAllUnapprovedResponse{Count: count}, nil
}
