package db

import (
	context "context"
)

func (s *Server) Get(ctx context.Context, req *GetRequest) (*GetResponse, error) {
	res, err := s.DBStore.Get(ctx, req.GetQuery())
	if err != nil {
		return nil, err
	}

	return &GetResponse{
		Data: &Entries{
			Data: res,
		},
	}, nil
}

func (s *Server) GetOne(ctx context.Context, req *GetOneRequest) (*GetOneResponse, error) {
	res, err := s.DBStore.GetById(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	return &GetOneResponse{
		Data: res,
	}, nil
}

func (s *Server) GetAll(ctx context.Context, req *GetAllRequest) (*GetAllResponse, error) {
	res, err := s.DBStore.GetAll(ctx, req.GetQuery())
	if err != nil {
		return nil, err
	}

	return &GetAllResponse{
		Data: &Entries{
			Data: res,
		},
	}, nil
}

func (s *Server) GetPending(ctx context.Context, req *GetPendingRequest) (*GetPendingResponse, error) {
	res, err := s.DBStore.GetAllEntriesWithoutTag(ctx, req.GetTag(), req.GetQuery())
	if err != nil {
		return nil, err
	}
	return &GetPendingResponse{
		Data: &Entries{
			Data: res,
		},
	}, err
}

func (s *Server) GetBySubmitter(ctx context.Context, req *GetBySubmitterRequest) (*GetBySubmitterResponse, error) {
	res, err := s.DBStore.GetBySubmitter(ctx, req.GetSubmitter(), req.GetQuery())
	if err != nil {
		return nil, err
	}
	return &GetBySubmitterResponse{
		Data: &Entries{
			Data: res,
		},
	}, err
}
