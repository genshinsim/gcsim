package db

import context "context"

func (s *Server) ReplaceConfig(ctx context.Context, req *ReplaceConfigRequest) (*ReplaceConfigResponse, error) {
	err := s.DBStore.ReplaceConfig(ctx, req.GetId(), req.GetConfig())
	if err != nil {
		return nil, err
	}
	return &ReplaceConfigResponse{
		Id: req.GetId(),
	}, nil
}
