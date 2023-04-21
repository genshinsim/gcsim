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

func (s *Server) GetUnfiltered(ctx context.Context, req *GetAllRequest) (*GetAllResponse, error) {
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
